package main

import (
	hangman "HangmanWeb"
	"encoding/hex"
	"fmt"
)

func main() {

	llave := hex.EncodeToString([]byte{0xb5, 0xd4, 0xe6, 0x8a, 0xa5, 0x3d, 0x54, 0x53, 0xc8, 0xd5, 0x77, 0x66, 0x31, 0xf5, 0x5, 0xf0, 0x99, 0xce, 0x5a, 0xc6, 0x10, 0x5e, 0xd8, 0xc6, 0xaf, 0x4a, 0xd5, 0xad, 0xc4, 0x47, 0x4e, 0xf8})
	fmt.Println()
	fmt.Println("Text to encrypt:")

	originalText := "Thorgan is the best ever!"
	fmt.Println(originalText)
	fmt.Println()

	fmt.Println("Encrypting.....")
	// encrypt value to base64
	cryptoText := hangman.Encrypt(llave + "\n" + originalText)
	fmt.Println(cryptoText)
	fmt.Println()

	fmt.Println("Decrypting.....")
	// encrypt base64 crypto to original value
	text := hangman.Decrypt(llave + "\n" + cryptoText)
	fmt.Println(text)
	fmt.Println()
}
