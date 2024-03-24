package controllers

// Basic imports
import (
	"errors"
	"pokt_gateway_server/mocks"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type RelayTestSuite struct {
	suite.Suite
	mockPocketService   *mocks.PocketService
	mockRelayController *RelayController
	context             *fasthttp.RequestCtx
}

func (suite *RelayTestSuite) SetupTest() {
	suite.mockPocketService = new(mocks.PocketService)
	suite.mockRelayController = NewRelayController(suite.mockPocketService, zap.NewNop())
	suite.context = &fasthttp.RequestCtx{} // mock the fasthttp.RequestCtx
}

// mock send relay request function
func (suite *RelayTestSuite) mockSendRelayRequest() *models.SendRelayRequest {

	chainID, path := getPathSegmented(suite.context.Path()) // get the chainID and path from the request path

	return &models.SendRelayRequest{
		Payload: &models.Payload{
			Data:   string(suite.context.PostBody()),
			Method: string(suite.context.Method()),
			Path:   path,
		},
		Chain: chainID,
	}
}

// test for the HandleRelay function in relay.go file using table driven tests to test different scenarios for the function
func (suite *RelayTestSuite) TestHandleRelay() {

	var testResponse string = "test"

	tests := []struct {
		name             string
		setupMocks       func(*fasthttp.RequestCtx)
		path             string
		expectedStatus   int
		expectedResponse *string
	}{
		{
			name: "EmptyChainID",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
			},
			path:             "/relay/",
			expectedStatus:   fasthttp.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "ChainIdLengthInvalid",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
			},
			path:             "/relay/1234555",
			expectedStatus:   fasthttp.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "ErrorSendingRelay",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
				suite.mockPocketService.EXPECT().SendRelay(suite.mockSendRelayRequest()).
					Return(nil, errors.New("relay error"))
			},
			path:             "/relay/1234",
			expectedStatus:   fasthttp.StatusInternalServerError,
			expectedResponse: nil,
		},
		{
			name: "Success",
			setupMocks: func(ctx *fasthttp.RequestCtx) {
				suite.mockPocketService.EXPECT().SendRelay(suite.mockSendRelayRequest()).
					Return(&models.SendRelayResponse{
						Response: testResponse,
					}, nil)

			},
			path:             "/relay/1234",
			expectedStatus:   fasthttp.StatusOK,
			expectedResponse: &testResponse,
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

			suite.Equal(test.expectedStatus, suite.context.Response.StatusCode())

			if test.expectedResponse != nil {
				suite.Equal(*test.expectedResponse, string(suite.context.Response.Body()))
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
