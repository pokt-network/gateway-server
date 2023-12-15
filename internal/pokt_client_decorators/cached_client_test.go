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
		name           string
		setupMocks     func()
		expectedResult *models.GetSessionResponse
		expectedError  error
	}{
		{
			name: "NotCached",
			setupMocks: func() {

				suite.newCachedClient.lastFailure = time.Now()

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

			},
			expectedResult: nil,
			expectedError:  ErrRecentlyFailed,
		},
		{
			name: "Cached",
			setupMocks: func() {

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

				suite.mockPocketService.EXPECT().GetSession(testRequest).Return(testResponse, nil)

				suite.mockTTLCachedService.EXPECT().Set("test-test", testResponse, ttlcache.DefaultTTL).Return(ttlcacheItem)

			},
			expectedResult: &models.GetSessionResponse{},
			expectedError:  nil,
		},
		{
			name: "Error",
			setupMocks: func() {

				suite.mockTTLCachedService.EXPECT().Get("test-test").Return(ttlcacheItem)

				suite.mockPocketService.EXPECT().GetSession(testRequest).Return(nil, errUnderlayingProvider)

			},
			expectedResult: nil,
			expectedError:  errUnderlayingProvider,
		},
	}

	// run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			suite.SetupTest() // reset mocks

			tc.setupMocks() // setup mocks

			session, err := suite.newCachedClient.GetSession(testRequest)

			// assert results
			suite.Equal(tc.expectedResult, session)
			suite.Equal(tc.expectedError, err)

		})
	}

}

func TestCachedClientTestSuite(t *testing.T) {
	suite.Run(t, new(CachedClientTestSuite))
}
