package localstorage

import (
	"errors"
)

var (
	ErrNoSuchBlob           = errors.New("blob with such ID could not be found")
	ErrBlobAlreadyFinalized = errors.New("this blob has already been finalized")
	ErrInvalidBlobID        = errors.New("invalid blob ID")
)
