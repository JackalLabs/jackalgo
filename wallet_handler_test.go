package jackalgo_test

import (
	"testing"

	"github.com/JackalLabs/jackalgo"
	"github.com/stretchr/testify/require"
)

func TestWalletHandler(t *testing.T) {
	r := require.New(t)
	handler := jackalgo.NewWalletHandler()

	id := handler.GetChainID()
	r.Equal("jackal-1", id)

}
