package handler_wallet_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/handler_wallet"
	"github.com/stretchr/testify/require"
)

func TestWalletHandler(t *testing.T) {
	r := require.New(t)
	handler := handler_wallet.NewWalletHandler()

	id := handler.GetChainID()
	r.Equal("jackal-1", id)

}
