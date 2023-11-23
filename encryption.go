package HangmanWeb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// Encrypt encrypts the data using the first line as an encryption key.
func Encrypt(data string) string {

	var quid, info string
	var quidDone bool

	for _, char := range data {
		if char == '\n' {
			quidDone = true
			continue
		}
		if quidDone {
			info += string(char)
		} else {
			quid += string(char)
		}
	}

	aliquid, _ := hex.DecodeString(quid)

	byteInfo := []byte(info)

	cipherBlock, err := aes.NewCipher(aliquid)
	if err != nil {
		panic(err.Error())
	}

	cipherInfo := make([]byte, aes.BlockSize+len(byteInfo))
	iv := cipherInfo[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	stream.XORKeyStream(cipherInfo[aes.BlockSize:], byteInfo)

	return base64.URLEncoding.EncodeToString(cipherInfo)
}

// Decrypt decrypts the data using the first line as a decryption key.
func Decrypt(data string) string {

	var quid, info string
	var quidDone bool

	for _, char := range data {
		if char == '\n' {
			quidDone = true
			continue
		}
		if quidDone {
			info += string(char)
		} else {
			quid += string(char)
		}
	}

	aliquid, _ := hex.DecodeString(quid)

	cipherInfo, _ := base64.URLEncoding.DecodeString(info)

	cipherBlock, err := aes.NewCipher(aliquid)
	if err != nil {
		panic(err)
	}

	if len(cipherInfo) < aes.BlockSize {
		panic("cipherInfo too short")
	}
	iv := cipherInfo[:aes.BlockSize]
	cipherInfo = cipherInfo[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)

	stream.XORKeyStream(cipherInfo, cipherInfo)
	fmt.Println(string(cipherInfo))
	return string(cipherInfo)
}
