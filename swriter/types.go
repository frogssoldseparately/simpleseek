package swriter

import "encoding/binary"

type Int interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type Float interface {
	float32
}

type Number interface {
	Int | Float
}

type SimpleWriter struct {
	buf        *[]byte
	offset     *uint32
	endianness *binary.ByteOrder
}

type Zippable interface {
	GetCompression() uint16
	GetFilename() string
	WriteLocalBody(*SimpleWriter) error
	GetBitFlag() uint16
	GetAttributes() uint32
}
