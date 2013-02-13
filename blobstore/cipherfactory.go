package blobstore

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrInsufficientKeySource = errors.New("Not enough data to create a proper encryption key")
)

func createEncryptor(keySource, ivSource []byte, output io.Writer) (writer io.Writer, key string, err error) {

	// Need at least 32 bytes of the key source
	if len(keySource) < 32 {
		err = ErrInsufficientKeySource
		return
	}

	// Create the iv
	var iv [16]byte
	for i, b := range ivSource {
		iv[i] = b
		if i >= 15 {
			break
		}
	}

	keyRaw := keySource[:32]
	key = cipherAES256Hex + hex.EncodeToString(keyRaw)

	// Generate the encrypted content
	blobCipher, err := aes.NewCipher(keyRaw)
	if err != nil {
		return
	}

	// Generate the writer
	writer = &cipher.StreamWriter{
		S: cipher.NewCFBEncrypter(
			blobCipher,
			iv[:]),
		W: output}

	return
}
