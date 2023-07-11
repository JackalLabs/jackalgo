package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/JackalLabs/jackalgo/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/libs/json"
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

func Encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, iv, data, nil)
	return cipherText, nil
}

func Decrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, iv, data, nil)
}

func ConvertFromEncryptedFile(data []byte, key []byte, iv []byte) (*types.File, error) {
	var details []byte
	parts := make([]byte, 0)
	var i int64
	for i = 0; i < int64(len(data)); {
		offset := i + 8
		segSize, err := strconv.ParseInt(string(data[i:offset]), 10, 64)
		if err != nil {
			return nil, err
		}
		last := offset + segSize
		segment := data[offset:last]

		raw, err := Decrypt(segment, key, iv)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			details = raw
		} else {
			parts = append(parts, raw...)
		}
		i = last
	}

	var detailStruct types.Details
	err := json.Unmarshal(details, &detailStruct)
	if err != nil {
		return nil, err
	}

	f := types.File{
		Buffer:  bytes.NewBuffer(parts),
		Details: detailStruct,
	}

	return &f, nil
}

func ConvertToEncryptedFile(workingFile types.File, key []byte, iv []byte) (*types.File, error) {
	chunkSize := int64(32 * 1024 * 1024)

	jsonDetails, err := json.Marshal(workingFile.Details) // TODO make sure details match json
	if err != nil {
		return nil, err
	}

	encryptedArray := []byte{}

	b, err := Encrypt(jsonDetails, key, iv)
	if err != nil {
		return nil, err
	}
	chunkedSize := int64(len(b) + 16)
	sizeData := []byte(fmt.Sprintf("%08d", chunkedSize))
	encryptedArray = append(encryptedArray, sizeData...)
	encryptedArray = append(encryptedArray, b...)

	fileBytes := workingFile.Buffer.Bytes()
	for i := int64(0); i < workingFile.Details.Size; i += chunkSize {
		end := i + chunkSize
		s := int64(len(fileBytes))
		if end >= s {
			end = s - 1
		}
		chunk := fileBytes[i:end]
		enc, err := Encrypt(chunk, key, iv)
		if err != nil {
			return nil, err
		}
		chunkedSize := int64(len(chunk) + 16)
		sizeData := []byte(fmt.Sprintf("%08d", chunkedSize))
		encryptedArray = append(encryptedArray, sizeData...)
		encryptedArray = append(encryptedArray, enc...)
	}

	hexedName := HashAndHex(fmt.Sprintf("%s%d", workingFile.Name(), time.Now().Unix()))

	finalName := fmt.Sprintf("%s.jkl", hexedName)

	details := types.Details{
		Name:         finalName,
		Size:         int64(len(encryptedArray)),
		FileType:     "text/plain",
		LastModified: time.Now(),
	}

	f := types.File{
		Buffer:  bytes.NewBuffer(encryptedArray),
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

func AesToString(wallet types.Wallet, pubKey cryptotypes.PubKey, key []byte, iv []byte) (string, error) {
	theIv, err := wallet.AsymmetricEncrypt(iv, pubKey)
	if err != nil {
		return "", err
	}
	theKey, err := wallet.AsymmetricEncrypt(key, pubKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s|%s", theIv, theKey), nil
}

func StringToAes(wallet types.Wallet, source string) (iv []byte, key []byte, err error) {
	if !strings.Contains(source, "|") {
		return nil, nil, fmt.Errorf("cannot have pipe before string start")
	}

	parts := strings.Split(source, "|")

	theIv, err := wallet.AsymmetricDecrypt(parts[0])
	if err != nil {
		return nil, nil, err
	}
	theKey, err := wallet.AsymmetricDecrypt(parts[1])
	if err != nil {
		return nil, nil, err
	}
	return theIv, theKey, nil
}
