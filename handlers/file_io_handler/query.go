package file_io_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
	"time"

	provider "github.com/JackalLabs/jackal-provider/jprov/types"
	"github.com/JackalLabs/jackalgo/handlers/file_download_handler"
	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils"
	"github.com/JackalLabs/jackalgo/utils/compression"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	"github.com/cosmos/cosmos-sdk/types/query"
	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
)

func (f *FileIoHandler) DownloadFolder(rawPath string) (folderHandler *folder_handler.FolderHandler, err error) {
	rawFolder, err := compression.ReadFileTreeEntry(f.walletHandler.GetAddress(), rawPath, f.walletHandler)
	if err != nil {
		return nil, err
	}

	var frame folder_handler.FolderFileFrame
	err = json.Unmarshal(rawFolder, &frame)
	if err != nil {
		return nil, err
	}

	return folder_handler.TrackFolder(frame, f.walletHandler), nil
}

func (f *FileIoHandler) DownloadFile(rawPath string) (fileHandler *file_download_handler.FileDownloadHandler, err error) {
	res, err := utils.GetFileTreeData(rawPath, f.walletHandler.GetAddress(), f.walletHandler)
	if err != nil {
		return nil, err
	}

	vacc := res.Files.ViewingAccess

	fmt.Println(vacc)

	var perms compression.EditorsViewers
	err = json.Unmarshal([]byte(vacc), &perms)
	if err != nil {
		return nil, err
	}

	user := crypt.HashAndHex(fmt.Sprintf("%s%s%s", "v", res.Files.TrackingNumber, f.walletHandler.GetAddress()))

	realPerms := perms[user]

	fmt.Println(realPerms)

	iv, key, err := crypt.StringToAes(f.walletHandler, realPerms)
	if err != nil {
		return nil, err
	}

	contents := res.Files.Contents

	var fids Fids
	err = json.Unmarshal([]byte(contents), &fids)
	if err != nil {
		return nil, err
	}

	if len(fids.Fids) == 0 {
		return nil, fmt.Errorf("not enough fids in the file")
	}

	fid := fids.Fids[0]

	queryClient := storagetypes.NewQueryClient(f.walletHandler.GetClientCtx())

	req := storagetypes.QueryFindFileRequest{
		Fid: fid,
	}

	urlRes, err := queryClient.FindFile(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	var ips []string

	err = json.Unmarshal([]byte(urlRes.ProviderIps), &ips)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		providerUrl, err := url2.Parse(ip)
		if err != nil {
			fmt.Println(err)
			continue
		}

		bytes, err := doDownload(providerUrl, fid)
		if err != nil {
			fmt.Println(err)
			continue
		}

		handler, err := file_download_handler.TrackFile(bytes, key, iv)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return handler, nil
	}

	return nil, fmt.Errorf("could not download file")
}

func (f *FileIoHandler) DownloadFileFromFid(fid string) (fileHandler *file_download_handler.FileDownloadHandler, err error) {
	queryClient := storagetypes.NewQueryClient(f.walletHandler.GetClientCtx())

	req := storagetypes.QueryFindFileRequest{
		Fid: fid,
	}

	urlRes, err := queryClient.FindFile(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	var ips []string

	err = json.Unmarshal([]byte(urlRes.ProviderIps), &ips)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		providerUrl, err := url2.Parse(ip)
		if err != nil {
			continue
		}

		bytes, err := doDownload(providerUrl, fid)
		if err != nil {
			continue
		}

		f := types.NewFile(bytes, types.Details{Name: fid, Size: int64(len(bytes)), LastModified: time.Now()})

		handler := file_download_handler.NewFileDownloadHandler(f)

		return handler, nil
	}

	return nil, fmt.Errorf("could not download file")
}

func fetchProviders(wallet *wallet_handler.WalletHandler) (provs []storagetypes.Providers, err error) {
	page := query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      500,
		CountTotal: false,
		Reverse:    false,
	}
	req := storagetypes.QueryAllProvidersRequest{
		Pagination: &page,
	}

	cli := storagetypes.NewQueryClient(wallet.GetClientCtx())

	res, err := cli.ProvidersAll(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res.Providers, nil
}

func sanitizeProviders(oldProvs []storagetypes.Providers, chainId string) (provs []storagetypes.Providers) {
	var wg sync.WaitGroup
	provs = make([]storagetypes.Providers, 0)
	for _, prov := range oldProvs {
		p := prov
		wg.Add(1)
		go func() {
			defer wg.Done()

			url := strings.Trim(p.Ip, "/")
			newUrl, err := url2.Parse(url)
			if err != nil {
				return
			}

			versionURL := newUrl.JoinPath("version")

			req, err := http.NewRequest("GET", versionURL.String(), nil)
			if err != nil {
				return
			}

			client := http.DefaultClient
			client.Timeout = 20 * time.Second

			res, err := client.Do(req)
			if err != nil {
				return
			}
			if res.StatusCode != 200 {
				return
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				return
			}

			var version provider.VersionResponse
			err = json.Unmarshal(b, &version)
			if err != nil {
				return
			}

			if version.ChainID != chainId {
				return
			}

			provs = append(provs, p)
		}()
	}

	wg.Wait()

	return
}
