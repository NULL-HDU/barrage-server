package user

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"errors"
	"sync"

	ws "golang.org/x/net/websocket"
	"io"
)

var (
	// errors
	errInvalidUser = errors.New("Invalid user error.")
)

var logger = b.Log

//User analyse the bytes from frontend and upload the result to Room.
//
//User will be in a room or hall until websocket of user is close fron frontend.
type User interface {
	ID() b.UserID
	Room() b.RoomID

	//Send is used by room to send bytes to frontend.
	//Send send bytes in a new goroutine.
	Send(bs []byte, itype m.InfoType)
	//SendError send bytes in a new goroutine.
	SendError(s string)

	//UploadInfo send infopkg to room via chan<- m.InfoPkg
	UploadInfo(infopkg m.InfoPkg)

	//BindRoom set infopkg channel and room id for user to binds room and user.
	BindRoom(id b.RoomID, c chan<- m.InfoPkg)

	// socket package should call Play to ready user(listen messages)
	// before call Play, socket should join user into hall, after call Play, socket
	// should left user from hall.
	Play() error
}

// NewUser create a User by websocket.Conn and userID.
func NewUser(wc *ws.Conn, id b.UserID) User {
	return &user{
		uid: id,
		wc:  wc,
	}
}

type user struct {
	nickname string
	uid      b.UserID
	wc       *ws.Conn

	roomM    sync.RWMutex
	rid      b.RoomID
	infoChan chan<- m.InfoPkg
}

// ID ...
func (u *user) ID() b.UserID {
	return u.uid
}

// Room ...
func (u *user) Room() b.RoomID {
	u.roomM.RLock()
	defer u.roomM.RUnlock()

	return u.rid
}

// UploadInfo ...
func (u *user) UploadInfo(ipkg m.InfoPkg) {
	go func() {
		u.roomM.RLock()
		defer u.roomM.RUnlock()

		if u.infoChan == nil {
			logger.Errorln(errInvalidUser)
			return
		}
		u.infoChan <- ipkg
	}()
}

// SendError ...
func (u *user) SendError(s string) {
	go func() {
		u.sendError(s)
	}()
}

// Send ...
func (u *user) Send(bs []byte, itype m.InfoType) {
	go func() {
		u.sendSync(bs, itype)
	}()
}

// BindRoom ...
func (u *user) BindRoom(id b.RoomID, c chan<- m.InfoPkg) {
	u.roomM.Lock()
	defer u.roomM.Unlock()

	u.rid = id
	u.infoChan = c
}

// sendSpecialMessage ...
func (u *user) sendError(s string) {
	si := &m.SpecialMsgInfo{Message: s}
	bs, _ := si.MarshalBinary()
	u.sendSync(bs, m.InfoSpecialMessage)
}

// sendSync construct message and write bytes to wc.
func (u *user) sendSync(body []byte, itype m.InfoType) error {
	msg, err := m.NewMessageFromInfo(itype, body)
	if err != nil {
		logger.Errorln(err)
	}

	bs, _ := msg.MarshalBinary()

	if err := ws.Message.Send(u.wc, bs); err != nil {
		logger.Errorf("Can't send: %s \n", err)
	}

	return nil
}

// play ...
func (u *user) recieveAndUploadMessage() {
	var cache []byte
	for {
		if err := ws.Message.Receive(u.wc, &cache); err != nil {
			if err != io.EOF {
				u.sendError(b.ErrServerError.Error())
				logger.Errorf("Error: %s \n", err)
			}
			break
		}

		msg, err := m.NewMessageFromBytes(cache)
		if err != nil {
			u.sendError("Invalid Message.")
			logger.Errorf("User %d Error: %s \n", u.uid, err)
		}

		ipkg, err := m.NewInfoPkgFromMsg(msg)
		if err != nil {
			// ignore empty info.
			if err == m.ErrEmptyInfo {
				continue
			} else {
				u.sendError("Invalid Message.")
				logger.Errorf("User %d Error: %s \n", u.uid, err)
			}
		}

		u.UploadInfo(ipkg)
	}

	logger.Infof("Close Connect from %v \n", u.wc.RemoteAddr())
}

// Play ...
func (u *user) Play() error {
	if u.wc == nil {
		return errInvalidUser
	}

	u.recieveAndUploadMessage()
	return nil
}
