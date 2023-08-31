package tests

import (
	"fmt"
	"github.com/JackalLabs/jackalgo/handlers/storage_handler"
	"os"
	"testing"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/stretchr/testify/require"
)

func TestFileUpload(t *testing.T) {
	r := require.New(t)

	wallet, err := wallet_handler.NewWalletHandler(
		"slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum",
		"https://jackal-testnet-rpc.polkachu.com:443",
		"lupulella-2")
	r.NoError(err)

	wallet = wallet.WithGas("500000")

	s := storage_handler.NewStorageHandler(wallet)
	res, err := s.BuyStorage(wallet.GetAddress(), 720, 1)
	r.NoError(err)
	fmt.Println(res.RawLog)

	fmt.Println(wallet.GetAddress())

	fileIO, err := file_io_handler.NewFileIoHandler(wallet)
	r.NoError(err)

	fileData, err := os.Open("test_data.txt")
	r.NoError(err)

	res, err = fileIO.GenerateInitialDirs([]string{"jackalgo"})
	r.NoError(err)

	r.Equal(uint32(0), res.Code)

	folder, err := fileIO.DownloadFolder("s/jackalgo")
	r.NoError(err)

	file, err := file_upload_handler.TrackFile(fileData, "s/jackalgo")
	r.NoError(err)

	r.Equal("test_data.txt", file.GetWhoAmI())

	failed, fids, cids, err := fileIO.StaggeredUploadFiles([]*file_upload_handler.FileUploadHandler{file}, folder, false)
	r.NoError(err)

	fmt.Println(fids)
	fmt.Println(cids)

	r.Equal(0, failed)

	folder, err = fileIO.DownloadFolder("s/jackalgo")
	r.NoError(err)

	children := folder.GetChildFiles()
	fmt.Println(children)

	f, err := fileIO.DownloadFile("s/jackalgo/test_data.txt")
	r.NoError(err)

	fmt.Println(f.File.Details)

	fmt.Println(f.GetFile().Buffer().String())

	err = fileIO.DeleteTargets([]string{"test_data.txt"}, folder)
	r.NoError(err)

	_, err = fileIO.DownloadFile("s/jackalgo/test_data.txt")
	r.Error(err)

	folder, err = fileIO.DownloadFolder("s/jackalgo")
	r.NoError(err)

	fmt.Println(folder.GetChildFiles())
	fmt.Println(folder.GetChildDirs())
}
