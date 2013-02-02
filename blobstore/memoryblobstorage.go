package blobstore

// TODO: Support for duplicates (let write the blob with same id as long as the content does match)

import (
	"bytes"
	"io"
)

func NewMemoryBlobStorage() BlobStorage {
	return &memoryBlobStorage{
		blobs: make(map[string][]byte)}
}

type memoryBlobStorage struct {
	blobs map[string][]byte
}

type memoryBlobWriter struct {
	storage *memoryBlobStorage
	buffer  bytes.Buffer
	bid     string
}

func (f *memoryBlobWriter) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}

func (f *memoryBlobWriter) Finalize() error {
	previous, exists := f.storage.blobs[f.bid]
	if exists {
		if !bytes.Equal(previous, f.buffer.Bytes()) {
			return ErrBIDCollision
		}
	} else {
		f.storage.blobs[f.bid] = f.buffer.Bytes()
	}
	return nil
}

func (f *memoryBlobWriter) Cancel() error {
	f.buffer.Reset()
	f.storage, f.bid = nil, ""
	return nil
}

func (s *memoryBlobStorage) NewBlobWriter(blobId string) (writer BlobWriter, err error) {
	return &memoryBlobWriter{
			storage: s,
			bid:     blobId},
		nil
}

func (s *memoryBlobStorage) NewBlobReader(blobId string) (reader io.Reader, err error) {
	blob, ok := s.blobs[blobId]
	if !ok {
		return nil, ErrBIDNotFound
	}

	return bytes.NewReader(blob), nil
}
