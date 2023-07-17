package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	ecies "github.com/ecies/go/v2"
)

type Wallet interface {
	GetAddress() string
	GetPubKey() types.PubKey
	GetClientCtx() client.Context
	GetChainID() string
	FindPubKey(address string) (*ecies.PublicKey, error)
	AsymmetricDecrypt(toDecrypt string) ([]byte, error)
	AsymmetricEncrypt(toEncrypt []byte, pubKey *ecies.PublicKey) (string, error)
	GetECIESPubKey() *ecies.PublicKey
}
