package blobstore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"io"
)

func createHashValidatedBlobFromReaderGenerator(readerGenerator func() io.Reader, storage BlobStorage) (bid string, key string, err error) {

	// Generate the key
	hasher := sha512.New()
	io.Copy(hasher, readerGenerator())
	keyRaw := hasher.Sum(nil)[:32]
	key = cipherAES256Hex + hex.EncodeToString(keyRaw)

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
	blobWriter, err := storage.NewBlobWriter(bid)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			blobWriter.Cancel()
		}
	}()
	if _, err = blobWriter.Write([]byte{validationMethodHash}); err != nil {
		return
	}
	if _, err = io.Copy(blobWriter, &encryptedBuffer); err != nil {
		return
	}
	if err = blobWriter.Finalize(); err != nil {
		return
	}

	// Ok, we're done here
	return
}
