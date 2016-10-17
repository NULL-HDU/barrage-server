package bufbo

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
	"testing"
)

const (
	byteL = 1 << iota
	int16L
	int32L
	int64L
)

var replyLen = byteL + int64L + int32L + int16L + byteL + int64L + int64L + int32L + byteL*9

//   99(Uint8) + 99(Uint64) + 99(Uint32) + 99(Uint16) + 99(Uint8) + 99(int64) + "9999 9999"
func generateBETestBytes() (reply []byte) {
	reply = make([]byte, replyLen)
	length := 0

	reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint64(reply[length:], uint64(99))
	length += int64L
	binary.BigEndian.PutUint32(reply[length:], uint32(99))
	length += int32L
	binary.BigEndian.PutUint16(reply[length:], uint16(99))
	length += int16L
	reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint64(reply[length:], uint64(99))
	length += int64L
	binary.BigEndian.PutUint64(reply[length:], math.Float64bits(float64(99)))
	length += int64L
	binary.BigEndian.PutUint32(reply[length:], math.Float32bits(float32(99)))
	length += int32L
	copy(reply[length:], []byte("9999 9999"))
	length += byteL * 9

	return
}

func TestBytesWriter(t *testing.T) {
	w := make([]byte, replyLen)
	bw := NewBEBytesWriter(w)

	bw.PutUint8(99)
	bw.PutUint64(99)
	bw.PutUint32(99)
	bw.PutUint16(99)
	bw.PutUint8(99)
	bw.PutUint64(99)
	bw.PutFloat64(99)
	bw.PutFloat32(99)
	bw.PutStr("9999 9999")

	if r := bytes.Compare(w, generateBETestBytes()); r != 0 {
		t.Errorf("w should be the same as the result of generateBETestBytes,"+
			"but the result of bytes.Compare is %v", r)
	}

	defer func() {
		if err := recover(); !strings.Contains(err.(error).Error(), "index out of range") {
			t.Error("w should panic and throw NotEnoughError!")
		}
	}()
	bw.PutUint8(99)
}

func TestBytesReader(t *testing.T) {
	r := generateBETestBytes()
	br := NewBEBytesReader(r)

	if result := br.Uint8(); result != 99 {
		t.Errorf("1st number should be 99(uint8), but get %v.", result)
	}
	if result := br.Uint64(); result != 99 {
		t.Errorf("2nd number should be 99(uint64), but get %v.", result)
	}
	if result := br.Uint32(); result != 99 {
		t.Errorf("3rd number should be 99(uint32), but get %v.", result)
	}
	if result := br.Uint16(); result != 99 {
		t.Errorf("4th number should be 99(uint16), but get %v.", result)
	}
	if result := br.Uint8(); result != 99 {
		t.Errorf("5th number should be 99(uint8), but get %v.", result)
	}
	if result := br.Uint64(); result != 99 {
		t.Errorf("6th number should be 99(uint64), but get %v.", result)
	}
	if result := br.Float64(); result != float64(99) {
		t.Errorf("7th number should be 99(float64), but get %v.", result)
	}

	if result := br.Float32(); result != float32(99) {
		t.Errorf("8th number should be 99(float32), but get %v.", result)
	}
	if result := br.Str(9 * byteL); result != "9999 9999" {
		t.Errorf("9th string should be '9999 9999', but get %v.", result)
	}

	defer func() {
		if err := recover(); !strings.Contains(err.(error).Error(), "index out of range") {
			t.Error("w should panic and throw NotEnoughError!")
		}
	}()
	br.Uint8()
}

func TestBufWriter(t *testing.T) {
	var w bytes.Buffer
	bfw := NewBEBufWriter(&w)

	bfw.PutUint8(99)
	bfw.PutUint64(99)
	bfw.PutUint32(99)
	bfw.PutUint16(99)
	bfw.PutUint8(99)
	bfw.PutUint64(99)
	bfw.PutFloat64(99)
	bfw.PutFloat32(99)
	bfw.PutStr("9999 9999")

	if r := bytes.Compare(w.Bytes(), generateBETestBytes()); r != 0 {
		t.Errorf("w should be the same as the result of generateBETestBytes,"+
			"but the result of bytes.Compare is %v", r)
	}
}
