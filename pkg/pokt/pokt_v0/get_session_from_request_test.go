package pokt_v0

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"
)

func TestGetSessionFromRequest(t *testing.T) {
	type args struct {
		pocketService PocketService
		req           *models.SendRelayRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Session
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSessionFromRequest(tt.args.pocketService, tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("GetSessionFromRequest(%v, %v)", tt.args.pocketService, tt.args.req)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetSessionFromRequest(%v, %v)", tt.args.pocketService, tt.args.req)
		})
	}
}
