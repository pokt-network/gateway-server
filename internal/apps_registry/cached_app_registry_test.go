package apps_registry

import (
	"github.com/pokt-network/gateway-server/internal/apps_registry/models"
	pokt_models "github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"testing"
)

func Test_arePoktApplicationSignersEqual(t *testing.T) {
	type args struct {
		slice1 []*models.PoktApplicationSigner
		slice2 []*models.PoktApplicationSigner
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "different length",
			args: args{
				slice1: []*models.PoktApplicationSigner{},
				slice2: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
			},
			want: false,
		},
		{
			name: "different address",
			args: args{
				slice1: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "1234",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
				slice2: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
			},
			want: false,
		},
		{
			name: "different chains",
			args: args{
				slice1: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "1234",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
				slice2: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
			},
			want: false,
		},
		{
			name: "same apps with exact chains",
			args: args{
				slice1: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
				slice2: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "123"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
			},
			want: true,
		},
		{
			name: "same apps with same chains different ordering",
			args: args{
				slice1: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "1234", "1235"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
				slice2: []*models.PoktApplicationSigner{{NetworkApp: &pokt_models.PoktApplication{
					Address:   "123",
					Chains:    []string{"123", "1235", "1234"},
					PublicKey: "",
					Status:    0,
					MaxRelays: 0,
				}}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arePoktApplicationSignersEqual(tt.args.slice1, tt.args.slice2); got != tt.want {
				t.Errorf("arePoktApplicationSignersEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
