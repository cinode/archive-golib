package blobstore

// TODO: Error handling

import (
	"io"
	"os"
)

func NewFileBlobStorage(path string) BlobStorage {
	os.MkdirAll(path, 0777)
	return &fileBlobStorage{path: path}
}

type fileBlobStorage struct {
	path string
}

func (s *fileBlobStorage) NewBlobWriter(blobId string) io.WriteCloser {
	fl, _ := os.Create(s.path + string(os.PathSeparator) + blobId)
	return fl
}
