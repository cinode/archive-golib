package blobstore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"io"
)

const validationMethodHash = 0x01

func createBlobFromReaderGenerator(readerGenerator func() io.Reader, storage BlobStorage) (bid string, key string, err error) {

	// Generate the key
	hasher := sha512.New()
	io.Copy(hasher, readerGenerator())
	keyRaw := hasher.Sum(nil)[:32]
	key = "AES:" + hex.EncodeToString(keyRaw)

	// Generate the encrypted content
	encryptedBuffer := bytes.Buffer{}
	blobCipher, _ := aes.NewCipher(keyRaw)
	io.Copy(
		&cipher.StreamWriter{
			S: cipher.NewCFBEncrypter(
				blobCipher,
				make([]byte, 16)),
			W: &encryptedBuffer},
		readerGenerator())

	// Generate blob id
	hasher.Reset()
	io.Copy(hasher, bytes.NewReader(encryptedBuffer.Bytes()))
	bid = hex.EncodeToString(hasher.Sum(nil))

	// Finally generate the blob itself
	blobWriter := storage.NewBlobWriter(bid)
	if _, err = blobWriter.Write([]byte{validationMethodHash}); err != nil {
		// TODO: Cleanup
		return "", "", err
	}
	if _, err = io.Copy(blobWriter, &encryptedBuffer); err != nil {
		// TODO: Cleanup
		return "", "", err
	}

	// Ok, we're done here
	return
}
