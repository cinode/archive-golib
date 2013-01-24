package blobstore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
)

func NewFileBlobWriter(st BlobStorage) io.WriteCloser {
	ret := &fileBlobWriter{
		keyHash: sha512.New(),
		storage: st,
	}

	byteBuff := []byte{0x01}
	ret.buff.Write(byteBuff)
	ret.keyHash.Write(byteBuff)

	return ret
}

type fileBlobWriter struct {
	buff    bytes.Buffer
	keyHash hash.Hash
	storage BlobStorage
}

func (w *fileBlobWriter) Write(b []byte) (n int, err error) {
	if w.buff.Len()+len(b) >= 16*1024*1024 {
		// TODO: Going into split file mode
	}

	n, err = w.buff.Write(b)
	w.keyHash.Write(b)

	return n, err
}

func (w *fileBlobWriter) Close() error {
	keyBytes := w.keyHash.Sum(nil)[:32]

	cph, _ := aes.NewCipher(keyBytes)
	stream := cipher.NewCFBEncrypter(cph, make([]byte, 16))

	encryptBuffer := &bytes.Buffer{}
	hash := sha512.New()

	writer := io.MultiWriter(hash, &cipher.StreamWriter{S: stream, W: encryptBuffer})
	writer.Write([]byte{0x01})

	io.Copy(writer, &w.buff)

	blobId := hex.EncodeToString(hash.Sum(nil))

	io.Copy(w.storage.NewBlobWriter(blobId), encryptBuffer)

	return nil
}
