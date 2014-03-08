package envelope

import (
	"github.com/cinode/golib/cipherfactory"
	"github.com/cinode/golib/localstorage"
	"testing"
)

func noError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func anyError(err error, t *testing.T) {
	if err == nil {
		t.Fatal("Expected error not found")
	}
}

func needError(err error, errRequired error, t *testing.T) {
	if err == nil {
		t.Fatal("Expected error not found, should be: %v", errRequired)
	}
	if err != errRequired {
		t.Fatal("Invalid error received, expected: %v, got: %v", errRequired, err)
	}
}

func TestHashTrivialStuff(t *testing.T) {

	e := &envelopeHash{}

	if tp := e.GetType(); tp != TypeHash {
		t.Fatalf("Expected type: %v, got: %v", TypeHash, tp)
	}

	if bid := e.GetBID(); bid != "" {
		t.Fatalf("BID of hash envelope without chunks must be empty but is: %v", bid)
	}

	if cnt := e.GetChunksCount(); cnt != 0 {
		t.Fatalf("Expected 0 chunks, got %v", cnt)
	}

	r, err := e.GetChunkReader(-123)
	if r != nil {
		t.Fatalf("Non-null reader for negatively-indexed chunk")
	}
	needError(err, ErrInvalidChunkNumber, t)

	r, err = e.GetChunkReader(0)
	if r != nil {
		t.Fatalf("Non-null reader for empty hash envelope")
	}
	needError(err, ErrInvalidChunkNumber, t)

}

var testVector = []struct {
	bid  string
	blob []byte
}{
	{
		bid:  "b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
		blob: []byte{0x01, 0xeb},
	},
	{
		bid:  "c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
		blob: []byte{0x01, 0x8f, 0x14},
	},
	{
		bid:  "82aeef202165cf11930ea44a9ad8337aea355d63751a7260552e3e014ad6313bca69c83fa4e3555531d44a1025708183784af0e2002562b7260559ce0e7af262",
		blob: []byte{0x01, 0x85, 0x5e, 0x29, 0x6f, 0x95, 0xd1, 0xea, 0xf3, 0xfe, 0xb7, 0xd4, 0x8c, 0xe0},
	},
}

func TestHashValidation(t *testing.T) {
	s := localstorage.InMemory()
	cf := cipherfactory.Create()

	// Fill blobstore with test blobs
	for _, tv := range testVector {

		w, err := s.GetBlobWriter()
		noError(err, t)

		_, err = w.Write(tv.blob)
		noError(err, t)

		err = w.Commit(tv.bid)
		noError(err, t)
	}

	// Validate blobs
	for _, tv := range testVector {
		e := &envelopeHash{bid: tv.bid, storage: s, cf: cf}
		err := e.Validate()
		noError(err, t)
	}
}
