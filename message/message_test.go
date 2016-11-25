package message

import (
	"math"
	"testing"
	"time"
)

var msgTime time.Time

const testMessageLength = msgHeadSize + 8

// generate ...
func generateTestMessage() Message {
	ci := &ConnectInfo{
		UID: 2333,
		RID: 1,
	}

	bs, _ := ci.MarshalBinary()
	msgTime = time.Now()
	return &msg{
		t:         MsgConnect,
		timestamp: msgTime,
		body:      bs,
	}
}

// TestMessageBase ...
func TestMessageBase(t *testing.T) {
	m := generateTestMessage()

	if mLength := m.Size(); mLength != testMessageLength {
		t.Errorf("Length of Message is wrong, hope %d, get %d.", testMessageLength, mLength)
	}
	if mType := m.Type(); mType != MsgConnect {
		t.Errorf("Type of Message is wrong, hope %d, get %d.", MsgConnect, mType)
	}
	if mTime := m.Timestamp(); mTime != msgTime {
		t.Errorf("Time of Message is wrong, hope %v, get %v", msgTime, mTime)
	}
}

// TestMessageMarshalAndUnmarshal ...
func TestMessageMarshalAndUnmarshal(t *testing.T) {
	m := generateTestMessage()

	// MarshalBinary
	bs, err := m.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	bSize, mSize := int(bs[3]), m.Size()
	if bSize != mSize {
		t.Errorf("Length value in bytes is wrong, hope %d, get %d.", mSize, bSize)
	}

	// UnmarshalBinary
	m = new(msg)
	err = m.UnmarshalBinary(bs)
	if err != nil {
		t.Error(err)
	}
	if newMSize := m.Size(); newMSize != mSize {
		t.Errorf("Size of unmarshaled Message is wrong, hope %d, get %d.", mSize, newMSize)
	}
	mTimeUnix, msgTimeUnix := m.Timestamp().UnixNano(), msgTime.UnixNano()
	if duration := math.Abs(float64(msgTimeUnix - mTimeUnix)); duration > 1000 {
		t.Errorf("Duration between new message and old message is wrong, duration: %v.", duration)
	}
}
