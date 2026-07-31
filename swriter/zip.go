package swriter

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"hash/crc32"
	"time"
)

var zipTime uint16
var zipDate uint16

func TakeTimestamp() {
	now := time.Now()
	year := uint16(now.Year() - 1980)
	month := uint16(now.Month())
	day := uint16(now.Day())
	zipDate = year<<9 | month<<5 | day
	hour := uint16(now.Hour())
	minute := uint16(now.Minute())
	zipTime = hour<<11 | minute<<5
}

func WriteZipEntry(zippable Zippable, modLocalWriter *SimpleWriter, modCentralWriter *SimpleWriter, offsetAdjustment uint32) error {
	offset := modLocalWriter.GetLength() + offsetAdjustment
	uncompressed, compressed, err := toWritable(zippable)
	if err != nil {
		return err
	}
	lw := NewEmptySimpleWriter(binary.LittleEndian)
	cw := NewEmptySimpleWriter(binary.LittleEndian)

	writeLocalZip(zippable, lw, uncompressed, compressed)
	writeCentralZipHeader(zippable, cw, uncompressed, compressed, offset)
	modLocalWriter.CopyFrom(lw)
	modCentralWriter.CopyFrom(cw)
	return nil
}

func WriteCentralDirectoryEndRecord(w *SimpleWriter, fn uint16, cDirLen uint32, cDirStart uint32) {
	Write[uint32](w, 0x06054B50)
	Write[uint16](w, 0x0)
	Write[uint16](w, 0x0)
	Write(w, fn)
	Write(w, fn)
	Write(w, cDirLen)
	Write(w, cDirStart)
	Write[uint16](w, 0x0)
}

func writeLocalZip(zippable Zippable, lw *SimpleWriter, uncompressed *[]byte, compressed *[]byte) {
	compLen := uint32(len(*compressed))
	writeLocalZipHeader(zippable, lw, uncompressed, compLen)
	WriteRaw(lw, compressed)
}

func writeLocalZipHeader(zippable Zippable, lw *SimpleWriter, uncompressed *[]byte, compressedLength uint32) {
	bitFlag := zippable.GetBitFlag()
	compressionMethod := zippable.GetCompression()
	uncompressedHash := crc32.ChecksumIEEE(*uncompressed)
	uncompressedLength := uint32(len(*uncompressed))
	filename := zippable.GetFilename()
	filenameLength := uint16(len(filename))

	Write[uint32](lw, 0x04034B50)
	Write[uint16](lw, 0x14)
	Write(lw, bitFlag)
	Write(lw, compressionMethod)
	Write(lw, zipTime)
	Write(lw, zipDate)
	Write(lw, uncompressedHash)
	Write(lw, compressedLength)
	Write(lw, uncompressedLength)
	Write(lw, filenameLength)
	Write[uint16](lw, 0x0)
	WriteString(lw, filename, false)
}

func writeCentralZipHeader(zippable Zippable, cw *SimpleWriter, uncompressed *[]byte, compressed *[]byte, localOffset uint32) {
	bitFlag := zippable.GetBitFlag()
	compressionMethod := zippable.GetCompression()
	uncompressedHash := crc32.ChecksumIEEE(*uncompressed)
	compressedLength := uint32(len(*compressed))
	uncompressedLength := uint32(len(*uncompressed))
	filename := zippable.GetFilename()
	filenameLength := uint16(len(filename))
	permissions := zippable.GetAttributes()
	Write[uint32](cw, 0x02014B50)
	Write[uint16](cw, 0x14)
	Write[uint16](cw, 0x14)
	Write(cw, bitFlag)
	Write(cw, compressionMethod)
	Write(cw, zipTime)
	Write(cw, zipDate)
	Write(cw, uncompressedHash)
	Write(cw, compressedLength)
	Write(cw, uncompressedLength)
	Write(cw, filenameLength)
	Write[uint64](cw, 0x0)

	Write(cw, permissions)
	Write(cw, localOffset)
	WriteString(cw, filename, false)
}

func toWritable(zippable Zippable) (*[]byte, *[]byte, error) {
	var uncompressed []byte
	if err := zippable.WriteLocalBody(NewSimpleWriter(&uncompressed, binary.LittleEndian)); err != nil {
		return nil, nil, err
	}
	var compressed []byte
	if zippable.GetCompression() != 0x8 {
		w := NewSimpleWriter(&compressed, binary.LittleEndian)
		WriteRaw(w, &uncompressed)
	} else {
		newBuf, err := toCompressed(&uncompressed)
		if err != nil {
			return nil, nil, err
		}
		compressed = *newBuf
	}
	return &uncompressed, &compressed, nil
}

func toCompressed(src *[]byte) (*[]byte, error) {
	var compressedEntry bytes.Buffer
	deflater, err := flate.NewWriter(&compressedEntry, 7)
	if err != nil {
		return nil, err
	}
	deflater.Write(*src)
	deflater.Close()
	buf := compressedEntry.Bytes()
	return &buf, nil
}
