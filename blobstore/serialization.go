package blobstore

import (
	"bytes"
)

func serializeInt(v int64, buff *bytes.Buffer) {
	for {
		b := byte(v & 0x7F)
		v = v >> 7
		if v != 0 {
			b |= 0x80
		}

		buff.WriteByte(b)

		if v == 0 {
			break
		}
	}
}

func serializeString(s string, buff *bytes.Buffer) {
	data := []byte(s)
	serializeInt(int64(len(data)),buff)
	buff.Write(data)
}