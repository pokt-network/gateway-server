package common

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test for Sha3_256Hash function in pkg/common/crypto.go file
func TestSha3_256Hash(t *testing.T) {

	type args struct {
		obj any
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Sha3_256_Payload1",
			args: args{
				obj: "test1",
			},
			want: []byte{161, 215, 105, 209, 47, 54, 86, 186, 167, 27, 175, 135, 81, 219, 180, 189, 25, 16, 73, 106, 74, 199, 100, 120, 143, 167, 64, 188, 208, 89, 118, 36},
		},
		{
			name: "Sha3_256_Payload2",
			args: args{
				obj: "test2",
			},
			want: []byte{239, 187, 198, 219, 95, 41, 245, 181, 210, 133, 45, 90, 15, 238, 124, 55, 162, 73, 233, 105, 12, 8, 178, 172, 35, 95, 179, 119, 38, 168, 182, 85},
		},
		{
			name: "Sha3_256_Payload3",
			args: args{
				obj: "", // empty string
			},
			want: []byte{40, 145, 227, 117, 2, 135, 48, 190, 87, 9, 94, 152, 9, 107, 184, 175, 66, 239, 2, 126, 70, 228, 39, 71, 139, 124, 167, 209, 51, 43, 216, 95},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Sha3_256Hash(tt.args.obj))
		})
	}

}

// test for Sha3_256HashHex function in pkg/common/crypto.go file
func TestSha3_256HashHex(t *testing.T) {

	type args struct {
		obj any
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Sha3_256_Payload1",
			args: args{
				obj: "test1",
			},
			want: "a1d769d12f3656baa71baf8751dbb4bd1910496a4ac764788fa740bcd0597624",
		},
		{
			name: "Sha3_256_Payload2",
			args: args{
				obj: "test2",
			},
			want: "efbbc6db5f29f5b5d2852d5a0fee7c37a249e9690c08b2ac235fb37726a8b655",
		},
		{
			name: "Sha3_256_Payload3",
			args: args{
				obj: "", // empty string
			},
			want: "2891e375028730be57095e98096bb8af42ef027e46e427478b7ca7d1332bd85f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Sha3_256HashHex(tt.args.obj))
		})
	}

}

// test for GetAddressFromPublicKey function in pkg/common/crypto.go file
func TestGetAddressFromPublicKey(t *testing.T) {

	type args struct {
		publicKey string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "invalid_public_key",
			args: args{
				publicKey: "--invalid-public-key--", // invalid public key
			},
			want:    "",
			wantErr: hex.InvalidByteError('-'),
		},
		{
			name: "valid_public_key_1",
			args: args{
				publicKey: "d992d8915443e85f620a67bce0928bba16cc349b6b6878698fd9518c6f49d5f2", // valid public key
			},
			want:    "f98920227f9a365320fcf2c8b750a8ab54dbf718", // address
			wantErr: nil,
		},
		{
			name: "valid_public_key_2",
			args: args{
				publicKey: "a59516f6f339699f6cfbe70046bc5cb5f1053cfee4c74e66be0ab1695e04a979", // valid public key
			},
			want:    "11188c2101fdc762e7c2ae7a77b819937809af53", // address
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetAddressFromPublicKey(tt.args.publicKey)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)

		})
	}

}
