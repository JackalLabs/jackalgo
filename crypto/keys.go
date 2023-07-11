package crypto

import (
	"encoding/hex"
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func Sign(priv *cryptotypes.PrivKey, msg []byte) ([]byte, error) {
	sig, err := priv.Sign(msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func ParsePrivKey(key string) (*cryptotypes.PrivKey, error) {
	keyData, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	k := cryptotypes.PrivKey{
		Key: keyData,
	}

	return &k, nil
}

func ExportPrivKey(priv *cryptotypes.PrivKey) string {
	return fmt.Sprintf("%x", priv.Key)
}
