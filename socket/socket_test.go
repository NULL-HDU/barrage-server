package socket

import (
	m "barrage-server/message"
	"encoding/binary"
	"golang.org/x/net/websocket"
	"sync"
	"testing"
	"time"
)

func testWebsocketClient(testFunc func(wc *websocket.Conn)) {
	origin := "http://localhost/"
	url := "ws://localhost:2333/test"
	wc, err := websocket.Dial(url, "", origin)
	if err != nil {
		logger.Fatalln(err.Error())
	}
	testFunc(wc)
}

// TestWebsocketConnect ...
func TestWebsocketConnect(t *testing.T) {
	var w sync.WaitGroup
	w.Add(2)

	go func() {
		w.Done()
		ListenAndServer("2333", "/test")
	}()

	time.Sleep(time.Millisecond * 100)

	testFunc := func(wc *websocket.Conn) {
		var bs []byte
		if err := websocket.Message.Receive(wc, &bs); err != nil {
			logger.Errorln(err)
		}

		msg, err := m.NewMessageFromBytes(bs)
		if err != nil {
			t.Error(err.Error())
		}

		if mType := msg.Type(); mType != m.MsgRandomUserID {
			t.Errorf("Type of messsage is wrong, hope %d, get %d.", m.MsgRandomUserID, mType)
		}
		binary.BigEndian.Uint32(msg.Body())

		w.Done()
	}

	testWebsocketClient(testFunc)
	w.Wait()
}
