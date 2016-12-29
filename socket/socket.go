package socket

import (
	b "barrage-server/base"
	m "barrage-server/message"
	r "barrage-server/room"
	"barrage-server/user"
	"errors"
	ws "golang.org/x/net/websocket"
	"net/http"
	"regexp"
)

var logger = b.Log

// ListenAndServer open a server.
func ListenAndServer(port, path string) {
	s := new(socket)
	if err := s.Open(port, path); err != nil {
		logger.Errorln(err)
	}
}

// Socket create a http server, and wapper websocket connect into
// user.
type Socket interface {
	// create a http server listen on 0.0.0.0:port, server websocket for path.
	Open(port, path string) error
	// handle websocket connnect.
	HandleFunc(wc *ws.Conn)
}

type socket struct{}

// Open ...
func (s *socket) Open(port, path string) error {
	// check if port and path are valid.
	matched, err := regexp.MatchString("^\\d{1,5}$", port)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Invalid parameter 'port'.")
	}

	matched, err = regexp.MatchString("^/\\w+$", path)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Invalid parameter 'path'.")
	}

	// provide websocket server
	http.Handle(path, ws.Handler(s.HandleFunc))

	logger.Infof("Service start, bind port: %v \n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalln("ListenAndServe:", err)
	}

	return nil
}

// HandleFunc ...
func (s *socket) HandleFunc(wc *ws.Conn) {
	logger.Infof("Connect flow service from %v \n", wc.RemoteAddr())

	defer func() {
		wc.Close()
		if v := recover(); v != nil {
			logger.Infof("HandleFunc Panic! %v\n", v)
		}
	}()

	// Random user Id (s -> c)
	msg, uid := m.NewRandomUserIDMsg()
	logger.Infof("random uid %d \n", uid)
	bs, _ := msg.MarshalBinary()
	if err := ws.Message.Send(wc, bs); err != nil {
		logger.Errorf("Can't send randID: %s \n", err)
		return
	}

	u := user.NewUser(wc, uid)
	logger.Infoln("user create success.")
	r.JoinHall(u)
	logger.Infoln("user start play.")
	u.Play()
	r.LeftHall(uid)

	logger.Infof("User %d left game. \n", uid)
	logger.Infof("Close Connect from %v \n", wc.RemoteAddr())
}
