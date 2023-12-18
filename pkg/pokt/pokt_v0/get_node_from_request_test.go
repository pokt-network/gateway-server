package pokt_v0

import (
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNodeFromRequest(t *testing.T) {
	// Prepare a mock session with nodes
	mockNodes := []*models.Node{
		{PublicKey: "pubKey1"},
		{PublicKey: "pubKey2"},
		{PublicKey: "pubKey3"},
	}
	mockSession := &models.Session{Nodes: mockNodes}

	testCases := []struct {
		Name               string
		SelectedNodePubKey string
		ExpectedError      error
		ExpectedNode       *models.Node
		ExpectedRandom     bool
	}{
		{
			"Get random node if selectedNodePubKey is empty",
			"",
			nil,
			nil, // The expectation for the node can be adjusted based on the test case.
			true,
		},
		{
			"Get specific node by public key",
			"pubKey2",
			nil,
			&models.Node{PublicKey: "pubKey2"},
			false,
		},
		{
			"Error if selectedNodePubKey is not found",
			"nonexistentKey",
			models.ErrNodeNotFound,
			nil,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			node, err := getNodeFromRequest(mockSession, tc.SelectedNodePubKey)
			assert.Equal(t, tc.ExpectedError, err)
			if tc.ExpectedRandom {
				assert.Contains(t, mockNodes, node)
			} else {
				assert.Equal(t, tc.ExpectedNode, node)
			}
		})
	}
}

func TestGetRandomNodeOrError(t *testing.T) {
	mockNodes := []*models.Node{
		{PublicKey: "pubKey1"},
		{PublicKey: "pubKey2"},
		{PublicKey: "pubKey3"},
	}

	testCases := []struct {
		Name          string
		Nodes         []*models.Node
		ExpectedError error
		ExpectedNode  *models.Node
	}{
		{
			"Get random node successfully",
			mockNodes,
			nil,
			nil,
		},
		{
			"Error if node list is empty",
			[]*models.Node{},
			models.ErrSessionHasZeroNodes,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			node, err := getRandomNodeOrError(tc.Nodes, tc.ExpectedError)

			assert.Equal(t, tc.ExpectedError, err)
			if err == nil {
				// use contains since random node to prevent flakiness
				assert.Contains(t, tc.Nodes, node)
			}
		})
	}
}

func TestFindNodeOrError(t *testing.T) {
	mockNodes := []*models.Node{
		{PublicKey: "pubKey1"},
		{PublicKey: "pubKey2"},
		{PublicKey: "pubKey3"},
	}

	testCases := []struct {
		Name          string
		Nodes         []*models.Node
		PubKeyToFind  string
		ExpectedError error
		ExpectedNode  *models.Node
	}{
		{
			"Find node by public key successfully",
			mockNodes,
			"pubKey2",
			nil,
			&models.Node{PublicKey: "pubKey2"},
		},
		{
			"Error if node is not found",
			mockNodes,
			"nonexistentKey",
			models.ErrNodeNotFound,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			node, err := findNodeOrError(tc.Nodes, tc.PubKeyToFind, tc.ExpectedError)

			assert.Equal(t, tc.ExpectedError, err)

			assert.Equal(t, tc.ExpectedNode, node)
		})
	}
}
