package handler_storage_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/handler_storage"
	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
	"github.com/stretchr/testify/require"
)

func TestQueryPayData(t *testing.T) {
	r := require.New(t)

	wallet := handler_wallet.NewWalletHandler()

	storageHandler := handler_storage.NewStorageHandler(wallet)

	_, err := storageHandler.QueryGetPayData("jkl1arsaayyj5tash86mwqudmcs2fd5jt5zgc3sexc")
	r.NoError(err)

}
