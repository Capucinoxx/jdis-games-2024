package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type Reader interface {
	io.Reader
	io.ByteReader
	io.Seeker
	ReadUint8() (uint8, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint64() (uint64, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadBool() (bool, error)
	ReadBytes(n int) ([]byte, error)
	ReadString() (string, error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
	WriteByte(b byte) error
	WriteUint8(value uint8) error
	WriteUint16(value uint16) error
	WriteUint32(value uint32) error
	WriteUint64(value uint64) error
	WriteInt8(value int8) error
	WriteInt16(value int16) error
	WriteInt32(value int32) error
	WriteInt64(value int64) error
	WriteFloat32(value float32) error
	WriteFloat64(value float64) error
	WriteBool(value bool) error
	WriteBytes(data []byte) (int, error)
	WriteString(s string) error
}

type BinaryDecoder interface {
	Decode(r Reader) error
}

type BinaryEncoder interface {
	Encode(w io.Writer) error
}

type ByteReader struct {
	data  []byte
	pos   int
	order binary.ByteOrder
}

func NewByteReader(data []byte, order binary.ByteOrder) *ByteReader {
	return &ByteReader{data: data, order: order}
}

func (r *ByteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *ByteReader) ReadByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *ByteReader) ResetPos() {
	r.pos = 0
}

func (r *ByteReader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(r.pos) + offset
	case io.SeekEnd:
		abs = int64(len(r.data)) + offset
	default:
		return 0, errors.New("binstruct: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("binstruct: negative position")
	}
	r.pos = int(abs)
	return abs, nil
}

func (r *ByteReader) ReadUint8() (uint8, error) {
	b, err := r.ReadByte()
	return uint8(b), err
}

func (r *ByteReader) ReadUint16() (uint16, error) {
	b := make([]byte, 2)
	if _, err := r.Read(b); err != nil {
		return 0, err
	}
	return r.order.Uint16(b), nil
}

func (r *ByteReader) ReadUint32() (uint32, error) {
	b := make([]byte, 4)
	if _, err := r.Read(b); err != nil {
		return 0, err
	}
	return r.order.Uint32(b), nil
}

func (r *ByteReader) ReadUint64() (uint64, error) {
	b := make([]byte, 8)
	if _, err := r.Read(b); err != nil {
		return 0, err
	}
	return r.order.Uint64(b), nil
}

func (r *ByteReader) ReadInt8() (int8, error) {
	b, err := r.ReadByte()
	return int8(b), err
}

func (r *ByteReader) ReadInt16() (int16, error) {
	b, err := r.ReadUint16()
	return int16(b), err
}

func (r *ByteReader) ReadInt32() (int32, error) {
	b, err := r.ReadUint32()
	return int32(b), err
}

func (r *ByteReader) ReadInt64() (int64, error) {
	b, err := r.ReadUint64()
	return int64(b), err
}

func (r *ByteReader) ReadFloat32() (float32, error) {
	bits, err := r.ReadUint32()
	return math.Float32frombits(bits), err
}

func (r *ByteReader) ReadFloat64() (float64, error) {
	bits, err := r.ReadUint64()
	return math.Float64frombits(bits), err
}

func (r *ByteReader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

func (r *ByteReader) ReadBytes(n int) ([]byte, error) {
	if r.pos+n > len(r.data) {
		return nil, io.EOF
	}
	b := r.data[r.pos : r.pos+n]
	r.pos += n
	return b, nil
}

func (r *ByteReader) ReadString() (string, error) {
	start := r.pos
	for r.pos < len(r.data) {
		if r.data[r.pos] == 0 {
			str := string(r.data[start:r.pos])
			r.pos++
			return str, nil
		}
		r.pos++
	}
	return "", io.EOF
}

type ByteWriter struct {
	buffer bytes.Buffer
	order  binary.ByteOrder
}

func NewByteWriter(order binary.ByteOrder) *ByteWriter {
	return &ByteWriter{buffer: bytes.Buffer{}, order: order}
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *ByteWriter) WriteByte(b byte) error {
	return w.buffer.WriteByte(b)
}

func (w *ByteWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

func (w *ByteWriter) WriteUint8(value uint8) error {
	return w.WriteByte(byte(value))
}

func (w *ByteWriter) WriteUint16(value uint16) error {
	var b [2]byte
	w.order.PutUint16(b[:], value)
	_, err := w.Write(b[:])
	return err
}

func (w *ByteWriter) WriteUint32(value uint32) error {
	var b [4]byte
	w.order.PutUint32(b[:], value)
	_, err := w.Write(b[:])
	return err
}

func (w *ByteWriter) WriteUint64(value uint64) error {
	var b [8]byte
	w.order.PutUint64(b[:], value)
	_, err := w.Write(b[:])
	return err
}

func (w *ByteWriter) WriteInt8(value int8) error {
	return w.WriteByte(byte(value))
}

func (w *ByteWriter) WriteInt16(value int16) error {
	return w.WriteUint16(uint16(value))
}

func (w *ByteWriter) WriteInt32(value int32) error {
	return w.WriteUint32(uint32(value))
}

func (w *ByteWriter) WriteInt64(value int64) error {
	return w.WriteUint64(uint64(value))
}

func (w *ByteWriter) WriteFloat32(value float32) error {
	bits := math.Float32bits(value)
	return w.WriteUint32(bits)
}

func (w *ByteWriter) WriteFloat64(value float64) error {
	bits := math.Float64bits(value)
	return w.WriteUint64(bits)
}

func (w *ByteWriter) WriteBool(value bool) error {
	var b byte
	if value {
		b = 1
	} else {
		b = 0
	}
	return w.WriteByte(b)
}

func (w *ByteWriter) WriteBytes(data []byte) (int, error) {
	return w.Write(data)
}

func (w *ByteWriter) WriteString(s string) error {
	if _, err := w.WriteBytes([]byte(s)); err != nil {
		return err
	}
	return w.WriteByte(0)
}
