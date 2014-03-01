package cipherfactory

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

var (
	ErrInsufficientKeySource = errors.New("Not enough data to create a proper encryption key")
	ErrInvalidKey            = errors.New("Invalid key")
	ErrUnknownKeyType        = errors.New("Unknown key type")
)

type defaultFactory struct {
}

func (d *defaultFactory) GetMinKeySourceBytes() int {
	return cipherAES256KeySourceLength
}

func (d *defaultFactory) CreateEncryptor(keySource, ivSource []byte, output io.Writer) (writer io.Writer, key string, err error) {

	// Need at least 32 bytes of the key source
	if len(keySource) < cipherAES256KeySourceLength {
		err = ErrInsufficientKeySource
		return
	}

	// Create the iv
	var iv [aes.BlockSize]byte
	copy(iv[:], ivSource)

	// Create AES-compatible key
	keyRaw := keySource[:cipherAES256KeySourceLength]
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

func (d *defaultFactory) CreateDecryptor(key string, ivSource []byte, input io.Reader) (reader io.Reader, err error) {
	keyRaw, err := hex.DecodeString(key)
	if err != nil || len(keyRaw) < 1 {
		return nil, ErrInvalidKey
	}

	switch keyRaw[0] {
	case cipherAES256:
		return d.createDecryptorAES256(keyRaw[1:], ivSource, input)
	}

	return nil, ErrUnknownKeyType
}

func (d *defaultFactory) createDecryptorAES256(key []byte, ivSource []byte, input io.Reader) (reader io.Reader, err error) {

	if len(key) != cipherAES256KeySourceLength {
		return nil, ErrInvalidKey
	}

	// Normalize the iv
	var iv [aes.BlockSize]byte
	copy( iv[:], ivSource );

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

func (d *defaultFactory) CreateHasher() (hasher hash.Hash, err error) {
	return sha512.New(), nil
}
