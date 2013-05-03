package blobstore

import (
	"bytes"
	"io"
)

// Helper structure for holding one directory entry
type DirEntry struct {
	Name, MimeType, Bid, Key string
}

func (d *DirEntry) serialize(b *bytes.Buffer) {
	serializeString(d.Name, b)
	serializeString(d.MimeType, b)
	serializeString(d.Bid, b)
	serializeString(d.Key, b)
}

func (d *DirEntry) deserialize(r io.Reader) (err error) {
	if d.Name, err = deserializeString(r, maxSaneNameLenght); err != nil {
		return
	}
	if d.MimeType, err = deserializeString(r, maxSaneMimeTypeLength); err != nil {
		return
	}
	if d.Bid, err = deserializeString(r, maxSaneBidLength); err != nil {
		return
	}
	if d.Key, err = deserializeString(r, maxSaneKeyLength); err != nil {
		return
	}
	return nil
}
