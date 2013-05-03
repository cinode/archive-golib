package blobstore

import (
	"io"
)

type DirBlobReader interface {
	// Open blob for reading
	Open(bid, key string) error

	// Test whether there's any more entry left
	IsNextEntry() bool

	// Get the next entry from the reader
	NextEntry() (DirEntry, error)
}

type dirBlobReader struct {
	baseBlobReader             // Inherit methods of base blob reader
	Storage        BlobStorage // Blob storage
	currentReader  io.Reader   // Current reader we work on
	entriesLeft    int64       // Number of directory entries left to read   
}

func NewDirBlobReader(storage BlobStorage) DirBlobReader {
	return &dirBlobReader{
		baseBlobReader: baseBlobReader{
			storage: storage}}
}

func (d *dirBlobReader) Open(bid, key string) error {

	// Get the raw blob reader
	reader, blobType, err := d.openInternal(bid, key, validationMethodHash)
	if err != nil {
		return err
	}

	// Validate the blob type
	switch blobType {

	case blobTypeSimpleStaticDir:
		d.currentReader = reader
		if d.entriesLeft, err = deserializeInt(reader); err != nil {
			return err
		}
		if d.entriesLeft < 0 || d.entriesLeft > maxSimpleDirEntries {
			return ErrMalformedDirInvalidEntriesCount
		}
		if err = d.eofTest(); err != nil {
			return err
		}
		return nil

	case blobTypeSplitStaticDir:
		panic("Split directory blobs are unimplemented")
	}

	return ErrInvalidFileBlobType
}

func (d *dirBlobReader) IsNextEntry() bool {
	return d.entriesLeft > 0
}

func (d *dirBlobReader) NextEntry() (entry DirEntry, err error) {

	// Test if there's anything else to read
	if !d.IsNextEntry() {
		err = ErrNoMoreDirEntries
		return
	}

	// Make sure the nober of entries left decreases
	// even in case of an error
	d.entriesLeft--

	// Read one entry
	if err = entry.deserialize(d.currentReader); err != nil {
		return
	}

	err = nil
	return
}

func (d *dirBlobReader) eofTest() error {
	if !d.IsNextEntry() && !d.atEOF() {
		// There must be no more data if we're at the end
		// of data stream
		return ErrMalformedDirExtraData
	}
	return nil
}

func (d *dirBlobReader) atEOF() bool {
	// TODO: We're using this for validation only, implement the proper version
	return true
}
