package envelope

import (
	"testing"
)

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
	if err != ErrInvalidChunkNumber {
		t.Fatalf("Expected error: %v, is: %v", ErrInvalidChunkNumber, err)
	}

	r, err = e.GetChunkReader(0)
	if r != nil {
		t.Fatalf("Non-null reader for empty hash envelope")
	}
	if err != ErrInvalidChunkNumber {
		t.Fatalf("Expected error: %v, is: %v", ErrInvalidChunkNumber, err)
	}

}
