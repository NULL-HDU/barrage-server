package message

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"errors"
	"fmt"
	"math"
)

var (
	// ErrEmptyInfo sign the returned info empty.
	ErrEmptyInfo = errors.New("This info is empty.")
)

// InfoType type of info
type InfoType MsgType

const (
	// Room -> User -----------------------------------------------------------

	// InfoGameOver is used when server will break off.
	InfoGameOver = InfoType(iota)
	// InfoSpecialMessage is used to send informations not related to game engine.
	InfoSpecialMessage
	// InfoAirplaneCreated is used to send airplane of userself to frontend while user connecting
	// into game.
	InfoAirplaneCreated
	// InfoConnected is used when room done the connect for user.
	InfoConnected

	// User -> Room -----------------------------------------------------------

	// InfoConnect is used when user want to connect game.
	InfoConnect
	// InfoDisconnect is used when user want to leave game early(now, just leave room to hall).
	InfoDisconnect

	// Room -> User, User -> Room -------------------------

	// InfoPlayground is used when backend send balls info to frontend.
	InfoPlayground
)

// Info is a interfase used as InfoPkg body.
type Info interface {
	b.CommunicationData
}

// InfoPkg is used to transfer data among major module(user, room and playground).
type InfoPkg interface {
	Type() InfoType
	Body() Info
}

// InfoList is used to marshal and unmarshal struct which contains length
// and a b.CommunicationData list.
type InfoList interface {
	Length() int
	Item(index int) b.CommunicationData

	NewItems(length uint32)
	Crop(length uint32)
}

// NewInfoPkgFromMsg ...
func NewInfoPkgFromMsg(msg Message) (InfoPkg, error) {
	var ipkg InfoPkg
	switch t := msg.Type(); t {
	case MsgSpecialMessage:
		ipkg = &SpecialMsgInfo{}
	case MsgDisconnect:
		ipkg = &DisconnectInfo{}
	case MsgConnect:
		ipkg = &ConnectInfo{}
	case MsgConnected:
		ipkg = &ConnectedInfo{}
	case MsgGameOver:
		ipkg = &GameOverInfo{}
	case MsgPlayground:
		fallthrough
	case MsgUserSelf:
		ipkg = &PlaygroundInfo{}
	default:
		return nil, fmt.Errorf("Not found mapped infopkg for the message(%v).", t)
	}

	if err := ipkg.(Info).UnmarshalBinary(msg.Body()); err != nil {
		return nil, err
	}
	return ipkg, nil
}

// MarshalListBinary marshal InfoList to bytes.
func MarshalListBinary(infolist InfoList) ([]byte, error) {
	var buffer bytes.Buffer
	// empty bytes of length first
	buffer.Write([]byte{0, 0, 0, 0})

	length := infolist.Length()
	if length > math.MaxUint32 {
		return nil, fmt.Errorf("Too many infos, get %d.", length)
	}

	count := length
	var bs []byte
	var err error
	for i := 0; i < length; i++ {
		bs, err = infolist.Item(i).MarshalBinary()
		if err != nil {
			logger.Errorln(err)
			count--
			continue
		}
		buffer.Write(bs)
	}

	// finally, write the length of list into begin of result bytes.
	result := buffer.Bytes()
	bw := bufbo.NewBEBytesWriter(result)
	bw.PutUint32(uint32(count))

	return result, nil
}

// UnmarshalListBinary unmarshal InfoList from bytes.
func UnmarshalListBinary(infolist InfoList, bs []byte) (unmarshaledLength int, err error) {
	br := bufbo.NewBEBytesReader(bs)
	length := br.Uint32()

	infolist.NewItems(length)

	unmarshaledLength = 4
	count := uint32(0)
	for i := uint32(0); i < length; i++ {
		item := infolist.Item(int(count))
		err := item.UnmarshalBinary(bs[unmarshaledLength:])
		size := item.Size()
		unmarshaledLength += size

		// ignore and drop fail marshaled bytes.
		if err != nil {
			logger.Errorln(err)
			continue
		}
		count++
	}

	length = count
	// if infolist doesn't unmarshal any CommunicationData, throw ErrEmptyInfo.
	if length == 0 {
		return unmarshaledLength, ErrEmptyInfo
	}

	// drop empty space
	infolist.Crop(length)
	return unmarshaledLength, nil
}
