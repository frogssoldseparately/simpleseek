package sreader

import (
	"archive/zip"
	"encoding/binary"
)

type Int interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type Float interface {
	float32
}

type Number interface {
	Int | Float
}

type SimpleReader struct {
	buf        *[]byte
	offset     *uint32
	endianness *binary.ByteOrder
}

type SimpleZipReader struct {
	name        string
	archive     *zip.ReadCloser
	entries     []*zip.File
	byName      map[string]*zip.File
	byExtension map[string][]*zip.File
}
