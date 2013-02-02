package blobstore

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strings"
	"testing"
)

var simpleTests = []struct {
	content string
	blobHex string
	key     string
	bid     string
}{
	{ // Empty file
		"",
		"01 eb",
		"01 7b54b668 36c1fbdd 13d2441d 9e1434dc 62ca677f b68f5fe6 6a464baa decdbd00",
		"b4f5a7bb 878c0cec 9cb4bd6a e8bb175a 7ea59c1a 048c5ab7 c119990d 0041cb9c fb67c2aa 9e6fada8 11271977 7b4b80ff ada80205 f8ebe698 1c0ade97 ff3df8e5"},
	{ // File with single 'a' character
		"a",
		"01 8f14",
		"01 504ce2f6 de7e3338 9deb73b2 1f765570 ad2b9f2a a8aaec83 28f47b48 bc3e841f",
		"c9d30a99 38ecea16 bed58efe 5ad5b998 927a56da 7c8c36c1 ee13292d ec79aa50 c5613fc9 0d80c37a 77a5a422 691d1967 693a1236 892e228a d95ed6fe 4b505d85"},
	{ // Programmer's challenge
		"Hello World!",
		"01 855e296f 95d1eaf3 feb7d48c e0",
		"01 ac9d2591 34ccef98 7f9f4df3 115b0b7a 24b379cb ebb2aaa9 1ed811c8 cf5e0907",
		"82aeef20 2165cf11 930ea44a 9ad8337a ea355d63 751a7260 552e3e01 4ad6313b ca69c83f a4e35555 31d44a10 25708183 784af0e2 002562b7 260559ce 0e7af262"},
	{ // Alphabet
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"01 f0ead942 12737b28 60ea35e3 1c7dd176 b5620968 2c3a6792 1d464823 13c245d4 551c765c 3ca851d7 f375911a 66e6b52b 650d51ea c3",
		"01 b11ef5de bd728940 485629e3 42c572bc c5b103d7 b56de27b 07f901b4 abcdb5d4",
		"4cfb056a 184d4377 eff9fc3e 8364906a f4b3b3c9 467c2fb8 245382bd d535ea17 f8a63abc 190a9253 9bd92951 52f112d3 365d4910 737b9f9f 3e0eb2f2 eef40648"},
}

func TestEmptyFile(t *testing.T) {

	for _, test := range simpleTests {

		key := strings.Replace(test.key, " ", "", -1)
		bid := strings.Replace(test.bid, " ", "", -1)
		blob, _ := hex.DecodeString(strings.Replace(test.blobHex, " ", "", -1))

		m := NewMemoryBlobStorage()
		bw := FileBlobWriter{Storage: m}
		bw.Write([]byte(test.content))

		rbid, rkey, err := bw.Finalize()
		if err != nil {
			t.Error(err)
		}

		if rbid != bid {
			t.Errorf("Invalid blob id generated, got: %v, expected: %v", rbid, bid)
		}

		if rkey != key {
			t.Errorf("Invalid key generated, got: %v, expected: %v", rkey, key)
		}

		reader, err := m.NewBlobReader(rbid)
		if err != nil {
			t.Errorf("Couldn't open the blob with id: %v for reading: %v", rbid, err)
		} else {

			readBytes, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Errorf("Couldn't read the blob with id: %v, error: %v", rbid, err)
			} else if !bytes.Equal(readBytes, blob) {
				t.Errorf("The blob with id: %v has invalid content", rbid)
			}
		}
	}
}
