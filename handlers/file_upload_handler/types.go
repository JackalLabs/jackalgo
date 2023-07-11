package file_upload_handler

import (
	"bytes"
	"os"

	"github.com/JackalLabs/jackalgo/types"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	"github.com/google/uuid"
)

type FileUploadHandler struct {
	File       types.File
	parentPath string
	uuid       string
	key        []byte
	iv         []byte
	cid        string
	fid        []string
}

func NewFileUploadHandler(file os.File, parentPath string, uuid string, savedKey []byte, savedIv []byte) (*FileUploadHandler, error) {
	fileDetails, err := file.Stat()
	if err != nil {
		return nil, err
	}

	details := types.Details{
		Name:         file.Name(),
		LastModified: fileDetails.ModTime(),
		FileType:     file.Name(),
		Size:         fileDetails.Size(),
	}

	b := bytes.NewBuffer([]byte{})
	_, err = b.ReadFrom(&file)
	if err != nil {
		return nil, err
	}

	newFile := types.File{
		Buffer:  b,
		Details: details,
	}

	f := FileUploadHandler{
		File:       newFile,
		parentPath: parentPath,
		uuid:       uuid,
		key:        savedKey,
		iv:         savedIv,
		cid:        "",
		fid:        make([]string, 0),
	}

	return &f, nil
}

func TrackFile(file os.File, parentPath string) (*FileUploadHandler, error) {
	savedKey := crypt.GenKey()
	savedIv := crypt.GenIv()
	uuid := uuid.New().String()

	return NewFileUploadHandler(file, parentPath, uuid, savedKey, savedIv)
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
		return &f.File, nil
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
