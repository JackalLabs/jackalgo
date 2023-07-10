package storage_handler_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestQueryPayData(t *testing.T) {
	r := require.New(t)

	wallet := wallet_handler.NewWalletHandler("slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum")

	storageHandler := storage_handler.NewStorageHandler(wallet)

	_, err := storageHandler.QueryGetPayData("jkl1arsaayyj5tash86mwqudmcs2fd5jt5zgc3sexc")
	r.NoError(err)
}
