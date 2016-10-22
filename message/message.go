package message

import (
	"barrage-server/base"
	"barrage-server/libs/bufbo"
	"errors"
	"time"
)

var logger = base.Log

type msgType uint8

const (
	// backend -> frontend

	// MsgGameOver is used when server will break off.
	MsgGameOver msgType = 0x0b
	// MsgSpecialMessage is used to send messages not related to game engine.
	MsgSpecialMessage msgType = 0x0a
	// MsgPlayground is used when backend send balls info to frontend.
	MsgPlayground msgType = 0x07
	// MsgAirplaneCreated is used to send airplane of userself to frontend while user connecting
	// into game.
	MsgAirplaneCreated msgType = 0x06

	// frontend -> backend

	// MsgUserSelf is used when frontend send balls info to backend.
	MsgUserSelf msgType = 0x0c
	// MsgConnect is used when user want to connect game.
	MsgConnect msgType = 0x09
	// MsgDisconnect is used when user want to leave game early.
	MsgDisconnect msgType = 0x08
)

const (
	// msgHeadSize is the size of message head
	// now, this includes length, timestamp, type.
	msgHeadSize = 13
)

var (
	// ErrInvalidMessage is the signature of invalid message.
	ErrInvalidMessage = errors.New("Invalid message error")
)

// Message is the interface implemented by an object that can analyze base form of message
// defined in protocal. And it also change itself to binary in form defined in protocal.
//
// Message is used to get message head, server will analyze message head to decide what task it should
// do, discard the message or continue to analyze its body according the message type to create InfoPkg.
type Message interface {
	base.CommunicationData
	Type() msgType
	Timestamp() time.Time
	Body() []byte
}

type msg struct {
	body      []byte
	t         msgType
	timestamp time.Time
}

// NewMessage create a new Message to wrap bytes of infos.
func NewMessage(t msgType, body []byte) Message {
	m := &msg{
		t:         t,
		timestamp: time.Now(),
		body:      body,
	}
	return m
}

// NewMessageFromBytes create a new Message from bytes.
func NewMessageFromBytes(bs []byte) (Message, error) {
	m := new(msg)
	err := m.UnmarshalBinary(bs)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Type ...
func (m *msg) Type() msgType {
	return m.t
}

// Timestamp ...
func (m *msg) Timestamp() time.Time {
	return m.timestamp
}

// Body ...
func (m *msg) Body() []byte {
	return m.body
}

// Size ...
func (m *msg) Size() int {
	return len(m.body) + msgHeadSize
}

// MarshalBinary ...
func (m *msg) MarshalBinary() ([]byte, error) {
	bs := make([]byte, m.Size())
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint32(uint32(m.Size()))
	bw.PutFloat64(float64(m.timestamp.UnixNano()))
	bw.PutUint8(uint8(m.t))

	copy(bs[msgHeadSize:], m.body)
	return bs, nil
}

// UnmarshalBinary ...
func (m *msg) UnmarshalBinary(bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)

	// first times checking message
	length := int(br.Uint32())
	if length != len(bs) {
		return ErrInvalidMessage
	}

	m.timestamp = time.Unix(0, int64(br.Float64()))
	m.t = msgType(br.Uint8())
	m.body = bs[msgHeadSize:]

	return nil
}

// CreateMessage creates instance of Message from given params, length and timestamp of the
// message will calculate automatically!
//
// This should be used to send message data to frontend
// func CreateMessage(t msgType, body []byte) Message {

// }

// NewMessage creates instance of Message from the binary.
//
// This should be used to receive message data from frontend.
// func NewMessage(message []byte) Message {

// }
