// Package bufbo provide the interface and the default implement about bufferByte reading and writing,
// just like encoding/binary.ByteOrder, but it remembers the reading or writing state.
package bufbo

import (
	"bytes"
	b "encoding/binary"
	"io"
	"math"
)

const (
	// NotEnoughError causes while buffer or bytes space is not enough.
	NotEnoughError = "runtime error: index out of range or some error throw out."
)

// OrderReader specifies how to read byte sequences into 16-, 32-, or 64-bit unsigned integers from
// buffer or other type which could remember reading state.
type OrderReader interface {
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Float32() float32
	Float64() float64
}

// OrderWriter specifies how to write 16-, 32-, or 64-bit unsigned integers as byte sequences from
// buffer or other type which could remember writing state.
type OrderWriter interface {
	PutUint8(uint8)
	PutUint16(uint16)
	PutUint32(uint32)
	PutUint64(uint64)
	PutFloat32(float32)
	PutFloat64(float64)
}

// BytesWriter implement OrderWriter using []byte and binary.ByteOrder.
// it is not concurrently safe.
type BytesWriter struct {
	w      []byte
	length int
	endian b.ByteOrder
}

// NewBEBytesWriter creates BytesWriter by binary.BigEndian.
func NewBEBytesWriter(writer []byte) *BytesWriter {
	return &BytesWriter{
		w:      writer,
		endian: b.BigEndian,
	}
}

// NewLEBytesWriter creates BytesWriter by binary.LittleEndian.
func NewLEBytesWriter(writer []byte) *BytesWriter {
	return &BytesWriter{
		w:      writer,
		endian: b.LittleEndian,
	}
}

// PutUint8 writes one byte into bufWriter.
func (bsw *BytesWriter) PutUint8(v uint8) {
	bsw.w[bsw.length] = v
	bsw.length++
}

// PutUint16 writes two bytes into BytesWriter.
func (bsw *BytesWriter) PutUint16(v uint16) {
	bsw.endian.PutUint16(bsw.w[bsw.length:], v)
	bsw.length += 2
}

// PutUint32 writes 4 bytes into BytesWriter.
func (bsw *BytesWriter) PutUint32(v uint32) {
	bsw.endian.PutUint32(bsw.w[bsw.length:], v)
	bsw.length += 4
}

// PutUint64 writes 8 bytes into BytesWriter.
func (bsw *BytesWriter) PutUint64(v uint64) {
	bsw.endian.PutUint64(bsw.w[bsw.length:], v)
	bsw.length += 8
}

// PutFloat32 writes 4 bytes into BytesWriter.
func (bsw *BytesWriter) PutFloat32(v float32) {
	bsw.endian.PutUint32(bsw.w[bsw.length:], math.Float32bits(v))
	bsw.length += 4
}

// PutFloat64 writes 8 bytes into BytesWriter.
func (bsw *BytesWriter) PutFloat64(v float64) {
	bsw.endian.PutUint64(bsw.w[bsw.length:], math.Float64bits(v))
	bsw.length += 8
}

// bufWriter implement OrderWriter using io.Buffer and encoding.binary.ByteOrder.
//
// This struct provide high level api for writer bytes.
type bufWriter struct {
	writer io.Writer
	endian b.ByteOrder
}

// NewBEBufWriter creates bufWriter using bytes.Buffer and binary.BigEndian or binary.BigEndian.
func NewBEBufWriter(buffer *bytes.Buffer) OrderWriter {
	return &bufWriter{
		writer: buffer,
		endian: b.BigEndian,
	}
}

// NewLEBufWriter creates bufWriter using bytes.Buffer and binary.BigEndian or binary.LittleEndian.
func NewLEBufWriter(buffer *bytes.Buffer) OrderWriter {
	return &bufWriter{
		writer: buffer,
		endian: b.LittleEndian,
	}
}

// PutUint8 writes one byte into bufWriter.
func (bfw *bufWriter) PutUint8(v uint8) {
	n, _ := bfw.writer.Write([]byte{v})
	if n < 1 {
		panic(NotEnoughError)
	}
}

// PutUint16 writes two bytes into bufWriter.
func (bfw *bufWriter) PutUint16(v uint16) {
	byte2 := make([]byte, 2)
	bfw.endian.PutUint16(byte2, v)

	n, _ := bfw.writer.Write(byte2)
	if n < 2 {
		panic(NotEnoughError)
	}
}

// PutUint32 writes 4 bytes into bufWriter.
func (bfw *bufWriter) PutUint32(v uint32) {
	byte4 := make([]byte, 4)
	bfw.endian.PutUint32(byte4, v)

	n, _ := bfw.writer.Write(byte4)
	if n < 4 {
		panic(NotEnoughError)
	}
}

// PutUint64 writes 8 bytes into bufWriter.
func (bfw *bufWriter) PutUint64(v uint64) {
	byte8 := make([]byte, 8)
	bfw.endian.PutUint64(byte8, v)

	n, _ := bfw.writer.Write(byte8)
	if n < 8 {
		panic(NotEnoughError)
	}
}

// PutFloat32 writes 4 bytes into bufWriter.
func (bfw *bufWriter) PutFloat32(v float32) {
	byte4 := make([]byte, 4)
	bfw.endian.PutUint32(byte4, math.Float32bits(v))

	n, _ := bfw.writer.Write(byte4)
	if n < 4 {
		panic(NotEnoughError)
	}
}

// PutFloat64 writes 8 bytes into bufWriter.
func (bfw *bufWriter) PutFloat64(v float64) {
	byte8 := make([]byte, 8)
	bfw.endian.PutUint64(byte8, math.Float64bits(v))

	n, _ := bfw.writer.Write(byte8)
	if n < 8 {
		panic(NotEnoughError)
	}
}

// BytesReader implement OrderReader using io.Buffer and encoding.binary.BigEndian.
type BytesReader struct {
	r      []byte
	length int
	endian b.ByteOrder
}

// NewBEBytesReader creates BytesReader by binary.BigEndian.
func NewBEBytesReader(r []byte) *BytesReader {
	return &BytesReader{
		r:      r,
		endian: b.BigEndian,
	}
}

// NewLEBytesReader creates BytesReader by binary.LittleEndian.
func NewLEBytesReader(r []byte) *BytesReader {
	return &BytesReader{
		r:      r,
		endian: b.LittleEndian,
	}
}

// Uint8 read one byte from BytesReader.
func (bsr *BytesReader) Uint8() (result uint8) {
	result = bsr.r[bsr.length]
	bsr.length++
	return
}

// Uint16 read two bytes then convert them to uint16.
func (bsr *BytesReader) Uint16() (result uint16) {
	result = bsr.endian.Uint16(bsr.r[bsr.length:])
	bsr.length += 2
	return
}

// Uint32 read four bytes then convert them to uint32
func (bsr *BytesReader) Uint32() (result uint32) {
	result = bsr.endian.Uint32(bsr.r[bsr.length:])
	bsr.length += 4
	return
}

// Uint64 read eight bytes then convert them to uint64
func (bsr *BytesReader) Uint64() (result uint64) {
	result = bsr.endian.Uint64(bsr.r[bsr.length:])
	bsr.length += 8
	return
}

// Float32 read eight bytes then convert them to float32
func (bsr *BytesReader) Float32() (result float32) {
	result = math.Float32frombits(bsr.endian.Uint32(bsr.r[bsr.length:]))
	bsr.length += 4
	return
}

// Float64 read eight bytes then convert them to float64
func (bsr *BytesReader) Float64() (result float64) {
	result = math.Float64frombits(bsr.endian.Uint64(bsr.r[bsr.length:]))
	bsr.length += 8
	return
}
