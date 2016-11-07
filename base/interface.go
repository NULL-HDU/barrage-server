package base

import (
	"encoding"
)

const (
	// InvaildBytesLenErrorFormat throw while the length of []byte parameter is not correct!
	InvaildBytesLenErrorFormat = "Invaild length of byte: hope %d, get %d."
)

// CommunicationData is the interface implemented by an object that can unmarshal a binay
// representation of itself and marshal itself into binary form.
//
// The implemented objects always are used to exchange data between client and server.
type CommunicationData interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Size provide the number of bytes after CommunicationData marshaled.
	Size() int
}

// // PoolItem is the interface implemented by an object that can clear itself.
// type PoolItem interface {
// 	Clear()
// }
