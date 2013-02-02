package blobstore

import (
	"bytes"
	"io"
)

// Structure used to generate static file blobs
type FileBlobWriter struct {

	// Buffer for storing data before we can hash it
	buffer bytes.Buffer

	// Storage object
	Storage BlobStorage

	// List of partial file blobs
	partialBids, partialKeys []string

	// Overall number of bytes written so far
	totalBytes int64
}

// Performing a write operation on the file blob
func (f *FileBlobWriter) Write(p []byte) (n int, err error) {

	bufferSpaceLeft := maxSimpleFileDataSize - f.buffer.Len()
	written := 0
	for len(p) > 0 {

		// Let's see how much can we chop this time
		partialSize := len(p)
		if partialSize > bufferSpaceLeft {
			partialSize = bufferSpaceLeft
		}

		// Chop off the next part
		f.buffer.Write(p[:partialSize])
		p = p[partialSize:]
		bufferSpaceLeft -= partialSize
		written += partialSize

		// Check out if we should emit next partial buffer
		if bufferSpaceLeft <= 0 {
			if err := f.finalizePartialBuffer(); err != nil {
				f.Cancel()
				return 0, err
			}
			bufferSpaceLeft = maxSimpleFileDataSize
		}
	}
	return written, nil
}

// Write the current content of internal buffer into a blob,
// save it's id and key in a list of partial blobs
func (f *FileBlobWriter) finalizePartialBuffer() error {

	// Create the header
	var hdr bytes.Buffer
	hdr.WriteByte(blobTypeSimpleStaticFile)

	// Generate the blob
	readerGen := func() io.Reader {
		headerReader := bytes.NewReader(hdr.Bytes())
		contentReader := bytes.NewReader(f.buffer.Bytes())
		return io.MultiReader(headerReader, contentReader)
	}
	bid, key, err := createHashValidatedBlobFromReaderGenerator(readerGen, f.Storage)
	if err != nil {
		return err
	}

	// Queue the blob on a list of partial blobs
	f.addPartialBlob(bid, key)

	// Increase the counter of bytes thrown out so far
	f.totalBytes += int64(f.buffer.Len())

	// Cleanup
	f.buffer.Reset()

	return nil
}

// Save bid and key into a list of partial blobs
func (f *FileBlobWriter) addPartialBlob(bid, key string) {
	f.partialBids = append(f.partialBids, bid)
	f.partialKeys = append(f.partialKeys, key)
}

// Finalize the generation of this file blob
func (f *FileBlobWriter) Finalize() (bid string, key string, err error) {

	// Throw out the last partial if needed
	if f.buffer.Len() > 0 || len(f.partialBids) == 0 {
		if err := f.finalizePartialBuffer(); err != nil {
			f.Cancel()
			return "", "", err
		}
	}

	// If there's only one partial in the list, we don't have to create
	// any split file blobs
	if len(f.partialBids) == 1 {
		return f.partialBids[0], f.partialKeys[0], nil
	}

	// Create split file blob
	return f.finalizeSplitFile()
}

// Finalize blob generation in case we've created split file blob
func (f *FileBlobWriter) finalizeSplitFile() (bid string, key string, err error) {
	var b bytes.Buffer

	// Blob type id
	b.WriteByte(blobTypeSplitStaticFile)

	// Total file size
	serializeInt(f.totalBytes, &b)

	// Number of partial blobs
	serializeInt(int64(len(f.partialBids)), &b)

	// Partial blobs list
	for i, bid := range f.partialBids {
		serializeString(bid, &b)
		serializeString(f.partialKeys[i], &b)
	}

	return createHashValidatedBlobFromReaderGenerator(
		func() io.Reader { return bytes.NewReader(b.Bytes()) },
		f.Storage)
}

// Cancel the generation of file blob, remove all blobs generated so far
func (f *FileBlobWriter) Cancel() {
	// TODO: Implement this
}
