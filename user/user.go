package user

import (
	b "barrage-server/base"
	m "barrage-server/message"
)

//User analyse the bytes from frontend and upload the result to Room.
type User interface {
	ID() b.UserID
	Room() b.RoomID

	Name() string

	//Send is used by room to send bytes to frontend.
	//Send send bytes in a new goroutine.
	Send(bs []byte, itype m.InfoType)
	//UploadInfo send infopkg to room via chan<- m.InfoPkg
	UploadInfo(infopkg m.InfoPkg)

	//BindRoom set infopkg channel and room id for user to binds room and user.
	BindRoom(id b.RoomID, c chan<- m.InfoPkg)
}
