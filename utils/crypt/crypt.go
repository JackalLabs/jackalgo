package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/JackalLabs/jackalgo/types"
)

func GenKey() []byte {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return token
}

func GenIv() []byte {
	token := make([]byte, 4)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return token
}

func ConvertToEncryptedFile(workingFile types.File, key []byte, iv []byte) (*types.File, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	b := bytes.NewBuffer([]byte{})
	_, err = b.ReadFrom(&workingFile)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, b.Bytes(), nil)
	//
	//const detailsBlob = new Blob([JSON.stringify(details)])
	//const encryptedArray: Blob[] = [
	//new Blob([(detailsBlob.size + 16).toString().padStart(8, '0')]),
	//await aesCrypt(detailsBlob, key, iv, 'encrypt')
	//]
	//for (let i = 0; i < workingFile.size; i += chunkSize) {
	//const blobChunk = workingFile.slice(i, i + chunkSize)
	//encryptedArray.push(
	//new Blob([(blobChunk.size + 16).toString().padStart(8, '0')]),
	//await aesCrypt(blobChunk, key, iv, 'encrypt')
	//)
	//}
	//const finalName = `${await hashAndHex(
	//  details.name + Date.now().toString()
	//)}.jkl`

	hexedName := HashAndHex(fmt.Sprintf("%s%d", workingFile.Name(), time.Now()))

	finalName := fmt.Sprintf("%s.jkl", hexedName)

	details := types.Details{
		Name:         finalName,
		Size:         int64(len(cipherText)),
		FileType:     "text/plain",
		LastModified: time.Now(),
	}

	f := types.File{
		Buffer:  bytes.NewBuffer(cipherText),
		Details: details,
	}

	return &f, nil
}

func HashAndHex(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	data := h.Sum(nil)

	return hex.EncodeToString(data)
}

func HexFullPath(path string, fileName string) string {
	return HashAndHex(fmt.Sprintf("%s%s", path, HashAndHex(fileName)))
}

func MerkleMeBro(rawpath string) string {
	pathArray := strings.Split(rawpath, "/")
	merkle := ""
	for i := 0; i < len(pathArray); i++ {
		merkle = HexFullPath(merkle, pathArray[i])
	}

	return merkle
}
