package sreader

import (
	"encoding/binary"
	"io"
	"log"
)

func NewSimpleReader(f io.Reader, endianness binary.ByteOrder) *SimpleReader {
	var offset uint32
	buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return &SimpleReader{&buf, &offset, &endianness}
}

// Doubles as r.GetOffset() if you use r.Seek(0, 0)
func (r *SimpleReader) Seek(offset uint32, whence int) uint32 {
	switch whence {
	case 0:
		*r.offset = offset
	case 1:
		*r.offset += offset
	case 2:
		*r.offset = uint32(len(*r.buf)) + offset
	}
	return *r.offset
}

func (r *SimpleReader) GetBuffer() *[]byte {
	return r.buf
}

func (r *SimpleReader) GetLength() uint32 {
	return uint32(len(*r.buf))
}

func Read[T Number](r *SimpleReader) T {
	var out T
	n, err := binary.Decode((*r.buf)[*r.offset:], *r.endianness, &out)
	if err != nil {
		log.Fatal(err)
	}
	*r.offset += uint32(n)
	return out
}
