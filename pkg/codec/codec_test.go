package codec_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/capucinoxx/forlorn/pkg/codec"
)

func floatEqual[T float32 | float64](a, b T) bool {
	const epsilon = 1e-6
	return math.Abs(float64(a)-float64(b)) < epsilon
}

type FakeStruct struct {
	name   string
	health float32
	points uint64
	posX   float64
	posY   float64

	tiles []struct {
		posX float64
		posY float64
	}
}

func (f *FakeStruct) Encode(w codec.Writer) (err error) {
	if err = w.WriteString(f.name); err != nil {
		return
	}

	if err = w.WriteFloat32(f.health); err != nil {
		return
	}

	if err = w.WriteUint64(f.points); err != nil {
		return
	}

	if err = w.WriteFloat64(f.posX); err != nil {
		return
	}

	if err = w.WriteFloat64(f.posY); err != nil {
		return
	}

	if err = w.WriteUint32(uint32(len(f.tiles))); err != nil {
		return
	}

	for _, tile := range f.tiles {
		if err = w.WriteFloat64(tile.posX); err != nil {
			return
		}

		if err = w.WriteFloat64(tile.posY); err != nil {
			return
		}
	}

	return
}

func (f *FakeStruct) Decode(r codec.Reader) (err error) {
	if f.name, err = r.ReadString(); err != nil {
		return
	}

	if f.health, err = r.ReadFloat32(); err != nil {
		return
	}

	if f.points, err = r.ReadUint64(); err != nil {
		return
	}

	if f.posX, err = r.ReadFloat64(); err != nil {
		return
	}

	if f.posY, err = r.ReadFloat64(); err != nil {
		return
	}

	var size uint32
	if size, err = r.ReadUint32(); err != nil {
		return
	}

	var px, py float64
	for i := uint32(0); i < size; i++ {
		if px, err = r.ReadFloat64(); err != nil {
			return
		}

		if py, err = r.ReadFloat64(); err != nil {
			return
		}

		f.tiles = append(f.tiles, struct {
			posX float64
			posY float64
		}{px, py})
	}

	return nil
}

func (f *FakeStruct) Equals(oth *FakeStruct) bool {
	ok := true
	ok = ok && (f.name == oth.name)
	ok = ok && floatEqual(f.health, oth.health)
	ok = ok && (f.points == oth.points)
	ok = ok && floatEqual(f.posX, oth.posX)
	ok = ok && floatEqual(f.posY, oth.posY)
	ok = ok && (len(f.tiles) == len(oth.tiles))

	if !ok {
		return false
	}

	for i := 0; i < len(f.tiles); i++ {
		ok = ok && floatEqual(f.tiles[i].posX, f.tiles[i].posX)
		ok = ok && floatEqual(f.tiles[i].posY, f.tiles[i].posY)
	}

	return ok
}

func TestByteReader(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	readerLittle := codec.NewByteReader(data, binary.LittleEndian)
	readerBig := codec.NewByteReader(data, binary.BigEndian)

	tests := []struct {
		name     string
		readFunc func() (interface{}, error)
		expected interface{}
	}{
		{"ReadUint8", func() (interface{}, error) { return readerLittle.ReadUint8() }, uint8(0x01)},
		{"ReadUint16LittleEndian", func() (interface{}, error) { return readerLittle.ReadUint16() }, uint16(0x0302)},
		{"ReadUint32LittleEndian", func() (interface{}, error) { return readerLittle.ReadUint32() }, uint32(0x07060504)},
		{"ReadUint64LittleEndian", func() (interface{}, error) { return readerLittle.ReadUint64() }, uint64(0x08)},
		{"ReadUint16BigEndian", func() (interface{}, error) { readerBig.ResetPos(); return readerBig.ReadUint16() }, uint16(0x0102)},
		{"ReadUint32BigEndian", func() (interface{}, error) { readerBig.ResetPos(); return readerBig.ReadUint32() }, uint32(0x01020304)},
		{"ReadUint64BigEndian", func() (interface{}, error) { readerBig.ResetPos(); return readerBig.ReadUint64() }, uint64(0x0102030405060708)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := tt.readFunc()
			if result != tt.expected {
				t.Errorf("Expected: %x, got: %x", tt.expected, result)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name       string
		data       FakeStruct
		shouldFail bool
	}{
		{"empty struct", FakeStruct{}, false},
		{"stuffed struct empty tiles", FakeStruct{name: "John Doe", health: 100.954, points: 111, posX: 123.765544, posY: 144.7665544}, false},
		{"stuffed with tiles", FakeStruct{name: "John Doe", health: 100.954, points: 111, posX: 123.765544, posY: 144.7665544, tiles: []struct {
			posX float64
			posY float64
		}{
			{123.665, 222.222222},
			{12333344.333, 22.0},
			{-1234.65, -33333.76},
		}}, false},
		{"corrupted message", FakeStruct{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := codec.NewByteWriter(binary.LittleEndian)

			if err := tt.data.Encode(writer); err != nil {
				t.Errorf("Unexpected error %v\n", err)
			}

			bytes := writer.Bytes()
			if tt.shouldFail {
				bytes = bytes[:len(bytes)-16]
			}

			reader := codec.NewByteReader(bytes, binary.LittleEndian)
			var res FakeStruct
			if err := res.Decode(reader); !tt.shouldFail && err != nil {
				t.Errorf("Unexpected error %v\n", err)
			} else if tt.shouldFail && err == nil {
				t.Errorf("Expected error")
			}

			if !tt.shouldFail && !tt.data.Equals(&res) {
				t.Errorf("Expect equals (%+v) != (%+v)", tt.data, res)
			}
		})
	}
}
