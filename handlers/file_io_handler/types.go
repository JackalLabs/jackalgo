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

type FileIoHandler struct {
	walletHandler *wallet_handler.WalletHandler
	providers     []storagetypes.Providers
}

func NewFileIoHandler(w *wallet_handler.WalletHandler) *FileIoHandler {
	f := FileIoHandler{
		walletHandler: w,
	}

	return &f
}

func (f *FileIoHandler) tumbleUpload(sender string, file *types.File) (fid string, cid string, err error) {
	var randProvs []storagetypes.Providers
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

		fid, cid, err := doUpload(uploadUrl, sender, file)
		if err == nil {
			return fid, cid, err
		}
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
	// copy the file into the fileWriter
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url.String(), &b)
	if err != nil {
		return
	}

	cli := &http.Client{Timeout: time.Second * 100}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed with code: %d", res.StatusCode)
		return
	}

	var pup provider.UploadResponse
	err = json.NewDecoder(res.Body).Decode(&pup)
	if err != nil {
		return
	}

	return pup.FID, pup.CID, nil
}

type Queue struct {
	details []*file_upload_handler.FileUploadHandler
	count   int
}

func NewQueue() *Queue {
	q := Queue{
		details: make([]*file_upload_handler.FileUploadHandler, 0),
		count:   0,
	}
	return &q
}

func (q *Queue) Push(f *file_upload_handler.FileUploadHandler) {
	q.details[q.count] = f
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
	Fids []string `json:"fids"`
}
