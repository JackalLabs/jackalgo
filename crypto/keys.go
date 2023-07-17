package crypto

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func Sign(priv cryptotypes.PrivKey, msg []byte) ([]byte, error) {
	sig, err := priv.Sign(msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
