package blobstore

import (
	"bytes"
	"io"
)

type FileBlobWriter struct {

	// Buffer for storing data before we can hash it
	buffer bytes.Buffer

	// Storage object
	Storage BlobStorage

	// List of partial file blobs
	partialBids, partialKeys []string
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
				f.cleanup()
				return 0, err
			}
			bufferSpaceLeft = maxSimpleFileDataSize
		}
	}
	return written, nil
}

func (f *FileBlobWriter) finalizePartialBuffer() error {

	// Create the header
	var hdr bytes.Buffer
	hdr.WriteByte(blobTypeSimpleStaticFile)
	serializeInt(int64(f.buffer.Len()), &hdr)

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

	// Cleanup
	f.buffer.Reset()

	return nil
}

func (f *FileBlobWriter) addPartialBlob(bid, key string) {
	f.partialBids = append(f.partialBids, bid)
	f.partialKeys = append(f.partialKeys, key)
}

func (f *FileBlobWriter) Finalize() (bid string, key string, err error) {

	// Throw out the last partial if needed
	if f.buffer.Len() > 0 || len(f.partialBids) == 0 {
		if err := f.finalizePartialBuffer(); err != nil {
			f.cleanup()
			return "", "", err
		}
	}

	// If there's only one partial in the list, we don't have to create
	// any split file blobs
	if len(f.partialBids) == 1 {
		return f.partialBids[0], f.partialKeys[0], nil
	}

	// Create split file blob
	panic("Unimplemented")
}

func (f *FileBlobWriter) cleanup() {
	// TODO: Remove all blobs generated so far
}
