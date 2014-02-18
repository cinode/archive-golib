package cipherfactory

import (
	"hash"
	"io"
)

// Interface for the provider of encryption primitives used in Cinode
type Factory interface {

	// TODO: Shouldn't we operate on io.WriterCloser here ?
	// Create io.Writer to encrypt data writter and save to provided writer,
	// Parameters:
	//   keySource - byte blob used as source for the key computation
	//   ivSource  - byte blob used as source for the iv computation
	//   output    - writer used as the output for produced bytes
	// Returns:
	//   writer - writer where plain data should be written
	//   key    - key in string form that can be used to create decryptor
	//   err    - error
	CreateEncryptor(keySource, ivSource []byte, output io.Writer) (writer io.Writer, key string, err error)

	// Create a decryptor from key (returned from CreateEncryptor function) and iv source,
	// Similarly to CreateEncryptor, this function returns reader that can be used to read
	// plain data, one must provide source data reader that should allow reading encrypted data
	CreateDecryptor(key string, ivSource []byte, input io.Reader) (reader io.Reader, err error)

	// Create default hasher
	CreateHasher() (hasher hash.Hash, err error)
}
