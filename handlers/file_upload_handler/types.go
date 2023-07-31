package file_upload_handler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	"github.com/google/uuid"
)

const GoVirtual = "govirtual"

type FileUploadHandler struct {
	File       *types.File
	parentPath string
	uuid       string
	key        []byte
	iv         []byte
	cid        string
	fid        []string
	public     bool
}

func NewFileUploadHandler(file *os.File, parentPath string, uuid string, savedKey []byte, savedIv []byte, public bool) (*FileUploadHandler, error) {
	fileDetails, err := file.Stat()
	if err != nil {
		return nil, err
	}

	details := types.Details{
		Name:         file.Name(),
		LastModified: fileDetails.ModTime(),
		FileType:     GoVirtual,
		Size:         fileDetails.Size(),
	}
	var b bytes.Buffer
	size, err := io.Copy(&b, file)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Loading in %d bytes\n", size)

	newFile := types.NewFile(b.Bytes(), details)

	f := FileUploadHandler{
		File:       newFile,
		parentPath: parentPath,
		uuid:       uuid,
		key:        savedKey,
		iv:         savedIv,
		cid:        "",
		fid:        make([]string, 0),
		public:     public,
	}

	return &f, nil
}

func TrackFile(file *os.File, parentPath string, public bool) (*FileUploadHandler, error) {
	savedKey := crypt.GenKey()
	savedIv := crypt.GenIv()
	uuid := uuid.New().String()

	return NewFileUploadHandler(file, parentPath, uuid, savedKey, savedIv, public)
}

func TrackVirtualFile(bytes []byte, fileName string, parentPath string, public bool) (*FileUploadHandler, error) {
	savedKey := crypt.GenKey()
	savedIv := crypt.GenIv()
	uuid := uuid.New().String()

	details := types.Details{
		Name:         fileName,
		LastModified: time.Now(),
		FileType:     GoVirtual,
		Size:         int64(len(bytes)),
	}

	newFile := types.NewFile(bytes, details)

	f := FileUploadHandler{
		File:       newFile,
		parentPath: parentPath,
		uuid:       uuid,
		key:        savedKey,
		iv:         savedIv,
		cid:        "",
		fid:        make([]string, 0),
		public:     public,
	}

	return &f, nil
}

func (f *FileUploadHandler) SetIds(cid string, fid []string) {
	f.cid = cid
	f.fid = fid
}

func (f *FileUploadHandler) SetUUID(uuid string) {
	f.uuid = uuid
}

func (f *FileUploadHandler) GetIds() (string, []string) {
	return f.cid, f.fid
}

func (f *FileUploadHandler) GetUUID() string {
	return f.uuid
}

func (f *FileUploadHandler) GetWhoAmI() string {
	return f.File.Name()
}

func (f *FileUploadHandler) GetWhereAmI() string {
	return f.parentPath
}

func (f *FileUploadHandler) GetForUpload(public bool) (*types.File, error) {
	if public {
		return f.File, nil
	}

	return crypt.ConvertToEncryptedFile(f.File, f.key, f.iv)
}

func (f *FileUploadHandler) GetEnc() (key []byte, iv []byte) {
	return f.key, f.iv
}

func (f *FileUploadHandler) GetFullMerkle() string {
	return crypt.HexFullPath(f.GetMerklePath(), f.GetWhoAmI())
}

func (f *FileUploadHandler) GetMerklePath() string {
	return crypt.MerkleMeBro(f.parentPath)
}

func (f *FileUploadHandler) GetMeta() types.Details {
	return f.File.Details
}
