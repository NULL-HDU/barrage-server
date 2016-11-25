// Package provide a test service for socket of frontend to test game protocal.
package main

import (
	"barrage-server/ball"
	"barrage-server/libs/log"
	"barrage-server/message"
	"barrage-server/testService/data"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"math"
	"net/http"
	"time"
)

var logger log.Logger

func init() {
	logger = log.NewStdLogger(log.InfoLevel)
}

func baseWs(ws *websocket.Conn) {
	logger.Infof("Connect base service from %v \n", ws.RemoteAddr())

	if err := websocket.Message.Send(ws, data.RandomUserID()); err != nil {
		logger.Errorf("Can't send: %s \n", err)
	}

	var cache []byte
	for {
		if err := websocket.Message.Receive(ws, &cache); err != nil {
			if err != io.EOF {
				logger.Errorf("Error: %s \n", err)
			}
			break
		}
		logger.Infof("Receive: % x \n", cache)

		if err := websocket.Message.Send(ws, cache); err != nil {
			logger.Errorf("Can't send: %s \n", err)
			break
		}
	}

	logger.Infof("Close Connect from %v \n", ws.RemoteAddr())
	ws.Close()
}

// echo response a test binary to client whenever client connect to server
func echo(ws *websocket.Conn) {
	logger.Infof("Connect echo service from %v \n", ws.RemoteAddr())

	var cache []byte
	for {
		if err := websocket.Message.Receive(ws, &cache); err != nil {
			if err != io.EOF {
				logger.Errorf("Error: %s \n", err)
			}
			break
		}
		logger.Infof("Receive: % x \n", cache)

		if err := websocket.Message.Send(ws, cache); err != nil {
			logger.Errorf("Can't send: %s \n", err)
			break
		}
	}

	logger.Infof("Close Connect from %v \n", ws.RemoteAddr())
	ws.Close()
}

// flow run a test service for flow testing.
func flow(ws *websocket.Conn) {
	logger.Infof("Connect flow service from %v \n", ws.RemoteAddr())

	flowStep := 0
	var cache []byte
	var err error

	// step 0
	// Random user Id (s -> c)
	if err = websocket.Message.Send(ws, data.RandomUserID()); err != nil {
		logger.Errorf("Can't send: %s \n", err)
	}
	flowStep++

FLOWOVER:
	for {
		if err = websocket.Message.Receive(ws, &cache); err != nil {
			if err != io.EOF {
				logger.Errorf("Error: %s \n", err)
			}
			break
		}

		// message parse check
		msg, err := message.NewMessageFromBytes(cache)
		if err != nil {
			logger.Infof("Clinet message error: %v. \n", err)
			sendSpecialMessage(ws, err.Error())
			break
		}

		// timestamp check
		ts, nowTs := msg.Timestamp().UnixNano(), time.Now().UnixNano()
		if math.Abs(float64(ts-nowTs)) > float64(600*time.Millisecond) {
			errString := "Timestamp error: the bias of Timestamp is bigger than 600ms."
			logger.Infof("Clinet message error: %v. \n", errString)
			sendSpecialMessage(ws, errString)
			break
		}

		switch flowStep {
		case 1: //(c→ s)[connect] → (s→ c)[connected]
			err = checkConnectMsgAndSendConnected(ws, msg)
		case 2: //(c→ s)[self info] → (s→ c)[playground info]
			err = checkSelfInfoMsgAndSendPlaygroundInfoMsg(ws, msg)
		case 3: //(c→ s)[disconnect] → (s→ c)[game over]
			err = checkDisconnectMsgAndSendGameOverMsg(ws, msg)
		default:
			logger.Infoln("Flow test success!")
			break FLOWOVER
		}

		if err != nil {
			logger.Infof("Clinet message error: %v. \n", err)
			sendSpecialMessage(ws, err.Error())
			break
		}
		flowStep++

	}

	logger.Infof("Close Connect from %v \n", ws.RemoteAddr())
	ws.Close()
}

func checkConnectMsgAndSendConnected(ws *websocket.Conn, msg message.Message) error {
	if msgType := msg.Type(); msgType != message.MsgConnect {
		return fmt.Errorf(
			"Message Type error: hope get connect message(type %v), but get %v.",
			message.MsgConnect, msgType)
	}

	ci := new(message.ConnectInfo)
	err := ci.UnmarshalBinary(msg.Body())
	if err != nil {
		return err
	}

	cedi := &message.ConnectedInfo{UID: ci.UID, RID: ci.RID}
	bs, err := cedi.MarshalBinary()
	if err != nil {
		return err
	}

	m := message.NewMessage(message.MsgConnected, bs)
	bs, err = m.MarshalBinary()
	if err != nil {
		return err
	}

	if err := websocket.Message.Send(ws, bs); err != nil {
		return err
	}

	return nil
}

func checkSelfInfoMsgAndSendPlaygroundInfoMsg(ws *websocket.Conn, msg message.Message) error {
	if msgType := msg.Type(); msgType != message.MsgUserSelf {
		return fmt.Errorf(
			"Message Type error: hope get connect message(type %v), but get %v.",
			message.MsgUserSelf, msgType)
	}

	pi := new(message.PlaygroundInfo)
	err := pi.UnmarshalBinary(msg.Body())
	if err != nil {
		return err
	}

	bs, err := pi.MarshalBinary()
	if err != nil {
		return err
	}

	m := message.NewMessage(message.MsgGameOver, bs)
	bs, err = m.MarshalBinary()
	if err != nil {
		return err
	}

	if err := websocket.Message.Send(ws, bs); err != nil {
		return err
	}

	return nil
}

func checkDisconnectMsgAndSendGameOverMsg(ws *websocket.Conn, msg message.Message) error {
	if msgType := msg.Type(); msgType != message.MsgDisconnect {
		return fmt.Errorf(
			"Message Type error: hope get connect message(type %v), but get %v.",
			message.MsgDisconnect, msgType)
	}

	ci := new(message.DisconnectInfo)
	err := ci.UnmarshalBinary(msg.Body())
	if err != nil {
		return err
	}

	goi := &message.GameOverInfo{Overtype: 1}
	bs, err := goi.MarshalBinary()
	if err != nil {
		return err
	}

	m := message.NewMessage(message.MsgGameOver, bs)
	bs, err = m.MarshalBinary()
	if err != nil {
		return err
	}

	if err := websocket.Message.Send(ws, bs); err != nil {
		return err
	}

	return nil
}

func sendSpecialMessage(ws *websocket.Conn, msg string) {
	smi := &message.SpecialMsgInfo{Message: msg}
	bs, err := smi.MarshalBinary()
	if err != nil {
		logger.Errorln(err)
		return
	}

	m := message.NewMessage(message.MsgSpecialMessage, bs)
	bs, err = m.MarshalBinary()
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err := websocket.Message.Send(ws, bs); err != nil {
		logger.Errorf("Can't send: %s \n", err)
	}
}

func main() {
	// provide file server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// provide websocket server
	http.Handle("/ws", websocket.Handler(baseWs))
	http.Handle("/echo", websocket.Handler(echo))
	http.Handle("/flow", websocket.Handler(flow))
	port := "1234"

	logger.Infof("Service start, bind port: %v \n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalln("ListenAndServe:", err)
	}
}
