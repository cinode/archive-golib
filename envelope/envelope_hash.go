package envelope

import ()

type envelopeHash struct {
	bid string
}

// Get envelope type
func (e *envelopeHash) GetType() int {
	return TypeHash
}

// Make sure the envelope is valid by analyzing the content in the
// assigned storage, null will be returned on success, error with
// validation result will be returned on failure
func (e *envelopeHash) Validate() error {
}

// Get the BID for this envelope, this can be null in case
// the blob has not been fully created yet (i.e. it may require
// at least one chunk to be written)
func (e *envelopeHash) GetBID() string {
	return e.bid
}

// Get number of blob chunks contained within that envelope
func (e *envelopeHash) GetChunksCount() int {

	// Uninitialized case
	if e.bid == "" {
		return 0
	}

	// Initialized case
	return 1
}

// Get chunk reader
func (e *envelopeHash) GetChunkReader(chunkNumber int) (reader io.Reader, err error) {

	if i < 0 {
		return nil, ErrInvalidChunkNumber
	}
	
	if e.bid == "" {
		return nil, ErrUninitialized
	}

}

// Get writer to new chunk that will be appended to list of chunks
//
// Note that order of chunks is not specified, each envelope type can
// decide by itself what does it mean to append new chunk to it's list
//
// Make sure to close the writer after writing all data. Closing the
// writer does actually materialize the chunk and in some cases finalize
// blob generation
func (e *envelopeHash) GetNewChunkWriter() (writer io.WriteCloser, err error) {}

// Get named attribute, null will be returned for invalid name
func (e *envelopeHash) GetAttribute(name string) interface{} {}

// Synchronize current blob with other node.
//
// The communicatino channel should be established between two
// envelopes of exactly the same type. This method should be very
// carefull though since the communication channel may be used to
// perform various attacks on our node.
func (e *envelopeHash) SynchronizeWithOtherNode(channel io.ReadWriteCloser, active bool) error {}
