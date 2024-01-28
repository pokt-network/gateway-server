package common

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/crypto/sha3"
)

func Sha3_256Hash(obj any) []byte {
	jsonStr, _ := ffjson.Marshal(obj)
	sha256Hash := sha3.New256()
	sha256Hash.Write(jsonStr)
	return sha256Hash.Sum(nil)
}

func Sha3_256HashHex(obj any) string {
	return hex.EncodeToString(Sha3_256Hash(obj))
}

func GetAddressFromPublicKey(publicKey string) (string, error) {
	bytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	hasher := sha256.New()
	hasher.Write(bytes)
	hashBytes := hasher.Sum(nil)
	address := hex.EncodeToString(hashBytes)
	return address[:40], nil
}
