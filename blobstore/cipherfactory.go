package blobstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrInsufficientKeySource = errors.New("Not enough data to create a proper encryption key")
	ErrInvalidKey            = errors.New("Invalid key")
	ErrUnknownKeyType        = errors.New("Unknown key type")
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

func createDecryptor(key string, ivSource []byte, input io.Reader) (reader io.Reader, err error) {
	keyRaw, err := hex.DecodeString(key)
	if err != nil || len(keyRaw) < 1 {
		return nil, ErrInvalidKey
	}

	switch keyRaw[0] {
	case cipherAES256:
		return createDecryptorAES256(keyRaw[1:], ivSource, input)
	}

	return nil, ErrUnknownKeyType
}

func createDecryptorAES256(key []byte, ivSource []byte, input io.Reader) (reader io.Reader, err error) {

	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	// Normalize the iv
	var iv [16]byte
	for i, b := range ivSource {
		iv[i] = b
		if i >= 15 {
			break
		}
	}

	// Create new base cipher
	blobCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate the reader in CFB mode
	return &cipher.StreamReader{
			S: cipher.NewCFBDecrypter(
				blobCipher,
				iv[:]),
			R: input},
		nil
}

func createDataHash(data []byte) []byte {
	hasher := sha512.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
