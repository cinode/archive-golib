package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var intSerializationData = []struct {
	value      uint64
	serialized []byte
}{
	{0x00, []byte{0x00}},
	{0x01, []byte{0x01}},
	{0x02, []byte{0x02}},
	{0x7F, []byte{0x7F}},
	{0x80, []byte{0x80, 0x01}},
	{0x81, []byte{0x81, 0x01}},
	{0x100, []byte{0x80, 0x02}},
	{0xFFFFFFFFFFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}},
}

func TestIntSerialization(t *testing.T) {

	for _, tst := range intSerializationData {

		var b bytes.Buffer

		err := serializeInt(tst.value, &b)

		if err != nil {
			t.Fatalf("Error while serializing int: %v", err)
			continue
		}

		s := b.Bytes()
		if !bytes.Equal(s, tst.serialized) {
			t.Fatalf(
				"Serialization of value %v failed, data not equal (%v vs %v)",
				tst.value,
				hex.EncodeToString(s),
				hex.EncodeToString(tst.serialized))
			continue
		}

		// Try to serialize back
		d, err := deserializeInt(&b)
		if err != nil {
			t.Fatalf("Error while deserializing int back: %v", err)
			continue
		}

		if d != tst.value {
			t.Fatalf("Incorrectly deserialized the value back, was: %v, got %v", tst.value, d)
			continue
		}

	}

}
