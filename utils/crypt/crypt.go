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
	ecies "github.com/ecies/go/v2"
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
	token := make([]byte, aes.BlockSize)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return token
}

func Encrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	bPlaintext := PKCS5Padding(plaintext, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return ciphertext, nil
}

func Decrypt(cipherText []byte, encKey []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	bPlaintext := make([]byte, len(cipherText))
	mode.CryptBlocks(bPlaintext, cipherText)
	plainText := RemovePKCS5Padding(bPlaintext, aes.BlockSize)
	return plainText, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func RemovePKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := int(ciphertext[len(ciphertext)-1])
	if padding > blockSize {
		return ciphertext
	}
	undoPadding := len(ciphertext) - padding
	return ciphertext[:undoPadding]
}

func ConvertFromEncryptedFile(data []byte, key []byte, iv []byte) (*types.File, error) {
	var details []byte
	parts := make([]byte, 0)
	for len(data) > 0 {
		counterBytes := data[:8]
		segSize, err := strconv.ParseInt(string(counterBytes), 10, 64)
		if err != nil {
			return nil, err
		}
		segment := data[8 : segSize+8]

		raw, err := Decrypt(segment, key, iv)
		if err != nil {
			return nil, err
		}
		if len(details) == 0 {
			details = raw
		} else {
			parts = append(parts, raw...)
		}

		data = data[segSize+8:]
	}

	fmt.Println(string(details))

	var detailStruct types.Details
	err := json.Unmarshal(details, &detailStruct)
	if err != nil {
		return nil, err
	}

	f := types.NewFile(parts, detailStruct)

	return f, nil
}

func ConvertToEncryptedFile(workingFile *types.File, key []byte, iv []byte) (*types.File, error) {
	chunkSize := int64(32 * 1024 * 1024)

	jsonDetails, err := json.Marshal(workingFile.Details) // TODO make sure details match json
	if err != nil {
		return nil, err
	}

	var encryptedArray []byte

	b, err := Encrypt(jsonDetails, key, iv)
	if err != nil {
		return nil, err
	}
	chunkedSize := int64(len(b))
	sizeData := []byte(fmt.Sprintf("%08d", chunkedSize))
	encryptedArray = append(encryptedArray, sizeData...)
	encryptedArray = append(encryptedArray, b...)

	fileBytes := workingFile.Buffer().Bytes()
	for len(fileBytes) > 0 {
		l := int64(len(fileBytes))
		cSize := chunkSize
		if cSize > l {
			cSize = l
		}

		chunk := fileBytes[:cSize]
		enc, err := Encrypt(chunk, key, iv)
		if err != nil {
			return nil, err
		}
		chunkedSize := int64(len(enc))
		sizeData := []byte(fmt.Sprintf("%08d", chunkedSize))
		encryptedArray = append(encryptedArray, sizeData...)
		encryptedArray = append(encryptedArray, enc...)

		fileBytes = fileBytes[cSize:]
	}

	hexedName := HashAndHex(fmt.Sprintf("%s%d", workingFile.Name(), time.Now().Unix()))

	finalName := fmt.Sprintf("%s.jkl", hexedName)

	details := types.Details{
		Name:         finalName,
		Size:         int64(len(encryptedArray)),
		FileType:     "text/plain",
		LastModified: time.Now(),
	}

	f := types.NewFile(encryptedArray, details)

	return f, nil
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

func AesToString(wallet types.Wallet, pubKey *ecies.PublicKey, key []byte, iv []byte) (string, error) {
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
