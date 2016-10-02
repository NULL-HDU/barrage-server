package socket

import (
	"barrage-server/base"
)

type msgType uint8

const (
	// GameOver is used when server will break off.
	GameOver msgType = 0x0b
	// SpecialMessage is used to send messages not related to game engine.
	SpecialMessage msgType = 0x0a
	// PlaygroundInfo is used when backend send balls info to frontend.
	PlaygroundInfo msgType = 0x07
	// AirplaneCreated is used to send airplane of userself to frontend while user connecting
	// into game.
	AirplaneCreated msgType = 0x06
	// SelfInfo is used when frontend send balls info to backend.
	SelfInfo msgType = 0x0c
	// Connect is used when user want to connect game.
	Connect msgType = 0x09
	// Disconnect is used when user want to leave game early.
	Disconnect msgType = 0x08
)

// Message is the interface implemented by an object that can analyze base form of message
// defined in protocal. And it also change itself to binary in form defined in protocal.
type Message interface {
	base.CommunicationData
	Len() uint32
	Type() msgType
	Timestmap() int64
	Body() []byte
}

// CreateMessage creates instance of Message from given params, length and timestamp of the
// message will calculate automatically!
//
// This should be used to send message data to frontend
func CreateMessage(t msgType, body []byte) Message {

}

// NewMessage creates instance of Message from the binary.
//
// This should be used to receive message data from frontend.
func NewMessage(message []byte) Message {

}
