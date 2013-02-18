// Copyright 2013 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blobstore

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

var (
	ErrDeserializeStringToLarge = errors.New("Could not deserialize string value due to invalid length")
	ErrDeserializeStringNotUTF8 = errors.New("Could not deserialize string value - not a UTF-8 sequence")
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
	serializeInt(int64(len(data)), buff)
	buff.Write(data)
}

func deserializeInt(r io.Reader) (v int64, err error) {
	// TODO: Overflows
	v, s := 0, uint(0)
	buff := []byte{0}
	for ; ; s += 7 {

		// Get next byte
		if _, err = r.Read(buff); err != nil {
			return
		}

		// Fill in the data in returned value
		v |= int64(buff[0]&0x7F) << s

		if (buff[0] & 0x80) == 0 {
			break
		}
	}
	return
}

func deserializeString(r io.Reader, maxLength int64) (v string, err error) {
	length, err := deserializeInt(r)
	if err != nil {
		return
	}

	if (length < 0) || (length > maxLength) {
		return ErrDeserializeStringToLarge
	}

	buffer := make([]byte, length)
	if _, err = io.ReadFull(r, buffer); err != nil {
		return
	}

	if !utf8.Valid(buffer) {
		return ErrDeserializeStringNotUTF8
	}

	return string(buffer)
}
