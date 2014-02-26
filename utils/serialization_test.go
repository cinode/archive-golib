package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func hexDump(data []byte) string {
	if data == nil {
		return "nil"
	}

	if len(data) > 8 {
		return hex.EncodeToString(data[:8]) + "..."
	}

	return hex.EncodeToString(data)
}

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
				hexDump(s),
				hexDump(tst.serialized))
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

func TestIntDeserialization(t *testing.T) {

	for _, tst := range intDeserializationData {

		v, err := DeserializeInt(bytes.NewBuffer(tst.serialized))

		if tst.bogus {
			if err == nil {
				t.Errorf("Dserialization of bogus data did not produce an error")
			}
			continue
		}

		if err != nil {
			t.Errorf("Could not deserialize data, expected value: %v, got error: %v", tst.value, err)
			continue
		}

		if v != tst.value {
			t.Errorf("Mismatch of deserialized value: %v, should be: %v", v, tst.value)
		}
	}
}

var bufferSerializationData = []struct {
	value      []byte
	serialized []byte
	err        error
	maxLength  uint64
}{
	{[]byte{}, []byte{0}, nil, 0},
	{[]byte{0}, []byte{1, 0}, nil, 1},
	{[]byte{0}, nil, ErrBufferToLarge, 0},
	{[]byte{7, 13}, []byte{2, 7, 13}, nil, 2},
	{[]byte{7, 13}, nil, ErrBufferToLarge, 1},
	{[]byte{7, 13}, nil, ErrBufferToLarge, 0},
}

func TestBufferSerialization(t *testing.T) {

	for _, tst := range bufferSerializationData {

		var b bytes.Buffer

		err := SerializeBuffer(tst.value, &b, tst.maxLength)

		if tst.err != nil {

			if err == nil {
				t.Errorf("Was able to serialize buffer but expected error")
				continue
			}

			continue
		}

		if err != nil {
			t.Errorf("Could not serialize buffer: %v", err)
			continue
		}

		if !bytes.Equal(b.Bytes(), tst.serialized) {
			t.Errorf("Serialized buffer value is incorrect, expected: %v, is: %v", hexDump(tst.serialized), hexDump(b.Bytes()))
			continue
		}

		v2, err := DeserializeBuffer(&b, tst.maxLength)
		if err != nil {
			t.Errorf("Unexpected deserialization error: %v", err)
			continue
		}

		if !bytes.Equal(v2, tst.value) {
			t.Errorf("Buffer value after deserialization is invalid, expected: %v, is: %v", hexDump(tst.value), v2)
			continue
		}
	}
}

var bufferDeserializationData = []struct {
	serialized []byte
	value      []byte
	err        error
	maxLength  uint64
}{
	{[]byte{0}, []byte{}, nil, 0},
	{[]byte{1, 0}, []byte{0}, nil, 1},
	{[]byte{1, 0}, nil, ErrBufferToLarge, 0},
	{[]byte{2, 7, 13}, []byte{7, 13}, nil, 2},
	{[]byte{2, 7, 13}, nil, ErrBufferToLarge, 1},
	{[]byte{2, 7, 13}, nil, ErrBufferToLarge, 0},
}

func TestBufferDeserialization(t *testing.T) {

	for _, tst := range bufferDeserializationData {

		v, err := DeserializeBuffer(bytes.NewBuffer(tst.serialized), tst.maxLength)

		if tst.err != nil {

			if err == nil {
				t.Errorf("Was able to serialize buffer but expected error")
				continue
			}

			continue
		}

		if err != nil {
			t.Errorf("Could not deserialize buffer: %v", err)
			continue
		}

		if !bytes.Equal(v, tst.value) {
			t.Errorf("Buffer value after deserialization is invalid, expected: %v, is: %v", hexDump(tst.value), hexDump(v))
			continue
		}

	}

}

var stringSerializationData = []struct {
	value      string
	serialized []byte
	err        error
	maxLength  uint64
}{
	{"", []byte{0}, nil, 0},
	{"a", []byte{1, 'a'}, nil, 1},
	{"a", nil, ErrBufferToLarge, 0},
	{"abc", []byte{3, 'a', 'b', 'c'}, nil, 3},
	{"abc", nil, ErrBufferToLarge, 2},
	{"abc", nil, ErrBufferToLarge, 1},
	{"abc", nil, ErrBufferToLarge, 0},
	{"\u2318", []byte{3, 0xe2, 0x8c, 0x98}, nil, 3},
}

func TestStringSerialization(t *testing.T) {

	for _, tst := range stringSerializationData {

		var b bytes.Buffer

		err := SerializeString(tst.value, &b, tst.maxLength)

		if tst.err != nil {

			if err == nil {
				t.Errorf("Was able to serialize string but expected error")
				continue
			}

			continue
		}

		if err != nil {
			t.Errorf("Could not serialize string: %v", err)
			continue
		}

		if !bytes.Equal(b.Bytes(), tst.serialized) {
			t.Errorf("Serialized string value is incorrect, expected: %v, is: %v", hexDump(tst.serialized), hexDump(b.Bytes()))
			continue
		}

		v2, err := DeserializeString(&b, tst.maxLength)
		if err != nil {
			t.Errorf("Unexpected deserialization error: %v", err)
			continue
		}

		if v2 != tst.value {
			t.Errorf("String value after deserialization is invalid, expected: %v, is: %v", tst.value, v2)
			continue
		}
	}
}

var stringDeserializationData = []struct {
	serialized []byte
	value      string
	err        error
	maxLength  uint64
}{
	{[]byte{0}, "", nil, 0},
	{[]byte{1, 'a'}, "a", nil, 1},
	{[]byte{1, 'a'}, "", ErrBufferToLarge, 0},
	{[]byte{2, 'a', 'b'}, "ab", nil, 2},
	{[]byte{2, 'a', 'b'}, "", ErrBufferToLarge, 1},
	{[]byte{2, 'a', 'b'}, "", ErrBufferToLarge, 0},
	{[]byte{2, 0x80, 'a'}, "", ErrStringNotUTF8, 2},
}

func TestStrigDeserialization(t *testing.T) {

	for _, tst := range stringDeserializationData {

		v, err := DeserializeString(bytes.NewBuffer(tst.serialized), tst.maxLength)

		if tst.err != nil {

			if err == nil {
				t.Errorf("Was able to serialize string but expected error")
				continue
			}

			continue
		}

		if err != nil {
			t.Errorf("Could not deserialize string: %v", err)
			continue
		}

		if v != tst.value {
			t.Errorf("String value after deserialization is invalid, expected: %v, is: %v", tst.value, v)
			continue
		}

	}

}
