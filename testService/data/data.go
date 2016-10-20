// Package data provide test data
package data

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"
)

const (
	byteL = 1 << iota
	int16L
	int32L
	int64L
)

// Reply is an immutable bytes for simple test
// value:
//   99(Uint32) + 99(Uint32) + 99(Uint16) + 99(Uint8) + 99(int32) + success!(bytes)
var Reply []byte

// init initial Reply bytes
func init() {
	t := time.Date(2016, time.October, 1, 12, 0, 0, 0, time.Local)
	replyLen := int32L + int64L + byteL + int32L + int32L + int16L + byteL + int32L + 8*byteL
	Reply = make([]byte, replyLen)
	length := 0

	binary.BigEndian.PutUint32(Reply[length:], uint32(replyLen))
	length += int32L
	binary.BigEndian.PutUint64(Reply[length:], math.Float64bits(float64(t.UnixNano())))
	length += int64L
	Reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint32(Reply[length:], uint32(99))
	length += int32L
	binary.BigEndian.PutUint32(Reply[length:], uint32(99))
	length += int32L
	binary.BigEndian.PutUint16(Reply[length:], uint16(99))
	length += int16L
	Reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint32(Reply[length:], uint32(99))
	length += int32L

	copy(Reply[length:], "success!")
}

// RandomUserID return a random user id message.
// value:
//   randomID(uint32)
func RandomUserID() (result []byte) {
	t := time.Now()
	replyLen := int32L + int64L + byteL + int32L
	result = make([]byte, replyLen)
	length := 0

	binary.BigEndian.PutUint32(result[length:], uint32(replyLen))
	length += int32L
	binary.BigEndian.PutUint64(result[length:], math.Float64bits(float64(t.UnixNano())))
	length += int64L
	result[length] = byte(212)
	length += byteL
	binary.BigEndian.PutUint32(result[length:], uint32(rand.Int31()))

	return
}
