package message

import (
	"barrage-server/ball"
	b "barrage-server/base"
	"barrage-server/libs/bufbo"
	"bytes"
	"fmt"
	"math"
)

// GameOverInfo send information from Room to User while server gonna shutdown.
type GameOverInfo struct {
	Overtype uint8
}

// Type return type of information
func (goi *GameOverInfo) Type() InfoType {
	return InfoGameOver
}

// Body return GameOverInfo self.
func (goi *GameOverInfo) Body() Info {
	return goi
}

// Size return the number of bytes after marshaled.
func (goi *GameOverInfo) Size() int {
	return 1
}

// MarshalBinary marshal GameOverInfo to bytes
func (goi *GameOverInfo) MarshalBinary() ([]byte, error) {
	bs := []byte{goi.Overtype}
	return bs, nil
}

// UnmarshalBinary unmarshal GameOverInfo from bytes
func (goi *GameOverInfo) UnmarshalBinary(bs []byte) error {
	goi.Overtype = bs[0]
	return nil
}

// SpecialMsgInfo send information from Room to User while special message generated.
type SpecialMsgInfo struct {
	Message string
}

// Type return type of information
func (smi *SpecialMsgInfo) Type() InfoType {
	return InfoSpecialMessage
}

// Body return SpecialMsgInfo self.
func (smi *SpecialMsgInfo) Body() Info {
	return smi
}

// Size return the number of bytes after marshaled.
func (smi *SpecialMsgInfo) Size() int {
	return 1 + len(smi.Message)
}

// MarshalBinary marshal SpecialMsgInfo to bytes
func (smi *SpecialMsgInfo) MarshalBinary() ([]byte, error) {
	msgBytes := []byte(smi.Message)
	msgLen := smi.Size() - 1

	if msgLen > math.MaxUint8 {
		return nil, fmt.Errorf("SpecialMsgInfo MarshalError: Special message is too long, hope 255, get %d.", msgLen)
	}

	bs := make([]byte, msgLen+1)
	bs[0] = uint8(msgLen)
	copy(bs[1:], msgBytes)

	return bs, nil
}

// UnmarshalBinary unmarshal SpecialMsgInfo from bytes
func (smi *SpecialMsgInfo) UnmarshalBinary(bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)

	smi.Message = br.Str(int(br.Uint8()))

	return nil
}

// AirplaneCreatedInfo send information from Room to User after user connecting to server.
type AirplaneCreatedInfo struct {
	Airplane ball.Ball
}

// Type return type of information
func (aci *AirplaneCreatedInfo) Type() InfoType {
	return InfoAirplaneCreated
}

// Body return AirplaneCreatedInfo self.
func (aci *AirplaneCreatedInfo) Body() Info {
	return aci
}

// Size return the number of bytes after marshaled.
func (aci *AirplaneCreatedInfo) Size() int {
	return aci.Airplane.Size()
}

// MarshalBinary marshal AirplaneCreatedInfo to bytes
func (aci *AirplaneCreatedInfo) MarshalBinary() ([]byte, error) {
	bs, err := aci.Airplane.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("AirplaneCreatedInfo MarshalError: %v", err)
	}
	return bs, nil
}

// UnmarshalBinary unmarshal AirplaneCreatedInfo from bytes
func (aci *AirplaneCreatedInfo) UnmarshalBinary(bs []byte) error {
	aci.Airplane = ball.NewBall()
	return aci.Airplane.UnmarshalBinary(bs)
}

// DisconnectInfo send information from User to Room while user disconnecting
// game early, holding data about user id and the id of that user's room.
type DisconnectInfo struct {
	// type of info
	// t InfoType

	UID b.UserID
	RID b.RoomID
}

// Type return type of information
func (di *DisconnectInfo) Type() InfoType {
	return InfoDisconnect
}

// Body return DisconnectInfo self.
func (di *DisconnectInfo) Body() Info {
	return di
}

// Size return the number of bytes after marshaled.
func (di *DisconnectInfo) Size() int {
	return 8
}

// MarshalBinary marshal DisconnectInfo to bytes
func (di *DisconnectInfo) MarshalBinary() ([]byte, error) {
	bs := make([]byte, di.Size())
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint32(uint32(di.UID))
	bw.PutUint32(uint32(di.RID))

	return bs, nil
}

// UnmarshalBinary unmarshal DisconnectInfo from bytes
func (di *DisconnectInfo) UnmarshalBinary(bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)

	di.UID = b.UserID(br.Uint32())
	di.RID = b.RoomID(br.Uint32())

	return nil
}

// ConnectInfo send information from User to Room while user joining
// game.
type ConnectInfo struct {
	UID      b.UserID
	Nickname string
	RID      b.RoomID
	Troop    uint8
}

// Type return type of information
func (ci *ConnectInfo) Type() InfoType {
	return InfoConnect
}

// Body return ConnectInfo self.
func (ci *ConnectInfo) Body() Info {
	return ci
}

// Size return the number of bytes after marshaled.
func (ci *ConnectInfo) Size() int {
	return 10 + len(ci.Nickname)
}

// MarshalBinary marshal ConnectInfo to bytes
func (ci *ConnectInfo) MarshalBinary() ([]byte, error) {
	bs := make([]byte, ci.Size())
	bw := bufbo.NewBEBytesWriter(bs)

	bw.PutUint32(uint32(ci.UID))
	bw.PutUint8(uint8(len(ci.Nickname)))
	bw.PutStr(ci.Nickname)
	bw.PutUint32(uint32(ci.RID))
	bw.PutUint8(ci.Troop)

	return bs, nil
}

// UnmarshalBinary unmarshal ConnectInfo from bytes
func (ci *ConnectInfo) UnmarshalBinary(bs []byte) error {
	br := bufbo.NewBEBytesReader(bs)

	ci.UID = b.UserID(br.Uint32())
	ci.Nickname = br.Str(int(br.Uint8()))
	ci.RID = b.RoomID(br.Uint32())
	ci.Troop = br.Uint8()

	return nil
}

// PlaygroundInfo exchange informations among User, Room and Playground.
type PlaygroundInfo struct {
	Sender     b.UserID
	Reciever   b.UserID
	CacheBytes []byte

	NewBalls      *BallsInfo
	Displacements *BallsInfo
	Collisions    *CollisionsInfo
	Disappears    *DisappearsInfo
}

// Type return type of information
func (pi *PlaygroundInfo) Type() InfoType {
	return InfoPlayground
}

// Body return PlaygroundInfo self.
func (pi *PlaygroundInfo) Body() Info {
	return pi
}

// Size return the number of bytes after marshaled.
func (pi *PlaygroundInfo) Size() int {
	return pi.Collisions.Size() + pi.Displacements.Size() + pi.NewBalls.Size() + pi.Disappears.Size()
}

// MarshalBinary marshal PlaygroundInfo to bytes
func (pi *PlaygroundInfo) MarshalBinary() ([]byte, error) {
	if pi.CacheBytes != nil {
		return pi.CacheBytes, nil
	}

	var buffer bytes.Buffer

	// NewBalls
	bs, err := MarshalListBinary(pi.NewBalls)
	if err != nil {
		return nil, fmt.Errorf("PlaygroundInfo MarshalError: %v", err)
	}
	buffer.Write(bs)

	// Displacements
	bs, err = MarshalListBinary(pi.Displacements)
	if err != nil {
		return nil, fmt.Errorf("PlaygroundInfo MarshalError: %v", err)
	}
	buffer.Write(bs)

	// Collisions
	bs, err = MarshalListBinary(pi.Collisions)
	if err != nil {
		return nil, fmt.Errorf("PlaygroundInfo MarshalError: %v", err)
	}
	buffer.Write(bs)

	// Disappears
	bs, err = pi.Disappears.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("PlaygroundInfo MarshalError: %v", err)
	}
	buffer.Write(bs)

	pi.CacheBytes = buffer.Bytes()
	return pi.CacheBytes, nil
}

// UnmarshalBinary unmarshal PlaygroundInfo from bytes
func (pi *PlaygroundInfo) UnmarshalBinary(bs []byte) error {
	pi.NewBalls = new(BallsInfo)
	pi.Displacements = new(BallsInfo)
	pi.Collisions = new(CollisionsInfo)
	pi.Disappears = new(DisappearsInfo)
	validPartsNum := 4

	length := 0
	n, err := UnmarshalListBinary(pi.NewBalls, bs[length:])
	if err != nil {
		if err == ErrEmptyInfo {
			validPartsNum--
		} else {
			return fmt.Errorf("PlaygroundInfo UnmarshalError: %v", err)
		}
	}

	length += n
	n, err = UnmarshalListBinary(pi.Displacements, bs[length:])
	if err != nil {
		if err == ErrEmptyInfo {
			validPartsNum--
		} else {
			return fmt.Errorf("PlaygroundInfo UnmarshalError: %v", err)
		}
	}

	length += n
	n, err = UnmarshalListBinary(pi.Collisions, bs[length:])
	if err != nil {
		if err == ErrEmptyInfo {
			validPartsNum--
		} else {
			return fmt.Errorf("PlaygroundInfo UnmarshalError: %v", err)
		}
	}

	length += n
	err = pi.Disappears.UnmarshalBinary(bs[length:])
	if err != nil {
		if err == ErrEmptyInfo {
			validPartsNum--
		} else {
			return fmt.Errorf("PlaygroundInfo UnmarshalError: %v", err)
		}
	}

	// Empty PlaygroundInfo should be drop.
	if validPartsNum == 0 {
		return ErrEmptyInfo
	}

	return nil
}
