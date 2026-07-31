package swriter

import (
	"encoding/binary"
	"log"
)

func NewEmptySimpleWriter(endianness binary.ByteOrder) *SimpleWriter {
	var buf []byte
	return NewSimpleWriter(&buf, endianness)
}

func NewSimpleWriter(buf *[]byte, endianness binary.ByteOrder) *SimpleWriter {
	var offset uint32
	return &SimpleWriter{buf, &offset, &endianness}
}

func (w *SimpleWriter) GetBuffer() *[]byte {
	return w.buf
}

func (w *SimpleWriter) GetLength() uint32 {
	return *w.offset
}

func Write[T Number](w *SimpleWriter, data T) {
	newBuf, err := binary.Append(*w.buf, *w.endianness, data)
	if err != nil {
		log.Fatal(err)
	}
	size := len(newBuf) - len(*w.buf)
	*w.buf = newBuf
	*w.offset += uint32(size)
}

func WriteString(w *SimpleWriter, s string, writeLength bool) {
	size := uint32(len(s))
	if writeLength {
		Write(w, size)
	}
	buf := []byte(s[:])
	WriteRaw(w, &buf)
}

func WriteRaw(w *SimpleWriter, buf *[]byte) {
	size := uint32(len(*buf))
	newBuf, err := binary.Append(*w.buf, *w.endianness, *buf)
	if err != nil {
		log.Fatal(err)
	}
	*w.buf = newBuf
	*w.offset += size
}

func (dest *SimpleWriter) CopyFrom(src *SimpleWriter) {
	WriteRaw(dest, src.GetBuffer())
}
