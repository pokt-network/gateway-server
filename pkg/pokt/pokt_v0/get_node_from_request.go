package pokt_v0

import (
	"github.com/pokt-network/gateway-server/pkg/common"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"slices"
)

// getNodeFromRequest obtains a node from a relay request.
// Parameters:
//   - req: SendRelayRequest instance containing the relay request parameters.
//
// Returns:
//   - (*models.Node): Node instance.
//   - (error): Error, if any.
func getNodeFromRequest(session *models.Session, selectedNodePubKey string) (*models.Node, error) {
	if selectedNodePubKey == "" {
		return getRandomNodeOrError(session.Nodes, models.ErrSessionHasZeroNodes)
	}
	return findNodeOrError(session.Nodes, selectedNodePubKey, models.ErrNodeNotFound)
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
	node, ok := common.GetRandomElement(nodes)
	if !ok || node == nil {
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
