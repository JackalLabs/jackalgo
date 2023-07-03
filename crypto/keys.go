package crypto

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func GetAddress(ctx client.Context) (string, error) {
	key, err := ReadKey(ctx)
	if err != nil {
		return "", err
	}

	return key.Address, nil
}

func ReadKey(ctx client.Context) (*StorPrivKey, error) {

	keyStruct := StorPrivKey{} // TODO: Get Key from seed phrase
	return &keyStruct, nil
}

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
