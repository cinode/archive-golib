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
	ErrIntegerOverflow     = errors.New("integer overflow")
	ErrIntegerPaddingZeros = errors.New("padding zeroes are not allowed")
	ErrBufferToLarge       = errors.New("buffer size is to large")
	ErrStringToLarge       = errors.New("string length is to large")
	ErrStringNotUTF8       = errors.New("buffer is not a valid utf-8 string")
)

const (
	// 64-bit integer can be represented in 10 bytes, each representing 7 bits of the number
	maxNumberBytes = 10
)

// SerializeInt serializes integer value into writer
func SerializeInt(v uint64, w io.Writer) error {

	buff := make([]byte, 0, maxNumberBytes)

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

// SerializeBuffer serializes array of bytes into writer
func SerializeBuffer(data []byte, w io.Writer, maxLength uint64) error {

	l := uint64(len(data))
	if l > maxLength {
		return ErrBufferToLarge
	}

	if err := SerializeInt(l, w); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

// SerializeString serializes string into writer
func SerializeString(s string, w io.Writer, maxLength uint64) error {
	return SerializeBuffer([]byte(s), w, maxLength)
}

// DeserializeInt tries to deserialize integer value from reader
func DeserializeInt(r io.Reader) (uint64, error) {

	v, s, buff := uint64(0), uint(0), []byte{0}

	for ; ; s += 7 {

		// Early overflow detection
		if s >= 64 {
			return 0, ErrIntegerOverflow
		}

		// Get next byte
		if _, err := r.Read(buff); err != nil {
			return 0, err
		}

		// Fill in the data in returned value
		v |= uint64(buff[0]&0x7F) << s

		// Overflow will cut some bits we won't be able to restore
		if (v >> s) != uint64(buff[0]&0x7F) {
			return 0, ErrIntegerOverflow
		}

		// Highest bit in a byte means we shall continue
		if (buff[0] & 0x80) == 0 {
			break
		}
	}

	// No padding zeros allowed
	if ((buff[0] & 0x7F) == 0) && (s > 0) {
		return 0, ErrIntegerPaddingZeros
	}

	return v, nil
}

// DeserializeBuffer tries to deserialize buffer form reader
func DeserializeBuffer(r io.Reader, maxLength uint64) ([]byte, error) {
	length, err := DeserializeInt(r)
	if err != nil {
		return nil, err
	}

	if (length < 0) || (length > maxLength) {
		return nil, ErrBufferToLarge
	}

	buffer := make([]byte, length)
	if _, err = io.ReadFull(r, buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

// DeserializeString tries to deserialize string from reader
func DeserializeString(r io.Reader, maxLength uint64) (string, error) {

	buffer, err := DeserializeBuffer(r, maxLength)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(buffer) {
		return "", ErrStringNotUTF8
	}

	return string(buffer), nil
}
