package message

import (
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
)

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
