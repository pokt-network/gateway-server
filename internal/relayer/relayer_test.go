package relayer

// Basic imports
import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	qos_models "pokt_gateway_server/internal/node_selector_service/models"
	altruist_registry_mock "pokt_gateway_server/mocks/altruist_registry"
	apps_registry_mock "pokt_gateway_server/mocks/apps_registry"
	node_selector_mock "pokt_gateway_server/mocks/node_selector"
	pocket_service_mock "pokt_gateway_server/mocks/pocket_service"
	session_registry_mock "pokt_gateway_server/mocks/session_registry"

	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"testing"
	"time"
)

type RelayerTestSuite struct {
	suite.Suite
	mockNodeSelectorService     *node_selector_mock.NodeSelectorService
	mockAltruistRegistryService *altruist_registry_mock.AltruistRegistryService
	mockSessionRegistryService  *session_registry_mock.SessionRegistryService
	mockPocketService           *pocket_service_mock.PocketService
	mockAppRegistry             *apps_registry_mock.AppsRegistryService
	relayer                     *Relayer
}

func (suite *RelayerTestSuite) SetupTest() {
	suite.mockPocketService = new(pocket_service_mock.PocketService)
	suite.mockNodeSelectorService = new(node_selector_mock.NodeSelectorService)
	suite.mockSessionRegistryService = new(session_registry_mock.SessionRegistryService)
	suite.mockAltruistRegistryService = new(altruist_registry_mock.AltruistRegistryService)
	suite.mockAppRegistry = new(apps_registry_mock.AppsRegistryService)
	suite.relayer = NewRelayer(suite.mockPocketService, suite.mockSessionRegistryService, suite.mockAppRegistry, suite.mockNodeSelectorService, suite.mockAltruistRegistryService, time.Minute, zap.NewNop())
}

func (suite *RelayerTestSuite) TestNodeSelectorRelay() {

	expectedResponse := &models.SendRelayResponse{Response: "response"}
	// create test cases
	testCases := []struct {
		name             string
		request          *models.SendRelayRequest
		setupMocks       func(*models.SendRelayRequest)
		expectedResponse *models.SendRelayResponse
		expectedError    error
	}{
		{
			name: "NodeSelectorFailed",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Chain:   "1234",
			},
			setupMocks: func(request *models.SendRelayRequest) {
				suite.mockNodeSelectorService.EXPECT().FindNode("1234").Return(nil, false)
			},
			expectedResponse: nil,
			expectedError:    errSelectNodeFail,
		},
		{
			name: "Success",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Chain:   "1234",
			},
			setupMocks: func(request *models.SendRelayRequest) {

				signer := &models.Ed25519Account{}
				node := &models.Node{PublicKey: "123"}
				session := &models.Session{}
				suite.mockNodeSelectorService.EXPECT().FindNode("1234").Return(&qos_models.QosNode{
					AppSigner:     signer,
					MorseNode:     node,
					PocketSession: session,
				}, true)
				// expect sendRelay to have same parameters as find node, otherwise validation will fail
				suite.mockPocketService.EXPECT().SendRelay(&models.SendRelayRequest{
					Payload:            request.Payload,
					Signer:             signer,
					Chain:              request.Chain,
					SelectedNodePubKey: node.PublicKey,
					Session:            session,
				}).Return(expectedResponse, nil)
			},
			expectedResponse: expectedResponse,
			expectedError:    nil,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks(tc.request) // setup mocks

			rsp, err := suite.relayer.sendNodeSelectorRelay(tc.request)

			// assert results
			suite.Equal(tc.expectedResponse, rsp)
			suite.Equal(tc.expectedError, err)

		})
	}

}

// test TestNodeSelectorRelay using table driven tests
func (suite *RelayerTestSuite) TestAltruistRelay() {

	// create test cases
	testCases := []struct {
		name             string
		request          *models.SendRelayRequest
		setupMocks       func(*models.SendRelayRequest)
		expectedResponse *models.SendRelayResponse
		expectedError    error
	}{
		{
			name: "Altruist Missing",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Chain:   "1234",
			},
			setupMocks: func(request *models.SendRelayRequest) {
				suite.mockAltruistRegistryService.EXPECT().GetAltruistURL(request.Chain).Return("", false)
			},
			expectedResponse: nil,
			expectedError:    errAltruistNotFound,
		},
		{
			name: "Altruist Registry successfully called",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Chain:   "1234",
			},
			setupMocks: func(request *models.SendRelayRequest) {
				// We can only check if altruist url
				suite.mockAltruistRegistryService.EXPECT().GetAltruistURL(request.Chain).Return("https://chain.com", true)
			},
			expectedResponse: nil,
			expectedError:    nil,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks(tc.request) // setup mocks

			_, err := suite.relayer.altruistRelay(tc.request)

			// Check if error matches expected
			suite.Equal(tc.expectedError, err)

		})
	}

}

func TestRelayerTestSuite(t *testing.T) {
	suite.Run(t, new(RelayerTestSuite))
}
