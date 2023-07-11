package storage_handler_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestQueryPayData(t *testing.T) {
	r := require.New(t)

	wallet, err := wallet_handler.DefaultWalletHandler("slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum")
	r.NoError(err)

	storageHandler := storage_handler.NewStorageHandler(wallet)

	res, err := storageHandler.QueryGetPayData("jkl102tsq6skvmr9d4qp06p09r02d36jxettdl7qfy")
	r.NoError(err)

	r.Equal(int64(3000000000000), res.Bytes)
}
