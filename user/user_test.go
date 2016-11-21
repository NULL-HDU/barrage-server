package user

import (
	m "barrage-server/message"
	tm "barrage-server/testLib/message"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// This example demonstrates a trivial echo server.
func createTestWebsocket(path string, testFunc func(ws *websocket.Conn)) {
	http.Handle("/"+path, websocket.Handler(testFunc))
	err := http.ListenAndServe("localhost:"+path, nil)
	if err != nil {
		logger.Panicln("ListenAndServe: " + err.Error())
	}
}

func testWebsocketClient(path string, testFunc func(ws *websocket.Conn)) {
	origin := "http://localhost/"
	url := "ws://localhost:" + path + "/" + path
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		logger.Fatalln(err)
	}
	testFunc(ws)
}

// TestUserBindRoomAndUpload ...
func TestUserBindRoomAndUpload(t *testing.T) {
	u := &user{uid: 99}
	testchan := make(chan m.InfoPkg)

	u.BindRoom(20, testchan)
	if uroom := u.Room(); uroom != 20 {
		t.Errorf("Room id of user is wrong, hope %d, get %d.", 20, uroom)
	}

	testInfopkg := &m.GameOverInfo{Overtype: 1}
	u.UploadInfo(testInfopkg)

	select {
	case ipkg := <-testchan:
		if iType := ipkg.Type(); iType != m.InfoGameOver {
			t.Errorf("Type of info package from testchan is wrong, hope %d, get %d.", iType, m.InfoGameOver)
		}
	case <-time.After(time.Millisecond * 100):
		t.Error("Conn't get infopkg from testchan.")
	}
}

// TestUserGuardMethods ...
func TestUserGuardMethods(t *testing.T) {
	var w sync.WaitGroup
	w.Add(2)

	serverCheckFunc := func(wc *websocket.Conn) {
		testchan := make(chan m.InfoPkg)
		u := &user{
			uid: 20,
			wc:  wc,
		}
		u.BindRoom(20, testchan)

		// test Play
		go func() {
			u.Play()
		}()

		time.Sleep(time.Second)

	HOPEDISCONNECT:
		for {
			select {
			case ipkg := <-testchan:
				if itype := ipkg.Type(); itype != m.InfoDisconnect {
					t.Errorf("Receive Error info! hope %d, but get %d.", m.InfoDisconnect, itype)
				}
				break HOPEDISCONNECT
			case <-time.After(time.Millisecond * 100):
			}
		}

		w.Done()
	}

	clientCheckFunc := func(wc *websocket.Conn) {
		var bs []byte
		count := 3

		invalidDi := &m.DisconnectInfo{RID: 1, UID: 10}
		msg, _ := m.NewMessageFromInfoPkg(invalidDi)
		bs, _ = msg.MarshalBinary()
		wc.Write(bs)

		invalidCi := &m.ConnectInfo{RID: 1, UID: 13}
		msg, _ = m.NewMessageFromInfoPkg(invalidCi)
		bs, _ = msg.MarshalBinary()
		wc.Write(bs)

		invalidBytes, _ := tm.GenerateTestRandomPlaygroundInfo(20, 10, 10, 10, 10).MarshalBinary()
		invalidBytes[2] = 1
		msg = m.NewMessage(m.MsgPlayground, invalidBytes)
		bs, _ = msg.MarshalBinary()
		wc.Write(bs)

		for {
			if count <= 0 {
				break
			}

			if err := websocket.Message.Receive(wc, &bs); err != nil {
				logger.Errorln(err)
			}
			count--

			msg, err := m.NewMessageFromBytes(bs)
			if err != nil {
				t.Error(err)
			}

			switch count {
			case 2:
				fallthrough
			case 1:
				if body, right := string(msg.Body()[1:]), errUserID.Error(); !strings.Contains(body, right) {
					t.Errorf("Wrong error message from server, hope contains '%s', get '%s'.", right, body)
				}
			case 0:
				if body, right := string(msg.Body()[1:]), m.ErrInvalidMessage.Error(); !strings.Contains(body, right) {
					t.Errorf("Wrong error message from server, hope contains '%s', get '%s'.", right, body)
				}
			}
		}

		di := &m.DisconnectInfo{RID: 1, UID: 20}
		bs, _ = di.MarshalBinary()
		msg = m.NewMessage(m.MsgDisconnect, bs)
		bs, _ = msg.MarshalBinary()
		wc.Write(bs)
		wc.WriteClose(0)

		w.Done()
	}

	go func() {
		createTestWebsocket("2333", serverCheckFunc)
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		testWebsocketClient("2333", clientCheckFunc)
	}()

	w.Wait()
}

// TestUserSendAndSendError ...
func TestUserSendAndSendErrorAndPlay(t *testing.T) {
	var w sync.WaitGroup
	w.Add(2)

	serverCheckFunc := func(wc *websocket.Conn) {
		testchan := make(chan m.InfoPkg)
		u := &user{
			uid: 20,
			wc:  wc,
		}
		u.BindRoom(20, testchan)

		pi := tm.GenerateTestPlaygroundInfo(0, 1, 1, 1, 1)
		// test Send
		u.sendSync(pi)
		// test SendError
		u.sendError("test_send_error")

		// test Play
		go func() {
			u.Play()
		}()

		time.Sleep(time.Second)

		select {
		case ipkg := <-testchan:
			if itype := ipkg.Type(); itype != m.InfoDisconnect {
				t.Errorf("Receive Error info! hope %d, but get %d.", m.InfoDisconnect, itype)
			}
		case <-time.After(time.Millisecond * 100):
			t.Error("Conn't get infopkg from testchan.")
		}

		time.Sleep(time.Millisecond * 200)

		w.Done()
	}

	clientCheckFunc := func(wc *websocket.Conn) {
		var bs []byte
		count := 2
		for {
			if count <= 0 {
				break
			}

			if err := websocket.Message.Receive(wc, &bs); err != nil {
				logger.Errorln(err)
			}
			count--

			msg, err := m.NewMessageFromBytes(bs)
			if err != nil {
				t.Error(err)
			}

			switch mtype := msg.Type(); mtype {
			// test SendError
			case m.MsgSpecialMessage:
				if body := string(msg.Body()[1:]); body != "test_send_error" {
					t.Errorf("Body of message should be '%s', but get '%s'.", "test_send_error", body)
				}
				// test Send
			case m.MsgPlayground:
				pi := new(m.PlaygroundInfo)
				if err := pi.UnmarshalBinary(msg.Body()); err != nil {
					t.Error(err)
				}
				if uid := pi.Sender; uid != 0 {
					t.Errorf("Sender of playgroundInfo is wrong, hope %d, get %d.", 0, uid)
				}
				if count := pi.Collisions.Length() + pi.Displacements.Length() + pi.NewBalls.Length() + len(pi.Disappears.IDs); count != 4 {
					t.Errorf("Number of info item is wrong, hope %d, get %d.", 4, count)
				}
			}
		}

		// test Play
		di := &m.DisconnectInfo{RID: 1, UID: 20}
		bs, _ = di.MarshalBinary()
		msg := m.NewMessage(m.MsgDisconnect, bs)
		bs, _ = msg.MarshalBinary()
		wc.Write(bs)

		w.Done()
	}

	go func() {
		createTestWebsocket("2222", serverCheckFunc)
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		testWebsocketClient("2222", clientCheckFunc)
	}()

	w.Wait()

}
