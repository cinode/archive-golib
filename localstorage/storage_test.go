package localstorage

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func genericStorageTest(s Storage, t *testing.T) {

	testContent, testBID := []byte("Hello world"), "hello"

	// Try to read blob that does not exist
	rdr, err := s.GetBlobReader("")
	if rdr != nil {
		t.Fatalf("Did receive reader although the blob id is invalid")
	}
	if err == nil {
		t.Fatalf("Did not get error although did not receive reader")
	}

	// Writer rollback
	wrt, err := s.GetBlobWriter()
	if err != nil {
		t.Fatalf("Unexpected error while creating new blob writer: %v", err)
	}
	if wrt == nil {
		t.Fatalf("Did not get proper writer object")
	}
	n, err := wrt.Write(testContent)
	if err != nil {
		t.Fatalf("Couldn't write to blob writer: %v", err)
	}
	if n != len(testContent) {
		t.Fatalf("Invalid number of bytes written, expected: %v, is: %v", len(testContent), n)
	}
	err = wrt.Rollback()
	if err != nil {
		t.Fatalf("Unexpected error while rolling back blob writer: %v", err)
	}
	n, err = wrt.Write(testContent)
	if err == nil {
		t.Fatalf("Could write to a blob that's been rolled back")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error received when writing to a rolled back blob: %v", err)
	}
	if n != 0 {
		t.Fatalf("Invalid number of bytes written to a rolled back blob: %v", n)
	}
	err = wrt.Rollback()
	if err == nil {
		t.Fatalf("Double rollback should not be allowed")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error for double rollback: %v", err)
	}
	err = wrt.Commit(testBID)
	if err == nil {
		t.Fatalf("Commit after rollback is not allowed")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error while commit-after-rollback: %v", err)
	}

	// Writer commit with invalid bid
	wrt, _ = s.GetBlobWriter()
	wrt.Write(testContent)
	err = wrt.Commit("")
	if err == nil {
		t.Fatal("Written blob with invalid empty blob id")
	}

	// Successfull write
	wrt, _ = s.GetBlobWriter()
	wrt.Write(testContent)
	err = wrt.Commit(testBID)
	if err != nil {
		t.Fatal("Couldn't commit blob: %v", err)
	}
	n, err = wrt.Write(testContent)
	if err == nil {
		t.Fatalf("Could write to a blob that's been committed")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error received when writing to a committed blob: %v", err)
	}
	if n != 0 {
		t.Fatalf("Invalid number of bytes written to a committed blob: %v", n)
	}
	err = wrt.Rollback()
	if err == nil {
		t.Fatalf("Rollback after commit should not be allowed")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error for rollback-after-commit: %v", err)
	}
	err = wrt.Commit(testBID)
	if err == nil {
		t.Fatalf("Double commits are not allowed")
	}
	if err != ErrBlobAlreadyFinalized {
		t.Fatalf("Invalid error while double-commit: %v", err)
	}

	// Trying to read the blob
	rdr, err = s.GetBlobReader(testBID)
	if err != nil {
		t.Fatalf("Couldn't open valid blob: %v", err)
	}
	if rdr == nil {
		t.Fatalf("Did not get a proper reader")
	}
	buff, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Fatalf("Couldn't read the contents of blob: %v", err)
	}
	if bytes.Compare(buff, testContent) != 0 {
		t.Fatalf("Invalid data read from the blob")
	}
}

func TestMemoryBlob(t *testing.T) {
	genericStorageTest(InMemory(), t)
}
