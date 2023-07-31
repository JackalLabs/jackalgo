package file_io_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"strings"
	"time"

	provider "github.com/JackalLabs/jackal-provider/jprov/types"
	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/types"
	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
)

const JackalGoVersion = 1

type FileIoHandler struct {
	walletHandler *wallet_handler.WalletHandler
	providers     []storagetypes.Providers
}

func NewFileIoHandler(w *wallet_handler.WalletHandler) (*FileIoHandler, error) {
	provs, err := fetchProviders(w)
	if err != nil {
		return nil, err
	}

	provs = sanitizeProviders(provs, w.GetChainID())

	f := FileIoHandler{
		walletHandler: w,
		providers:     provs,
	}

	return &f, nil
}

func (f *FileIoHandler) tumbleUpload(sender string, file *types.File) (fid string, cid string, err error) {
	randProvs := make([]storagetypes.Providers, len(f.providers))
	copy(randProvs, f.providers)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(randProvs), func(i, j int) { randProvs[i], randProvs[j] = randProvs[j], randProvs[i] })

	for _, prov := range randProvs {
		url := strings.Trim(prov.Ip, "/")
		newUrl, err := url2.Parse(url)
		if err != nil {
			continue
		}

		uploadUrl := newUrl.JoinPath("upload")

		fmt.Printf("Doing upload to %s...\n", uploadUrl)

		fid, cid, err := doUpload(uploadUrl, sender, file)
		if err == nil || len(fid) == 0 {
			return fid, cid, nil
		}
		fmt.Println(err)
		fmt.Printf("Failed upload to %s.\n", uploadUrl)

	}

	return "", "", fmt.Errorf("failed to upload to any providers")
}

func doUpload(url *url2.URL, sender string, file *types.File) (fid string, cid string, err error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	defer writer.Close()

	err = writer.WriteField("sender", sender)
	if err != nil {
		return
	}

	fileWriter, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return
	}

	size, err := io.Copy(fileWriter, file.Buffer())
	if err != nil {
		return
	}
	fmt.Printf("Posting file of size: %d\n", size)
	err = writer.Close()
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url.String(), &b)
	if err != nil {
		return
	}

	cli := &http.Client{Timeout: 60 * time.Second}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := cli.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		var errRes provider.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			return
		}
		err = fmt.Errorf("code: %d -> %s", res.StatusCode, errRes.Error)
		return
	}

	var pup provider.UploadResponse
	err = json.NewDecoder(res.Body).Decode(&pup)
	if err != nil {
		return
	}

	err = res.Body.Close()
	if err != nil {
		return
	}

	return pup.FID, pup.CID, nil
}

func doDownload(url *url2.URL, fid string) ([]byte, error) {
	newUrl := url.JoinPath("download", fid)

	fmt.Printf("downloading file from %s\n", newUrl.String())

	req, err := http.NewRequest("GET", newUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	cli := &http.Client{Timeout: 0}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var errRes provider.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			return nil, err
		}
		err = fmt.Errorf("code: %d -> %s", res.StatusCode, errRes.Error)
		return nil, err
	}

	var n bytes.Buffer
	_, err = io.Copy(&n, res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return n.Bytes(), nil
}

type Queue struct {
	details []*file_upload_handler.FileUploadHandler
}

func NewQueue() *Queue {
	q := Queue{
		details: make([]*file_upload_handler.FileUploadHandler, 0),
	}
	return &q
}

func (q *Queue) Push(f *file_upload_handler.FileUploadHandler) {
	q.details = append(q.details, f)
}

func (q *Queue) Pop() *file_upload_handler.FileUploadHandler {
	details := q.details[0]
	q.details = q.details[1:]

	return details
}

func (q *Queue) Empty() bool {
	return len(q.details) == 0
}

type Fids struct {
	Fids            []string `json:"fids"`
	JackalGoVersion int      `json:"jackal_go_version"`
}
