package controllers

// Basic imports
import (
	"errors"
	"os-gateway/mocks"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type RelayTestSuite struct {
	suite.Suite
	mockPocketService       *mocks.PocketService
	mockRelayController     *RelayController
	context                 *fasthttp.RequestCtx
	mockAppStakePrivateKeys []*models.Ed25519Account
}

func (suite *RelayTestSuite) SetupTest() {
	suite.mockPocketService = new(mocks.PocketService)
	suite.context = &fasthttp.RequestCtx{} // mock the fasthttp.RequestCtx
	suite.mockAppStakePrivateKeys = mockAppStakePrivateKeys()
}

// mock app stake private keys
func mockAppStakePrivateKeys() []*models.Ed25519Account {

	var appStakePrivateKeys []*models.Ed25519Account

	appStake, _ := models.NewAccount("3fe64039816c44e8872e4ef981725b968422e3d49e95a1eb800707591df30fe374039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849")

	appStakePrivateKeys = append(appStakePrivateKeys, appStake)

	return appStakePrivateKeys

}

// test for the HandleRelay function in relay.go file using table driven tests to test different scenarios for the function
func (suite *RelayTestSuite) TestHandleRelay() {
	// table driven tests
	tests := []struct {
		name             string
		setupMocks       func(*fasthttp.RequestCtx)
		path             string
		expectedSatus    int
		expectedResponse func() *string
	}{
		{
			name: "EmptyChainID",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
				suite.mockRelayController = NewRelayController(suite.mockPocketService, suite.mockAppStakePrivateKeys, zap.NewNop())
			},
			path:          "/relay/",
			expectedSatus: fasthttp.StatusBadRequest,
			expectedResponse: func() *string {
				return nil
			},
		},
		{
			name: "AppStakeNotProvided",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
				suite.mockRelayController = NewRelayController(suite.mockPocketService, []*models.Ed25519Account{}, zap.NewNop())
			},
			path:          "/relay/1234",
			expectedSatus: fasthttp.StatusInternalServerError,
			expectedResponse: func() *string {
				return nil
			},
		},
		{
			name: "ErrorSendingRelayRequest",
			setupMocks: func(ctx *fasthttp.RequestCtx) {

				chainID, path := getPathSegmented(ctx.Path())

				mockAppStakePrivateKeys := suite.mockAppStakePrivateKeys

				suite.mockRelayController = NewRelayController(suite.mockPocketService, mockAppStakePrivateKeys, zap.NewNop())

				suite.mockPocketService.EXPECT().SendRelay(&models.SendRelayRequest{
					Payload: &models.Payload{
						Data:   string(ctx.PostBody()),
						Method: string(ctx.Method()),
						Path:   path,
					},
					Signer: mockAppStakePrivateKeys[0],
					Chain:  chainID,
				}).Return(nil, errors.New("error"))

			},
			path:          "/relay/1234",
			expectedSatus: fasthttp.StatusInternalServerError,
			expectedResponse: func() *string {
				return nil
			},
		},
		{
			name: "SuccessfullRelayRequest",
			setupMocks: func(ctx *fasthttp.RequestCtx) {

				chainID, path := getPathSegmented(ctx.Path())

				mockAppStakePrivateKeys := suite.mockAppStakePrivateKeys

				suite.mockRelayController = NewRelayController(suite.mockPocketService, mockAppStakePrivateKeys, zap.NewNop())

				suite.mockPocketService.EXPECT().SendRelay(&models.SendRelayRequest{
					Payload: &models.Payload{
						Data:   string(ctx.PostBody()),
						Method: string(ctx.Method()),
						Path:   path,
					},
					Signer: mockAppStakePrivateKeys[0],
					Chain:  chainID,
				}).Return(&models.SendRelayResponse{
					Response: "test",
				}, nil)

			},
			path:          "/relay/1234",
			expectedSatus: fasthttp.StatusOK,
			expectedResponse: func() *string {

				response := "test"

				return &response

			},
		},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {

			suite.SetupTest() // reset the test suite

			suite.context.Request.SetBody([]byte("test"))
			suite.context.Request.Header.SetMethod("POST")
			suite.context.Request.SetRequestURI(test.path)

			test.setupMocks(suite.context) // setup the mocks for the test

			suite.mockRelayController.HandleRelay(suite.context)

			suite.Equal(test.expectedSatus, suite.context.Response.StatusCode())

			if test.expectedResponse() != nil {
				suite.Equal(*test.expectedResponse(), string(suite.context.Response.Body()))
			}

		})
	}
}

// test for the getPathSegmented function in relay.go file using table driven tests to test different scenarios for the function
func (suite *RelayTestSuite) TestGetPathSegmented() {

	tests := []struct {
		name         string
		path         string
		expectedPath string
		expectedRest string
	}{
		{
			name:         "EmptyPath",
			path:         "",
			expectedPath: "",
			expectedRest: "",
		},
		{
			name:         "LessThanTwoSegments",
			path:         "/segment1",
			expectedPath: "",
			expectedRest: "",
		},
		{
			name:         "TwoSegments",
			path:         "/segment1/1234",
			expectedPath: "1234",
			expectedRest: "",
		},
		{
			name:         "MoreThanTwoSegments",
			path:         "/segment1/1234/segment2",
			expectedPath: "1234",
			expectedRest: "/segment2",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {

			path, rest := getPathSegmented([]byte(test.path))

			suite.Equal(test.expectedPath, path)
			suite.Equal(test.expectedRest, rest)

		})
	}

}

func TestRelayTestSuite(t *testing.T) {
	suite.Run(t, new(RelayTestSuite))
}
