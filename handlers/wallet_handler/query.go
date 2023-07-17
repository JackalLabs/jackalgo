package wallet_handler

import (
	"context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	ecies "github.com/ecies/go/v2"
	filetreetypes "github.com/jackalLabs/canine-chain/v3/x/filetree/types"
)

func (w *WalletHandler) GetChainID() string {
	return w.clientCtx.ChainID
}

func (w *WalletHandler) GetAddress() string {
	return w.address
}

func (w *WalletHandler) GetPubKey() types.PubKey {
	return w.key.PubKey()
}

func (w *WalletHandler) GetECIESPubKey() *ecies.PublicKey {
	return w.eciesKey.PublicKey
}

func (w *WalletHandler) getPrivKey() types.PrivKey {
	return w.key
}

func (w *WalletHandler) GetClientCtx() client.Context {
	return w.clientCtx
}

func (w *WalletHandler) FindPubKey(address string) (*ecies.PublicKey, error) {
	cli := filetreetypes.NewQueryClient(w.clientCtx)

	req := filetreetypes.QueryPubkeyRequest{Address: address}

	res, err := cli.Pubkey(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	r := res.Pubkey.Key

	hexKey, err := hex.DecodeString(r)
	if err != nil {
		return nil, err
	}

	newPkey, err := ecies.NewPublicKeyFromBytes(hexKey)
	if err != nil {
		return nil, err
	}
	return newPkey, nil
}
