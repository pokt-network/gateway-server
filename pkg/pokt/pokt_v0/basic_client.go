package pokt_v0

import (
	"encoding/hex"
	"errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"math/rand"
	"os-gateway/pkg/common"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"slices"
	"time"
)

const (
	endpointClientPrefix = "/v1/client"
	endpointQueryPrefix  = "/v1/query"
	endpointDispatch     = endpointClientPrefix + "/dispatch"
	endpointSendRelay    = endpointClientPrefix + "/relay"
	endpointGetHeight    = endpointQueryPrefix + "/height"
)

var (
	ErrMissingFullNodes          = errors.New("require full node host")
	ErrSessionHasZeroNodes       = errors.New("session missing valid nodes")
	ErrNodeNotFound              = errors.New("node not found")
	ErrMalformedSendRelayRequest = errors.New("malformed send relay request")
)

// BasicClient represents a basic client with a logging, full node host, and a global request timeout.
type BasicClient struct {
	fullNodeHost         string
	globalRequestTimeout time.Duration
}

// NewBasicClient creates a new BasicClient instance.
// Parameters:
//   - fullNodeHost: Full node host address.
//   - logging: Logger instance.
//   - timeout: Global request timeout duration.
//
// Returns:
//   - (*BasicClient): New BasicClient instance.
//   - (error): Error, if any.
func NewBasicClient(fullNodeHost string, timeout time.Duration) (*BasicClient, error) {
	if len(fullNodeHost) == 0 {
		return nil, ErrMissingFullNodes
	}
	return &BasicClient{
		fullNodeHost:         fullNodeHost,
		globalRequestTimeout: timeout,
	}, nil
}

// GetSession obtains a session from the full node.
// Parameters:
//   - req: GetSessionRequest instance containing the request parameters.
//
// Returns:
//   - (*GetSessionResponse): Session response.
//   - (error): Error, if any.
func (r BasicClient) GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error) {
	var sessionResponse models.GetSessionResponse
	err := r.makeRequest(endpointDispatch, "POST", req, &sessionResponse, nil)
	if err != nil {
		return nil, err
	}
	return &sessionResponse, nil
}

// SendRelay sends a relay request to the full node.
// Parameters:
//   - req: SendRelayRequest instance containing the relay request parameters.
//
// Returns:
//   - (*SendRelayResponse): Relay response.
//   - (error): Error, if any.
func (r BasicClient) SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

	if req.Payload == nil || req.Signer == nil {
		return nil, ErrMalformedSendRelayRequest
	}
	// Get a session from the request or retrieve from full node
	session, err := r.getSessionFromRequest(req)

	if err != nil {
		return nil, err
	}

	// Get the preferred selected node, or chose a random one.
	node, err := r.getNodeFromRequest(session, req.SelectedNodePubKey)

	if err != nil {
		return nil, err
	}

	currentSessionHeight := session.SessionHeader.SessionHeight

	relayMetadata := &models.RelayMeta{BlockHeight: currentSessionHeight}

	relayProof := r.generateRelayProof(req.Chain, currentSessionHeight, node.PublicKey, relayMetadata, req.Payload, req.Signer)

	// Relay created, generating a request to the servicer
	var sessionResponse models.SendRelayResponse
	err = r.makeRequest(endpointSendRelay, "POST", &models.Relay{
		Payload:    req.Payload,
		Metadata:   relayMetadata,
		RelayProof: &relayProof,
	}, &sessionResponse, &node.ServiceUrl)

	if err != nil {
		return nil, err
	}

	return &sessionResponse, nil
}

// GetLatestBlockHeight gets the latest block height from the full node.
// Returns:
//   - (*GetLatestBlockHeightResponse): Latest block height response.
//   - (error): Error, if any.
func (r BasicClient) GetLatestBlockHeight() (*models.GetLatestBlockHeightResponse, error) {

	var height models.GetLatestBlockHeightResponse
	err := r.makeRequest(endpointGetHeight, "POST", nil, &height, nil)

	if err != nil {
		return nil, err
	}

	return &height, nil
}

func (r BasicClient) makeRequest(endpoint string, method string, requestData any, responseModel any, hostOverride *string) error {
	reqPayload, err := ffjson.Marshal(requestData)
	if err != nil {
		return err
	}

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	if hostOverride != nil {
		request.SetRequestURI(*hostOverride + endpoint)
	} else {
		request.SetRequestURI(r.fullNodeHost + endpoint)
	}
	request.Header.SetMethod(method)

	if method == "POST" {
		request.SetBody(reqPayload)
	}

	err = fasthttp.DoTimeout(request, response, r.globalRequestTimeout)
	if err != nil {
		return err
	}

	// Check for a successful HTTP status code
	if !common.IsHttpOk(response.StatusCode()) {
		var pocketError models.PocketRPCError
		err := ffjson.Unmarshal(response.Body(), pocketError)
		// failed to unmarshal, not sure what the response code is
		if err != nil {
			return models.PocketRPCError{HttpCode: uint64(response.StatusCode()), Message: string(response.Body())}
		}
		return pocketError
	}
	return ffjson.Unmarshal(response.Body(), responseModel)
}

// generateRelayProof generates a relay proof.
// Parameters:
//   - chainId: Blockchain ID.
//   - sessionHeight: Session block height.
//   - servicerPubKey: Servicer public key.
//   - requestMetadata: Request metadata.
//   - account: Ed25519 account used for signing.
//
// Returns:
//   - models.RelayProof: Generated relay proof.
func (r BasicClient) generateRelayProof(chainId string, sessionHeight uint, servicerPubKey string, relayMetadata *models.RelayMeta, reqPayload *models.Payload, account *models.Ed25519Account) models.RelayProof {
	entropy := uint64(rand.Int63())
	aat := account.GetAAT()

	requestMetadata := models.RequestHashPayload{
		Metadata: relayMetadata,
		Payload:  reqPayload,
	}

	requestHash := requestMetadata.Hash()

	unsignedAAT := &models.AAT{
		Version:      aat.Version,
		AppPubKey:    aat.AppPubKey,
		ClientPubKey: aat.ClientPubKey,
		Signature:    "",
	}

	proofObj := &models.RelayProofHashPayload{
		RequestHash:        requestHash,
		Entropy:            entropy,
		SessionBlockHeight: sessionHeight,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         chainId,
		Signature:          "",
		UnsignedAAT:        unsignedAAT.Hash(),
	}

	hashedPayload := common.Sha3_256Hash(proofObj)
	hashSignature := hex.EncodeToString(account.Sign(hashedPayload))
	return models.RelayProof{
		RequestHash:        requestHash,
		Entropy:            entropy,
		SessionBlockHeight: sessionHeight,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         chainId,
		AAT:                aat,
		Signature:          hashSignature,
	}
}

// getSessionFromRequest obtains a session from a relay request.
// Parameters:
//   - req: SendRelayRequest instance containing the relay request parameters.
//
// Returns:
//   - (*GetSessionResponse): Session response.
//   - (error): Error, if any.
func (r BasicClient) getSessionFromRequest(req *models.SendRelayRequest) (*models.Session, error) {
	if req.Session != nil {
		return req.Session, nil
	}
	sessionResp, err := r.GetSession(&models.GetSessionRequest{
		AppPubKey: req.Signer.PublicKey,
		Chain:     req.Chain,
	})
	if err != nil {
		return nil, err
	}
	return sessionResp.Session, nil
}

// getNodeFromRequest obtains a node from a relay request.
// Parameters:
//   - req: SendRelayRequest instance containing the relay request parameters.
//
// Returns:
//   - (*models.Node): Node instance.
//   - (error): Error, if any.
func (r BasicClient) getNodeFromRequest(session *models.Session, selectedNodePubKey string) (*models.Node, error) {
	if selectedNodePubKey == "" {
		return getRandomNodeOrError(session.Nodes, ErrSessionHasZeroNodes)
	}
	return findNodeOrError(session.Nodes, selectedNodePubKey, ErrNodeNotFound)
}

// getRandomNodeOrError gets a random node or returns an error if the node list is empty.
// Parameters:
//   - nodes: List of nodes.
//   - err: Error to be returned if the node list is empty.
//
// Returns:
//   - (*models.Node): Random node.
//   - (error): Error, if any.
func getRandomNodeOrError(nodes []*models.Node, err error) (*models.Node, error) {
	node := common.GetRandomElement(nodes)
	if node == nil {
		return nil, err
	}
	return node, nil
}

// findNodeOrError finds a node by public key or returns an error if the node is not found.
// Parameters:
//   - nodes: List of nodes.
//   - pubKey: Public key of the node to find.
//   - err: Error to be returned if the node is not found.
//
// Returns:
//   - (*models.Node): Found node.
//   - (error): Error, if any.
func findNodeOrError(nodes []*models.Node, pubKey string, err error) (*models.Node, error) {
	idx := slices.IndexFunc(nodes, func(node *models.Node) bool {
		return node.PublicKey == pubKey
	})
	if idx == -1 {
		return nil, err
	}
	return nodes[idx], nil
}
