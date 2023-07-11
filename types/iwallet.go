package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type Wallet interface {
	GetAddress() string
	GetPubKey() types.PubKey
	GetClientCtx() client.Context
	GetChainID() string
	FindPubKey(address string) (types.PubKey, error)
	AsymmetricDecrypt(toDecrypt string) ([]byte, error)
	AsymmetricEncrypt(toEncrypt []byte, pubKey cryptotypes.PubKey) (string, error)
}
