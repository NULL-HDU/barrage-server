package message

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"errors"
)

var (
	// ErrEmptyInfo sign the returned info empty.
	ErrEmptyInfo = errors.New("This info is empty.")
)

var logger = b.Log

type infoType msgType

// Info is a interfase used as InfoPkg body.
type Info interface {
	b.CommunicationData
}

// InfoPkg is used to transfer data among major module(user, room and playground).
type InfoPkg interface {
	Type() infoType
	Body() Info
}

type InfoListUnit interface {
	Length() uint32
	SizeOfItem() int
	Item(index int) b.CommunicationData

	NewItems(length uint32)
	Crop(length uint32)
}

// MarshalListBinary marshal InfoListUnit to bytes.
func MarshalListBinary(infolist InfoListUnit) ([]byte, error) {
	var buffer bytes.Buffer
	length := int(infolist.Length())

	var bs []byte
	var err error
	for i := 0; i < length; i++ {
		bs, err := infolist.Item(i).MarshalBinary()
		if err != nil {
			logger.Errorln(err)
			length--
			continue
		}
		buffer.Write(bs)
	}

	bw := bufbo.NewBEBufWriter(&buffer)
	bw.PutUint32(uint32(length))

	return buffer.Bytes(), nil
}

// UnmarshalListBinary unmarshal InfoListUnit from bytes.
func UnmarshalListBinary(infolist InfoListUnit, bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)
	length := br.Uint32()
	sizeOfItem := infolist.SizeOfItem()

	infolist.NewItems(length)

	unmarsLength := int(length) * sizeOfItem
	marshaledLength := 0
	count := 0
	for unmarsLength > 0 {
		err := infolist.Item(count).UnmarshalBinary(
			bs[marshaledLength : marshaledLength+sizeOfItem])
		unmarsLength -= sizeOfItem
		marshaledLength += sizeOfItem

		// ignore and drop fail marshaled bytes.
		if err != nil {
			logger.Errorln(err)
			continue
		}
		count++
	}

	length = uint32(count)
	// if infolist doesn't unmarshal any CommunicationData, throw ErrEmptyInfo.
	if length == 0 {
		return ErrEmptyInfo
	}

	// drop empty space
	infolist.Crop(length)
	return nil
}
