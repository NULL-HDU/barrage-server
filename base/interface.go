package base

import (
	"encoding"
)

// CommunicationData is the interface implemented by an object that can unmarshal a binay
// representation of itself and marshal itself into binary form.
//
// The implemented objects always are used to exchange data between client and server.
type CommunicationData interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
