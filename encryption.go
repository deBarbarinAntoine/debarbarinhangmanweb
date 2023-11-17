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

	//Create a new Cipher Block from the key
	cipherBlock, err := aes.NewCipher(aliquid)
	if err != nil {
		panic(err.Error())
	}

	/* The IV needs to be unique, but not secure. Therefore it's common to
	   include it at the beginning of the ciphertext.  */
	cipherInfo := make([]byte, aes.BlockSize+len(byteInfo))
	iv := cipherInfo[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	stream.XORKeyStream(cipherInfo[aes.BlockSize:], byteInfo)

	// convert to base64
	return base64.URLEncoding.EncodeToString(cipherInfo)
}

// decrypt from base64 to decrypted string
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

	/* The IV needs to be unique, but not secure. Therefore it's common to
	   include it at the beginning of the ciphertext.  */
	if len(cipherInfo) < aes.BlockSize {
		panic("cipherInfo too short")
	}
	iv := cipherInfo[:aes.BlockSize]
	cipherInfo = cipherInfo[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherInfo, cipherInfo)
	fmt.Println(string(cipherInfo))
	return string(cipherInfo)
}
