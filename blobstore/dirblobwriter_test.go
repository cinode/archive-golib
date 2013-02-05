package blobstore

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strings"
	"testing"
)

type testDirEntry struct{ name, bid, key string }
type testDirEntries []testDirEntry

var simpleDirTests = []struct {
	entries testDirEntries
	blobHex string
	key     string
	bid     string
}{
	{ // Empty Directory
		testDirEntries{},
		"01c659",
		"0129d7159641f64847d66fc4091d1320ff201147e2ca7e221080ce08933f1e1fd3",
		"cc347605074b230f9ca42f53c0f16475e3560df75c9378c0e9f7608781a6a04127f178bd428a10c1442b608e239148283a9e52f3bf0efdf514dfd7e1f9326372"},
}

func TestSimpleDirs(t *testing.T) {
	for _, test := range simpleDirTests {
		key := strings.Replace(test.key, " ", "", -1)
		bid := strings.Replace(test.bid, " ", "", -1)
		blob, _ := hex.DecodeString(strings.Replace(test.blobHex, " ", "", -1))

		m := NewMemoryBlobStorage()
		dw := DirBlobWriter{Storage: m}
		for _, entry := range test.entries {
			dw.AddEntry(entry.name, entry.bid, entry.key)
		}

		rbid, rkey, err := dw.Finalize()
		if err != nil {
			t.Error(err)
		}

		if rbid != bid {
			t.Errorf("Invalid blob id generated, got: %v..., expected: %v...", rbid[:16], bid[:16])
		}

		if rkey != key {
			t.Errorf("Invalid key generated, got: %v..., expected: %v...", rkey[:16], key[:16])
		}

		reader, err := m.NewBlobReader(rbid)
		if err != nil {
			t.Errorf("Couldn't open the blob with id: %v... for reading: %v", rbid[:16], err)
		} else {

			readBytes, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Errorf("Couldn't read the blob with id: %v..., error: %v", rbid[:16], err)
			} else if !bytes.Equal(readBytes, blob) {

				readBytesHex, blobHex := "", ""

				if len(readBytes) > 10 {
					readBytesHex = hex.EncodeToString(readBytes[:10]) + "..."
				} else {
					readBytesHex = hex.EncodeToString(readBytes)
				}

				if len(blob) > 10 {
					blobHex = hex.EncodeToString(blob[:10]) + "..."
				} else {
					blobHex = hex.EncodeToString(blob)
				}

				t.Errorf("The blob with id: %v... has invalid content, got: %v, expected %v",
					rbid[:16],
					readBytesHex,
					blobHex)
			}
		}
	}
}
