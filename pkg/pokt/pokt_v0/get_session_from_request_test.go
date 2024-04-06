package pokt_v0

import (
	"errors"
	pocket_service_mock "github.com/pokt-network/gateway-server/mocks/pocket_service"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetSessionFromRequest(t *testing.T) {

	mockSession := &models.Session{
		Nodes: []*models.Node{
			{
				ServiceUrl: "example-node-1",
				PublicKey:  "example-pub-key-1",
			},
		},
		SessionHeader: &models.SessionHeader{
			SessionHeight: uint(1),
		},
	}

	mockErr := errors.New("failure")

	type args struct {
		pocketService PocketService
		req           *models.SendRelayRequest
	}
	tests := []struct {
		name            string
		generateArgs    func() args
		expectedSession *models.Session
		expectedErr     error
	}{
		{
			name: "SessionFromInnerRequest",
			generateArgs: func() args {
				return args{
					pocketService: nil,
					req:           &models.SendRelayRequest{Session: mockSession},
				}
			},
			expectedErr:     nil,
			expectedSession: mockSession,
		},
		{
			name: "SessionFromPocketServiceError",
			generateArgs: func() args {
				mockPocketService := new(pocket_service_mock.PocketService)
				mockPocketService.EXPECT().GetSession(mock.Anything).Return(nil, mockErr).Times(1)
				return args{
					pocketService: mockPocketService,
					req:           &models.SendRelayRequest{Session: nil, Signer: &models.Ed25519Account{}},
				}
			},
			expectedErr:     mockErr,
			expectedSession: nil,
		},
		{
			name: "SessionFromPocketServiceSuccess",
			generateArgs: func() args {
				mockPocketService := new(pocket_service_mock.PocketService)
				mockPocketService.EXPECT().GetSession(mock.Anything).Return(&models.GetSessionResponse{Session: mockSession}, nil).Times(1)
				return args{
					pocketService: mockPocketService,
					req:           &models.SendRelayRequest{Session: nil, Signer: &models.Ed25519Account{}},
				}
			},
			expectedErr:     nil,
			expectedSession: mockSession,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.generateArgs()
			session, err := GetSessionFromRequest(args.pocketService, args.req)
			assert.Equal(t, err, tt.expectedErr)
			assert.Equal(t, session, tt.expectedSession)
		})
	}
}
