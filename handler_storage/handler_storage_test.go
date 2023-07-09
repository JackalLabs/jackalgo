package handler_storage_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handler_storage"
	"github.com/JackalLabs/jackalgo/handler_wallet"
	"github.com/stretchr/testify/require"
)

func TestBuyStorage(t *testing.T) {
	r := require.New(t)

	wallet := handler_wallet.NewWalletHandler()

	storageHandler := handler_storage.NewStorageHandler(wallet)

	_, err := storageHandler.BuyStorage("jkl1arsaayyj5tash86mwqudmcs2fd5jt5zgc3sexc", 1, 1_000_000)
	r.NoError(err)

}
