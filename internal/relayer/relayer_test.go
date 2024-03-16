package relayer

// Basic imports
import (
	"errors"
	"github.com/stretchr/testify/suite"
	"pokt_gateway_server/mocks"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"testing"
)

type CachedClientTestSuite struct {
	suite.Suite
	mockAltruistRegistryService *mocks.AltruistRegistryService
	mockSessionRegistryService  *mocks.SessionRegistryService
	mockPocketService           *mocks.PocketService
	cachedClient                *Relayer
}

func (suite *CachedClientTestSuite) SetupTest() {
	suite.mockPocketService = new(mocks.PocketService)
	suite.mockSessionRegistryService = new(mocks.SessionRegistryService)
	// suite.cachedClient = NewRelayer(suite.mockPocketService, suite.mockSessionRegistryService, suite.mockAltruistRegistryService, time.Minute, zap.NewNop())
}

// test SendRelay using table driven tests
func (suite *CachedClientTestSuite) TestSendRelay() {

	testGetSessionRequest := &models.GetSessionRequest{
		AppPubKey: "test",
		Chain:     "test",
	}

	testResponse := &models.GetSessionResponse{}

	testSendRelayResponse := &models.SendRelayResponse{
		Response: "test",
	}

	// create test cases
	testCases := []struct {
		name             string
		request          *models.SendRelayRequest
		setupMocks       func(*models.SendRelayRequest)
		expectedResponse *models.SendRelayResponse
		expectedError    error
	}{
		{
			name: "InvalidRequest",
			request: &models.SendRelayRequest{
				Payload:            nil, // invalid request
				Signer:             nil, // invalid request
				Chain:              "test",
				SelectedNodePubKey: "test",
				Session:            &models.Session{},
			},
			setupMocks: func(request *models.SendRelayRequest) {

				suite.mockPocketService.EXPECT().SendRelay(request).Return(nil, models.ErrMalformedSendRelayRequest).Times(1)

			},
			expectedResponse: nil,
			expectedError:    models.ErrMalformedSendRelayRequest,
		},
		{
			name: "SessionError",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Signer: &models.Ed25519Account{
					PublicKey: "test",
				},
				Chain:              "test",
				SelectedNodePubKey: "test",
			},
			setupMocks: func(request *models.SendRelayRequest) {
				suite.mockSessionRegistryService.EXPECT().GetSession(testGetSessionRequest).Return(nil, errors.New("error")).Times(1)
			},
			expectedResponse: nil,
			expectedError:    errors.New("error"),
		},
		{
			name: "WithSessionInRequestSuccess",
			request: &models.SendRelayRequest{
				Payload:            &models.Payload{},
				Signer:             &models.Ed25519Account{},
				Chain:              "test",
				SelectedNodePubKey: "test",
				Session:            &models.Session{},
			},
			setupMocks: func(request *models.SendRelayRequest) {

				suite.mockPocketService.EXPECT().SendRelay(request).Return(testSendRelayResponse, nil).Times(1)

			},
			expectedResponse: testSendRelayResponse,
			expectedError:    nil,
		},
		{
			name: "WithoutSessionInRequestSuccess",
			request: &models.SendRelayRequest{
				Payload: &models.Payload{},
				Signer: &models.Ed25519Account{
					PublicKey: "test",
				},
				Chain:              "test",
				SelectedNodePubKey: "test",
			},
			setupMocks: func(request *models.SendRelayRequest) {

				suite.mockSessionRegistryService.EXPECT().GetSession(testGetSessionRequest).Return(testResponse, nil).Times(1)

				suite.mockPocketService.EXPECT().SendRelay(request).Return(testSendRelayResponse, nil).Times(1)

			},
			expectedResponse: testSendRelayResponse,
			expectedError:    nil,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks(tc.request) // setup mocks

			session, err := suite.cachedClient.SendRelay(tc.request)

			// assert results
			suite.Equal(tc.expectedResponse, session)
			suite.Equal(tc.expectedError, err)

		})
	}

}

func TestCachedClientTestSuite(t *testing.T) {
	suite.Run(t, new(CachedClientTestSuite))
}
