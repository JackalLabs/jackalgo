package tests

import (
	"fmt"
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestDashboardCompat(t *testing.T) {
	r := require.New(t)

	wallet, err := wallet_handler.NewWalletHandler(
		"slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum",
		"https://jackal-testnet-rpc.polkachu.com:443",
		"lupulella-2")
	r.NoError(err)

	fmt.Println(wallet.GetAddress())

	fileIO, err := file_io_handler.NewFileIoHandler(wallet.WithGas("500000"))
	r.NoError(err)

	folder, err := fileIO.DownloadFolder("s/Home", wallet.GetAddress())
	r.NoError(err)

	children := folder.GetChildFiles()
	fmt.Println(children)

	f, err := fileIO.DownloadFile("s/Home/test_data.txt", wallet.GetAddress())
	r.NoError(err)

	fmt.Println(f.File.Details)

	fmt.Println(f.GetFile().Buffer().String())
}
