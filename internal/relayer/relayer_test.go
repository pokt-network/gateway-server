package relayer

// Basic imports
import (
	"github.com/jackc/pgtype"
	"github.com/pokt-network/gateway-server/internal/db_query"
	qos_models "github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	apps_registry_mock "github.com/pokt-network/gateway-server/mocks/apps_registry"
	chain_configurations_registry_mock "github.com/pokt-network/gateway-server/mocks/chain_configurations_registry"
	global_config_mock "github.com/pokt-network/gateway-server/mocks/global_config"
	node_selector_mock "github.com/pokt-network/gateway-server/mocks/node_selector"
	pocket_service_mock "github.com/pokt-network/gateway-server/mocks/pocket_service"
	session_registry_mock "github.com/pokt-network/gateway-server/mocks/session_registry"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"time"

	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"testing"
)

type RelayerTestSuite struct {
	suite.Suite
	mockNodeSelectorService        *node_selector_mock.NodeSelectorService
	mockChainConfigurationsService *chain_configurations_registry_mock.ChainConfigurationsService
	mockSessionRegistryService     *session_registry_mock.SessionRegistryService
	mockPocketService              *pocket_service_mock.PocketService
	mockAppRegistry                *apps_registry_mock.AppsRegistryService
	mockConfigProvider             *global_config_mock.GlobalConfigProvider
	relayer                        *Relayer
}

func (suite *RelayerTestSuite) SetupTest() {
	suite.mockPocketService = new(pocket_service_mock.PocketService)
	suite.mockNodeSelectorService = new(node_selector_mock.NodeSelectorService)
	suite.mockSessionRegistryService = new(session_registry_mock.SessionRegistryService)
	suite.mockChainConfigurationsService = new(chain_configurations_registry_mock.ChainConfigurationsService)
	suite.mockAppRegistry = new(apps_registry_mock.AppsRegistryService)
	suite.mockConfigProvider = new(global_config_mock.GlobalConfigProvider)
	suite.relayer = NewRelayer(suite.mockPocketService, suite.mockSessionRegistryService, suite.mockAppRegistry, suite.mockNodeSelectorService, suite.mockChainConfigurationsService, "", suite.mockConfigProvider, zap.NewNop())
}

func (suite *RelayerTestSuite) TestNodeSelectorRelay() {

	expectedResponse := &models.SendRelayResponse{Response: "response"}
	// create test cases
	testCases := []struct {
		name             string
		request          *models.SendRelayRequest
		setupMocks       func(*models.SendRelayRequest)
		expectedResponse *models.SendRelayResponse
		expectedNodeHost string
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
			expectedNodeHost: "",
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
				node := &models.Node{PublicKey: "123", ServiceUrl: "http://complex.subdomain.root.com/test/123"}
				session := &models.Session{}
				suite.mockNodeSelectorService.EXPECT().FindNode("1234").Return(qos_models.NewQosNode(node, session, signer), true)
				// expect sendRelay to have same parameters as find node, otherwise validation will fail
				suite.mockPocketService.EXPECT().SendRelay(&models.SendRelayRequest{
					Payload:            request.Payload,
					Signer:             signer,
					Chain:              request.Chain,
					SelectedNodePubKey: node.PublicKey,
					Session:            session,
				}).Return(expectedResponse, nil)
			},
			expectedNodeHost: "root.com",
			expectedResponse: expectedResponse,
			expectedError:    nil,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks(tc.request) // setup mocks

			rsp, host, err := suite.relayer.sendNodeSelectorRelay(tc.request)

			// assert results
			suite.Equal(tc.expectedResponse, rsp)
			suite.Equal(tc.expectedNodeHost, host)
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
				suite.mockChainConfigurationsService.EXPECT().GetChainConfiguration(request.Chain).Return(db_query.GetChainConfigurationsRow{}, false)
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
				pgTypeStr := pgtype.Varchar{}
				pgTypeStr.Set("https://example.com")
				// We can only check if altruist url and if proper config is called
				suite.mockConfigProvider.EXPECT().GetAltruistRequestTimeout().Return(time.Second * 15)
				suite.mockChainConfigurationsService.EXPECT().GetChainConfiguration(request.Chain).Return(db_query.GetChainConfigurationsRow{AltruistUrl: pgTypeStr}, true)
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
