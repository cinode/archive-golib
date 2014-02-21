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

		err := SerializeInt(tst.value, &b)

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
		d, err := DeserializeInt(&b)
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

var intDeserializationData = []struct {
	serialized []byte
	value      uint64
	bogus      bool
}{
	{[]byte{0x00}, 0x00, false},
	{[]byte{0x01}, 0x01, false},
	{[]byte{0x02}, 0x02, false},
	{[]byte{0x7F}, 0x7F, false},
	{[]byte{0x80, 0x01}, 0x80, false},
	{[]byte{0x81, 0x01}, 0x81, false},
	{[]byte{0x80, 0x02}, 0x100, false},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}, 0xFFFFFFFFFFFFFFFF, false},

	{[]byte{0x80, 0x00}, 0x00, true},
	{[]byte{0x81, 0x00}, 0x01, true},
	{[]byte{0x80, 0x82, 0x00}, 0x100, true},
	{[]byte{0x80}, 0x00, true},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}, 0, true},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x02}, 0, true},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}, 0, true},
}

func TestIntDesrialization(t *testing.T) {

	for _, tst := range intDeserializationData {

		v, err := DeserializeInt(bytes.NewBuffer(tst.serialized))

		if tst.bogus {
			if err == nil {
				t.Errorf("Dserialization of bogus data did not produce an error")
			}
			continue
		}

		if err != nil {
			t.Errorf("Couldn not deserialize data, expected value: %v, got error: %v", tst.value, err)
			continue
		}

		if v != tst.value {
			t.Errorf("Mismatch of deserialized value: %v, should be: %v", v, tst.value)
		}
	}
}
