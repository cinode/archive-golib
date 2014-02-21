// Copyright 2013-2014 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"io"
	"unicode/utf8"
)

var (
	ErrDeserializeIntegerOverflow     = errors.New("Couldn not deserialize integer value - overflow of a 64-bit value detected")
	ErrDeserializeIntegerPaddingZeros = errors.New("Couldn not deserialize integer value - padding zeroes detected")
	ErrDeserializeStringToLarge       = errors.New("Could not deserialize string value due to invalid length")
	ErrDeserializeStringNotUTF8       = errors.New("Could not deserialize string value - not a UTF-8 sequence")
)

func SerializeInt(v uint64, w io.Writer) error {

	buff := make([]byte, 0, 10)

	for {
		b := byte(v & 0x7F)
		v = v >> 7
		if v != 0 {
			b |= 0x80
		}

		buff = append(buff, b)

		if v == 0 {
			break
		}
	}

	_, err := w.Write(buff)
	return err
}

func SerializeBuffer(data []byte, w io.Writer) error {

	if err := SerializeInt(uint64(len(data)), w); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func SerializeString(s string, w io.Writer) error {
	return SerializeBuffer([]byte(s), w)
}

func DeserializeInt(r io.Reader) (v uint64, err error) {

	v, s := 0, uint(0)
	buff := []byte{0}

	for ; ; s += 7 {

		// Early overflow detection
		if s >= 64 {
			return 0, ErrDeserializeIntegerOverflow
		}

		// Get next byte
		if _, err = r.Read(buff); err != nil {
			return 0, err
		}

		// Fill in the data in returned value
		v |= uint64(buff[0]&0x7F) << s

		// Overflow will cut some bits we won't be able to restore
		if (v >> s) != uint64(buff[0]&0x7F) {
			return 0, ErrDeserializeIntegerOverflow
		}

		// Highest bit in a byte means we shall continue
		if (buff[0] & 0x80) == 0 {
			break
		}
	}

	// No padding zeros allowed
	if ((buff[0] & 0x7F) == 0) && (s > 0) {
		return 0, ErrDeserializeIntegerPaddingZeros
	}

	return
}

func DeserializeBuffer(r io.Reader, maxLength uint64) (data []byte, err error) {
	length, err := DeserializeInt(r)
	if err != nil {
		return
	}

	if (length < 0) || (length > maxLength) {
		return nil, ErrDeserializeStringToLarge
	}

	buffer := make([]byte, length)
	if _, err = io.ReadFull(r, buffer); err != nil {
		return
	}

	return buffer, nil
}

func DeserializeString(r io.Reader, maxLength uint64) (s string, err error) {

	buffer, err := DeserializeBuffer(r, maxLength)
	if err != nil {
		return
	}

	if !utf8.Valid(buffer) {
		return "", ErrDeserializeStringNotUTF8
	}

	return string(buffer), nil
}
