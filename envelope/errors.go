package envelope

import "errors"

var (
	ErrInvalidChunkNumber  = errors.New("Invalid chunk number")
	ErrUninitialized       = errors.New("Chunk has not been initialized properly")
	ErrInvalidEnvelopeType = errors.New("Invalid envelope type")
)
