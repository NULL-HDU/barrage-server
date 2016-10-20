package message

import (
	"barrage-server/base"
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

// Message is the interface implemented by an object that can analyze base form of message
// defined in protocal. And it also change itself to binary in form defined in protocal.
//
// Message is used to get message head, server will analyze message head to decide what task it should
// do, discard the message or continue to analyze its body according the message type to create InfoPkg.
type Message interface {
	base.CommunicationData
	Len() int
	Type() msgType
	Timestamp() int64
	Body() []byte
}

type msg struct {
	body      []byte
	t         msgType
	timestamp int64
}

// Len ...
func (m *msg) Len() int {
	return len(m.body) + 4 + 1 + 8
}

// Type ...
func (m *msg) Type() msgType {
	return m.t
}

// Timestamp ...
func (m *msg) Timestamp() int64 {
	return m.timestamp
}

// Body ...
func (m *msg) Body() []byte {
	return m.body
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
