package user

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"errors"
	"sync"

	"fmt"
	ws "golang.org/x/net/websocket"
	"io"
	"time"
)

var (
	// errors
	errInvalidUser   = errors.New("Invalid user error")
	errNotAllowedMsg = errors.New("Not allowed message")
	errUserID        = errors.New("User ID error")
)

var logger = b.Log
var interval = time.Second * 2

// constructErrorStringForMsg construct error string after receiving and unmarshaling message
// according to message type and running environment.
func constructErrorStringForMsg(msg m.Message, err string) string {
	if b.RunningEnv == b.Production {
		return err
	}

	if msg == nil {
		return err + ". Failed to unmarshal message."
	}

	return fmt.Sprintf("%s. Type of received message is %d.", err, msg.Type())
}

//User analyse the bytes from frontend and upload the result to Room.
//
//User will be in a room or hall until websocket of user is close fron frontend.
type User interface {
	ID() b.UserID
	Room() b.RoomID

	//Send is used by room to send bytes to frontend.
	//Send send bytes in a new goroutine.
	Send(ipkg m.InfoPkg)
	//SendError send bytes in a new goroutine.
	SendError(s string)

	//UploadInfo send infopkg to room via chan<- m.InfoPkg
	UploadInfo(infopkg m.InfoPkg) error

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
		uid:       id,
		wc:        wc,
		writeChan: make(chan []byte, 50),
	}
}

type user struct {
	nickname string
	uid      b.UserID
	wc       *ws.Conn

	stateM sync.RWMutex
	state  uint8 // 0 not play, 1 play

	roomM    sync.RWMutex
	rid      b.RoomID
	infoChan chan<- m.InfoPkg

	writeChan chan []byte
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

// UploadInfo do base check and add it to infoChan.
func (u *user) UploadInfo(ipkg m.InfoPkg) error {
	u.roomM.RLock()
	defer u.roomM.RUnlock()

	if u.infoChan == nil {
		return errInvalidUser
	}

	go func() {
		u.infoChan <- ipkg
	}()

	return nil
}

// convertBytesToInfopkg ...
func (u *user) convertBytesToInfopkg(cache []byte) (ipkg m.InfoPkg, msg m.Message, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
	}()

	msg, err = m.NewMessageFromBytes(cache)
	if err != nil {
		return nil, nil, err
	}

	ipkg, err = m.NewInfoPkgFromMsg(msg)
	if err != nil {
		return nil, msg, err
	}

	return ipkg, msg, nil
}

// preOperationForIpkg is a guard fucntion to filter invalid infopkgs and do some
// pre oreration.
func (u *user) preOperationForIpkg(ipkg m.InfoPkg) error {
	switch t := ipkg.Type(); t {
	case m.InfoPlayground:
		// add Sender to playground info
		pi := ipkg.(*m.PlaygroundInfo)
		pi.Sender = u.uid
		pi.Receiver = b.SysID
		return nil
	case m.InfoConnect:
		return u.checkConnectInfo(ipkg.Body().(*m.ConnectInfo))
	case m.InfoDisconnect:
		return u.checkDisconnectInfo(ipkg.Body().(*m.DisconnectInfo))
	default:
		return errNotAllowedMsg
	}
}

// checkConnectInfo ...
func (u *user) checkConnectInfo(ci *m.ConnectInfo) error {
	if ci.UID != u.uid {
		return errUserID
	}
	return nil
}

// checkDisconnectInfo ...
func (u *user) checkDisconnectInfo(di *m.DisconnectInfo) error {
	if di.UID != u.uid {
		return errUserID
	}
	return nil
}

// BindRoom ...
func (u *user) BindRoom(id b.RoomID, c chan<- m.InfoPkg) {
	u.roomM.Lock()
	defer u.roomM.Unlock()

	u.rid = id
	u.infoChan = c
}

// SendError ...
func (u *user) SendError(s string) {
	go func() {
		u.stateM.RLock()
		defer u.stateM.RUnlock()

		if u.state == 1 {
			u.sendError(s)
		}
	}()
}

// Send ...
func (u *user) Send(ipkg m.InfoPkg) {
	go func() {
		u.stateM.RLock()
		defer u.stateM.RUnlock()

		if u.state == 1 {
			u.sendInfoPkg(ipkg)
		}
	}()
}

// sendSpecialMessage ...
func (u *user) sendError(s string) {
	si := &m.SpecialMsgInfo{Message: s}
	u.sendInfoPkg(si)
}

// sendInfoPkg construct message and write bytes to wc.
func (u *user) sendInfoPkg(ipkg m.InfoPkg) error {
	msg, err := m.NewMessageFromInfoPkg(ipkg)
	if err != nil {
		return err
	}

	bs, _ := msg.MarshalBinary()

	u.writeChan <- bs
	return nil
}

// sendMessage ...
func (u *user) sendMessage() {
	for bs := range u.writeChan {
		u.wc.SetWriteDeadline(time.Now().Add(interval))
		if err := ws.Message.Send(u.wc, bs); err != nil {
			logger.Errorf("Can't send: %s \n", err)
		}
	}
}

// play ...
func (u *user) receiveAndUploadMessage() {
	var cache []byte
	for {
		// receive bytes
		u.wc.SetReadDeadline(time.Now().Add(interval))
		if err := ws.Message.Receive(u.wc, &cache); err != nil {
			if err != io.EOF {
				logger.Errorf("Websocket Message Receive Error: %s \n", err)
			}
			break
		}

		// convert bytes to infopkg
		ipkg, msg, err := u.convertBytesToInfopkg(cache)
		if err != nil {
			if err != m.ErrEmptyInfo {
				logger.Infof("Client Message Error: %v.\n", err)
				u.sendError(
					constructErrorStringForMsg(msg, m.ErrInvalidMessage.Error()))
			}
			continue
		}

		// pre operation for infopkg
		if err := u.preOperationForIpkg(ipkg); err != nil {
			logger.Infof("Client Message Error: %v.\n", err)
			u.sendError(constructErrorStringForMsg(msg, err.Error()))
			continue
		}

		// upload infopkg
		if err := u.UploadInfo(ipkg); err != nil {
			logger.Errorf("InfoChan of the user %d is nil.", u.ID())
			u.sendError(b.ErrServerError.Error())
			break
		}
	}
}

// overPlay ...
func (u *user) overPlay() {
	close(u.writeChan)

	u.stateM.Lock()
	defer u.stateM.Unlock()
	u.state = 0
}

// Play ...
func (u *user) Play() error {
	if u.wc == nil {
		return errInvalidUser
	}
	u.stateM.Lock()
	u.state = 1
	u.stateM.Unlock()

	go u.sendMessage()
	u.receiveAndUploadMessage()
	u.overPlay()
	return nil
}
