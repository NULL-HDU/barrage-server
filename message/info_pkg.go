package message

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
)

// DisconnectInfo is used while user disconnect game early, holding data about
// user id and the id of that user's room.
type DisconnectInfo struct {
	// type of info
	t infoType

	uid b.UserID
	rid b.RoomID
}

// Type return type of information
func (di *DisconnectInfo) Type() infoType {
	return di.t
}

// Body return DisconnectInfo self.
func (di *DisconnectInfo) Body() Info {
	return di
}

// MarshalBinary marshal DisconnectInfo to bytes
func (di *DisconnectInfo) MarshalBinary() ([]byte, error) {
	bs := make([]byte, 12)
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint64(uint64(di.uid))
	bw.PutUint32(uint32(di.rid))

	return bs, nil
}

// UnmarshalBinary unmarshal DisconnectInfo from bytes
func (di *DisconnectInfo) UnmarshalBinary(bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)

	di.uid = b.UserID(br.Uint64())
	di.rid = b.RoomID(br.Uint32())

	return nil
}
