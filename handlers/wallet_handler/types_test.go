package wallet_handler_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	r := require.New(t)

	seed := "slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum"

	wallet, err := wallet_handler.NewWalletHandler(seed, "https://testnet.jackalprotocol.com:443", "lupulella-2")
	r.NoError(err)

	address := "jkl15cwg7teruwldgelxdg96g4cqxhsar7ye3zhv9f"

	r.Equal(wallet.GetAddress(), address)
}
