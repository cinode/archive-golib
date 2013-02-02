package blobstore

// TODO: Support for duplicates (let write the blob with same id as long as the content does match)

import (
	"os"
	"io"
)

func NewFileBlobStorage(path string) BlobStorage {
	os.MkdirAll(path, 0777)
	return &fileBlobStorage{path: path}
}

type fileBlobStorage struct {
	path string
}

type fileBlobWriter struct {
	fl *os.File
}

func (f *fileBlobWriter) Write(p []byte) (n int, err error) {
	return f.fl.Write(p)
}

func (f *fileBlobWriter) Finalize() error {
	return f.fl.Close()
}

func (f *fileBlobWriter) Cancel() error {
	f.fl.Close()
	os.Remove(f.fl.Name())
	return nil
}

func (s *fileBlobStorage) blobPath(blobId string) string {
	return s.path + string(os.PathSeparator) + blobId
}

func (s *fileBlobStorage) NewBlobWriter(blobId string) (writer BlobWriter, err error) {
	fl, err := os.OpenFile(s.blobPath(blobId), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &fileBlobWriter{fl}, nil
}

func (s *fileBlobStorage) NewBlobReader(blobId string) (reader io.Reader, err error) {
	return os.OpenFile(s.blobPath(blobId), os.O_RDONLY, 0666)
}
