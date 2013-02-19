package blobstore

import (
	"errors"
	"io"
)

var (
	ErrFileNeedsHashValidationMethod    = errors.New("Invalid blob validation method used - the file blob requires a hash-based validation")
	ErrInvalidBlobType                  = errors.New("Invalid blob type - not a file blob")
	ErrInvalidSplitFileSize             = errors.New("Invalid size of a split file")
	ErrMalformedSplitFileSizePartsCount = errors.New("Invalid split file blob - number of partial blobs is incorrect")
	ErrMalformedSplitFileExtraData      = errors.New("Invalid split file blob - extra bytes found at the end of the blob")
	ErrMalformedSplitFileExtraDataPart  = errors.New("Invalid split file blob - extra bytes found at the end of the partial blob")
	ErrInvalidSubBlobType               = errors.New("Invalid sub blob type - not a file blob")
)

// FileBlobReader is a structure that can be used to easily read from file blobs
type FileBlobReader struct {
	Storage             BlobStorage // Storage object
	currentReader       io.Reader   // Reader object currently used
	isSplit             bool        // Flag indicating whether this is a split file
	totalSize           int64       // Total file size. Valid for split files only, used for validation purposes only
	thisBlobBytesLeft   int         // Number of bytes left to read from this particular blob
	otherBlobsBytesLeft int64       // Number of bytes left to read in all blobs but this particular one
	otherBlobsBidsLeft  []string    // Bids for blobs not yet read
	otherBlobsKeysLeft  []string    // Keys for blobs not yet read
}

// Open does open blob with given bid and key
func (f *FileBlobReader) Open(bid, key string) error {

	// Get the raw blob reader
	reader, blobType, err := f.openInternal(bid, key)
	if err != nil {
		return err
	}

	switch blobType {

	// For simple type blob just read the rest of the unencrypted content
	case blobTypeSimpleStaticFile:
		f.isSplit = false
		f.totalSize = -1
		f.currentReader = reader
		return nil

	// For split file blob we have to read all entries and queue them
	case blobTypeSplitStaticFile:
		return f.loadSplitFileData(reader)
	}

	return ErrInvalidBlobType
}

// Internal function, try to open a blob having it's bid and key,
// don't interpret anything but blob's type
func (f *FileBlobReader) openInternal(bid, key string) (reader io.Reader, blobType int64, err error) {

	// Get the raw blob reader
	if reader, err = f.Storage.NewBlobReader(bid); err != nil {
		return
	}

	// Find out the validation method
	validationMethod, err := deserializeInt(reader)
	if err != nil {
		return
	}

	// File blobs must use the hash-based validation
	// TODO: We may relax this if we start using links and decide to dereference links here
	if validationMethod != validationMethodHash {
		return nil, 0, ErrFileNeedsHashValidationMethod
	}

	// Get the unencrypted stream
	if reader, err = createReaderForHashBlobData(reader, key); err != nil {
		return
	}

	// See what type of a blob this is
	blobType, err = deserializeInt(reader)
	if err != nil {
		return
	}

	return
}

// Setup the reader for loading split file content
func (f *FileBlobReader) loadSplitFileData(masterBlobReader io.Reader) error {

	// Read the size
	totalSize, err := deserializeInt(masterBlobReader)
	if err != nil {
		return err
	}

	// Read all sub-blob entries
	subBlobsCnt, err := deserializeInt(masterBlobReader)
	if err != nil {
		return err
	}

	// Make sure the sub blobs count is sane value
	if (subBlobsCnt < 2) || (subBlobsCnt > maxSaneSplitFileParts) {
		return ErrMalformedSplitFileSizePartsCount
	}

	// We can validate the total file size, subBlobsCnt-1 blobs must be of size
	// maxSimpleFileDataSize and the last one must be of size in range 1..maxSimpleFileDataSize
	maxSize := subBlobsCnt * maxSimpleFileDataSize
	minSize := maxSize - maxSimpleFileDataSize + 1
	if (totalSize < minSize) || (totalSize > maxSize) {
		return ErrInvalidSplitFileSize
	}

	// Read all sub-blob entries
	var bids, keys []string
	for i := int64(0); i < subBlobsCnt; i++ {
		bid, err := deserializeString(masterBlobReader, maxSaneBidLength)
		if err != nil {
			return err
		}
		key, err := deserializeString(masterBlobReader, maxSaneKeyLength)
		if err != nil {
			return err
		}

		bids = append(bids, bid)
		keys = append(keys, key)
	}

	// We must have read everything from the split file blob by now
	if !f.atEOF(masterBlobReader) {
		return ErrMalformedSplitFileExtraData
	}

	// Fill in the data
	f.isSplit = true
	f.totalSize = totalSize
	f.thisBlobBytesLeft = 0
	f.otherBlobsBytesLeft = totalSize
	f.otherBlobsBidsLeft = bids
	f.otherBlobsKeysLeft = keys

	return nil
}

func (f *FileBlobReader) Read(p []byte) (n int, err error) {

	// Simple case for the non-split file
	if !f.isSplit {
		return f.currentReader.Read(p)
	}

	// Make sure to advance to next partial blob if the current one is exhausted
	if f.thisBlobBytesLeft <= 0 {
		if err = f.switchToNextPartialBlob(); err != nil {
			return
		}
	}

	// Reduce the number of bytes we will read at this call
	// to prevent crossing the one partial blob border
	if f.thisBlobBytesLeft < n {
		n = f.thisBlobBytesLeft
	}

	n, err = f.currentReader.Read(p)
	f.thisBlobBytesLeft -= n

	// Make sure not to throw any error between partial blobs switch
	if n > 0 {
		err = nil
	}

	return
}

func (f *FileBlobReader) switchToNextPartialBlob() error {

	// Return EOF if no more blobs left
	if len(f.otherBlobsBidsLeft) == 0 {
		return io.EOF
	}

	// Make sure partial blobs did not contain any extra data
	if f.currentReader != nil && !f.atEOF(f.currentReader) {
		return ErrMalformedSplitFileExtraDataPart
	}

	// Try to open the next blob
	reader, blobType, err := f.openInternal(f.otherBlobsBidsLeft[0], f.otherBlobsKeysLeft[0])
	if err != nil {
		return err
	}
	if blobType != blobTypeSimpleStaticFile {
		return ErrInvalidSubBlobType
	}

	// Update structures
	f.otherBlobsBidsLeft = f.otherBlobsBidsLeft[1:]
	f.otherBlobsKeysLeft = f.otherBlobsKeysLeft[1:]
	if f.otherBlobsBytesLeft > maxSimpleFileDataSize {
		f.thisBlobBytesLeft = maxSimpleFileDataSize
	} else {
		f.thisBlobBytesLeft = int(f.otherBlobsBytesLeft)
	}
	f.otherBlobsBytesLeft -= int64(f.thisBlobBytesLeft)
	f.currentReader = reader

	return nil
}

func (f *FileBlobReader) atEOF(r io.Reader) bool {
	// TODO: We're using this for validation only, implement the proper version
	return true
}
