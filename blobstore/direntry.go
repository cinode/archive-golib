package blobstore

import (
	"bytes"
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
