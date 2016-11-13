package message

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var logger = b.Log

// MsgType type for message
type MsgType uint8

const (
	// backend -> frontend

	// MsgRandomUserID is used when websocket connect is created.
	MsgRandomUserID MsgType = 0xd4

	// MsgGameOver is used when server will break off.
	MsgGameOver MsgType = 0x0b
	// MsgSpecialMessage is used to send messages not related to game engine.
	MsgSpecialMessage MsgType = 0x0a
	// MsgPlayground is used when backend send balls info to frontend.
	MsgPlayground MsgType = 0x07
	// MsgAirplaneCreated is used to send airplane of userself to frontend while user connecting
	// into game.
	MsgAirplaneCreated MsgType = 0x06

	// frontend -> backend

	// MsgUserSelf is used when frontend send balls info to backend.
	MsgUserSelf MsgType = 0x0c
	// MsgConnect is used when user want to connect game.
	MsgConnect MsgType = 0x09
	// MsgDisconnect is used when user want to leave game early.
	MsgDisconnect MsgType = 0x08
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

var infoMsgSendMap = map[InfoType]MsgType{
	InfoGameOver:        MsgGameOver,
	InfoSpecialMessage:  MsgSpecialMessage,
	InfoAirplaneCreated: MsgAirplaneCreated,
	InfoPlayground:      MsgPlayground,
	InfoConnect:         MsgConnect,
	InfoDisconnect:      MsgDisconnect,
}

// Message is the interface implemented by an object that can analyze base form of message
// defined in protocal. And it also change itself to binary in form defined in protocal.
//
// Message is used to get message head, server will analyze message head to decide what task it should
// do, discard the message or continue to analyze its body according the message type to create InfoPkg.
type Message interface {
	b.CommunicationData
	Type() MsgType
	Timestamp() time.Time
	Body() []byte
}

type msg struct {
	body      []byte
	t         MsgType
	timestamp time.Time
}

// NewMessageFromInfoPkg creates instance of Message from given InfoPkg.
//
// This should be used to send message data to frontend
func NewMessageFromInfoPkg(ipkg InfoPkg) (Message, error) {
	iType := ipkg.Type()
	body := ipkg.Body()

	bs, err := body.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("Info Marshal Error: %s.", err)
	}

	mType, ok := infoMsgSendMap[iType]
	if !ok {
		return nil, fmt.Errorf("Not found mapped message for the infoType(%v).", iType)
	}

	return NewMessage(mType, bs), nil
}

// NewMessage creates instance of Message from given params.
// length and timestamp of the message will calculate automatically!
//
// This should be used to send message data to frontend
func NewMessage(t MsgType, body []byte) Message {
	m := &msg{
		t:         t,
		timestamp: time.Now(),
		body:      body,
	}
	return m
}

// NewMessageFromBytes create a new Message from bytes.
//
// This should be used to receive message data from frontend.
func NewMessageFromBytes(bs []byte) (Message, error) {
	m := new(msg)
	err := m.UnmarshalBinary(bs)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Type ...
func (m *msg) Type() MsgType {
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
	m.t = MsgType(br.Uint8())
	m.body = bs[msgHeadSize:]

	return nil
}

// NewRandomUserIDMsg create a new message whose type is MsgRandomUserID containing a randomID.
func NewRandomUserIDMsg() (Message, b.UserID) {
	randID := uint32(rand.Int31())

	bs := make([]byte, 4)
	bw := bufbo.NewBEBytesWriter(bs)
	bw.PutUint32(randID)

	return NewMessage(MsgRandomUserID, bs), b.UserID(randID)
}
