package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test for IsHttpOk function in pkg/common/http_test.go file
func TestIsHttpOk(t *testing.T) {

	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "Informational",
			statusCode: 100,
			want:       false,
		},
		{
			name:       "Successful",
			statusCode: 200,
			want:       true,
		},
		{
			name:       "Redirection",
			statusCode: 300,
			want:       false,
		},
		{
			name:       "ClientError",
			statusCode: 400,
			want:       false,
		},
		{
			name:       "ServerError",
			statusCode: 500,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsHttpOk(tt.statusCode))
		})
	}

}
