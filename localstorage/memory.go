package localstorage

import (
	"bytes"
	"io/ioutil"
)

// memory represents simple blob storage holding data inside it's memory
type memory struct {
	blobs map[string][]byte
}

func (m *memory) GetBlobReader(blobID string) (reader Reader, err error) {
	buff, ok := m.blobs[blobID]
	if !ok {
		return nil, ErrNoSuchBlob
	}

	return ioutil.NopCloser(bytes.NewBuffer(buff)), nil
}

func (m *memory) GetBlobWriter() (writer Writer, err error) {
	return &memoryWriter{m: m}, nil
}

type memoryWriter struct {
	b bytes.Buffer // Buffer holding part of the data written so far
	m *memory      // Parent memory storage object
}

func (w *memoryWriter) Commit(blobID string) error {
	if w.m == nil {
		return ErrBlobAlreadyFinalized
	}
	defer func() { w.m = nil }()
	if blobID == "" {
		return ErrInvalidBlobID
	}
	w.m.blobs[blobID] = w.b.Bytes()
	return nil
}

func (w *memoryWriter) Rollback() error {
	if w.m == nil {
		return ErrBlobAlreadyFinalized
	}
	w.m = nil
	return nil
}

func (w *memoryWriter) Write(b []byte) (n int, err error) {
	if w.m == nil {
		return 0, ErrBlobAlreadyFinalized
	}
	return w.b.Write(b)
}

func InMemory() Storage {
	return &memory{make(map[string][]byte)}
}
