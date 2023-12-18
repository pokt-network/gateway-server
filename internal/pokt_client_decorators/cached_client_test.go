package pokt_client_decorators

// Basic imports
import (
	"errors"
	"os-gateway/mocks"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/suite"
)

type CachedClientTestSuite struct {
	suite.Suite
	mockPocketService    *mocks.PocketService
	newCachedClient      *CachedClient
	mockTTLCachedService *mocks.TTLCacheService[string, *models.GetSessionResponse]
}

func (suite *CachedClientTestSuite) SetupTest() {
	suite.mockPocketService = new(mocks.PocketService)
	suite.mockTTLCachedService = new(mocks.TTLCacheService[string, *models.GetSessionResponse])
	suite.newCachedClient = NewCachedClient(suite.mockPocketService, suite.mockTTLCachedService)
}

// test GetSession using table driven tests
func (suite *CachedClientTestSuite) TestGetSession() {

	testRequest := &models.GetSessionRequest{
		AppPubKey: "test",
		Chain:     "test",
		Height:    1,
	}

	ttlcacheItem := &ttlcache.Item[string, *models.GetSessionResponse]{}

	testResponse := &models.GetSessionResponse{}

	errUnderlayingProvider := errors.New("error underlaying provider")

	// create test cases
	testCases := []struct {
		name             string
		setupMocks       func()
		expectedResponse *models.GetSessionResponse
		expectedError    error
	}{
		{
			name: "NotCached",
			setupMocks: func() {

				suite.newCachedClient.lastFailure = time.Now()

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

			},
			expectedResponse: nil,
			expectedError:    ErrRecentlyFailed,
		},
		{
			name: "Cached",
			setupMocks: func() {

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

				suite.mockPocketService.EXPECT().GetSession(testRequest).Return(testResponse, nil)

				suite.mockTTLCachedService.EXPECT().Set("test-test", testResponse, ttlcache.DefaultTTL).Return(ttlcacheItem)

			},
			expectedResponse: testResponse,
			expectedError:    nil,
		},
		{
			name: "Error",
			setupMocks: func() {

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

				suite.mockPocketService.EXPECT().GetSession(testRequest).Return(nil, errUnderlayingProvider)

			},
			expectedResponse: nil,
			expectedError:    errUnderlayingProvider,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks() // setup mocks

			session, err := suite.newCachedClient.GetSession(testRequest)

			// assert results
			suite.Equal(tc.expectedResponse, session)
			suite.Equal(tc.expectedError, err)

		})
	}

}

// test SendRelay using table driven tests
func (suite *CachedClientTestSuite) TestSendRelay() {

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

				suite.mockPocketService.EXPECT().SendRelay(request).Return(nil, models.ErrMalformedSendRelayRequest)

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

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(&ttlcache.Item[string, *models.GetSessionResponse]{})

				suite.mockPocketService.EXPECT().GetSession(&models.GetSessionRequest{
					AppPubKey: "test",
					Chain:     "test",
				}).Return(nil, errors.New("error"))

			},
			expectedResponse: nil,
			expectedError:    errors.New("error"),
		},
		{
			name: "Success",
			request: &models.SendRelayRequest{
				Payload:            &models.Payload{},
				Signer:             &models.Ed25519Account{},
				Chain:              "test",
				SelectedNodePubKey: "test",
				Session:            &models.Session{},
			},
			setupMocks: func(request *models.SendRelayRequest) {

				suite.mockPocketService.EXPECT().SendRelay(request).Return(testSendRelayResponse, nil)

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

			session, err := suite.newCachedClient.SendRelay(tc.request)

			// assert results
			suite.Equal(tc.expectedResponse, session)
			suite.Equal(tc.expectedError, err)

		})
	}

}

func TestCachedClientTestSuite(t *testing.T) {
	suite.Run(t, new(CachedClientTestSuite))
}
