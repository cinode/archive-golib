package envelope

import (
	"io"
)

// Envelope is a type used to pass arbitrary data chunks between nodes
// The data is encrypted and all operations done to propagate them don't
// require knowledge about what's inside. Each blob though must fulfill
// set of validation rules which guarantees consistency of data inside
// the network. This also allows things such as validation of chunk author,
// removal of outdated chunkgs etc.
// The envelope is responsible for BID generation.
type Envelope interface {

	// Get envelope type
	GetType() int

	// Make sure the envelope is valid by analyzing the content in the
	// assigned storage, nil will be returned on success, error with
	// validation result will be returned on failure
	Validate() error

	// Get the BID for this envelope, this can be empty string in case
	// the blob has not been fully created yet (i.e. it may require
	// at least one chunk to be written)
	GetBID() string

	// Get number of blob chunks contained within that envelope
	GetChunksCount() int

	// Get chunk reader
	GetChunkReader(chunkNumber int) (reader io.Reader, err error)

	// Get writer to new chunk that will be appended to list of chunks
	//
	// Note that order of chunks is not specified, each envelope type can
	// decide by itself what does it mean to append new chunk to it's list
	//
	// Make sure to close the writer after writing all data. Closing the
	// writer does actually materialize the chunk and in some cases finalize
	// blob generation
	GetNewChunkWriter() (writer io.WriteCloser, err error)

	// Get named attribute, null will be returned for invalid name
	GetAttribute(name string) interface{}

	// Synchronize current blob with other node.
	//
	// The communicatino channel should be established between two
	// envelopes of exactly the same type. This method should be very
	// carefull though since the communication channel may be used to
	// perform various attacks on our node.
	SynchronizeWithOtherNode(channel io.ReadWriteCloser, active bool) error
}
