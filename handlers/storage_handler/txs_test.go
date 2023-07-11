package storage_handler_test

import (
	"fmt"
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestBuyStorage(t *testing.T) {
	r := require.New(t)

	wallet, err := wallet_handler.NewWalletHandler(
		"slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum",
		"https://testnet-rpc.jackalprotocol.com:443",
		"lupulella-2")
	r.NoError(err)

	storageHandler := storage_handler.NewStorageHandler(wallet)

	res, err := storageHandler.BuyStorage("jkl1xzlwuc79dt4g2kxezwpqk5m6eh4wc0zwcpjsyf", 1, 1)
	r.NoError(err)

	fmt.Println(res.RawLog)

	r.Equal(uint32(0), res.Code)

}
