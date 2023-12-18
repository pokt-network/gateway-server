package pokt_v0

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"
)

func Test_findNodeOrError(t *testing.T) {
	type args struct {
		nodes  []*models.Node
		pubKey string
		err    error
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Node
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findNodeOrError(tt.args.nodes, tt.args.pubKey, tt.args.err)
			if !tt.wantErr(t, err, fmt.Sprintf("findNodeOrError(%v, %v, %v)", tt.args.nodes, tt.args.pubKey, tt.args.err)) {
				return
			}
			assert.Equalf(t, tt.want, got, "findNodeOrError(%v, %v, %v)", tt.args.nodes, tt.args.pubKey, tt.args.err)
		})
	}
}

func Test_getNodeFromRequest(t *testing.T) {
	type args struct {
		session            *models.Session
		selectedNodePubKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Node
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNodeFromRequest(tt.args.session, tt.args.selectedNodePubKey)
			if !tt.wantErr(t, err, fmt.Sprintf("getNodeFromRequest(%v, %v)", tt.args.session, tt.args.selectedNodePubKey)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getNodeFromRequest(%v, %v)", tt.args.session, tt.args.selectedNodePubKey)
		})
	}
}

func Test_getRandomNodeOrError(t *testing.T) {
	type args struct {
		nodes []*models.Node
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Node
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRandomNodeOrError(tt.args.nodes, tt.args.err)
			if !tt.wantErr(t, err, fmt.Sprintf("getRandomNodeOrError(%v, %v)", tt.args.nodes, tt.args.err)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getRandomNodeOrError(%v, %v)", tt.args.nodes, tt.args.err)
		})
	}
}
