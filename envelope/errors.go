package envelope

import "errors"

var ErrInvalidChunkNumber = errors.New("Invalid chunk number")
var ErrUninitialized = errors.New("Chunk has not been initialized properly")
