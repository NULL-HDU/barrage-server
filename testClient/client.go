package main

import (
	b "barrage-server/base"
	"barrage-server/libs/cmdface"
	m "barrage-server/message"
	tm "barrage-server/testLib/message"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"strconv"
)

var ipkgsLinkList *infoPkgNode
var tailOfLinkList *infoPkgNode
var websocketConn *websocket.Conn
var infoTypeMap = map[m.InfoType]string{
	m.InfoDisconnect:      "disconnect info",
	m.InfoAirplaneCreated: "airplanecreated info",
	m.InfoGameOver:        "gameover info",
	m.InfoPlayground:      "playground info",
	m.InfoSpecialMessage:  "specialmessage info",
	m.InfoConnect:         "connect info",
}
var uid b.UserID

var isWsConnected = false

type infoPkgNode struct {
	ipkg m.InfoPkg
	next *infoPkgNode
}

func pushInfoPkg(ipkg m.InfoPkg) {
	newNode := &infoPkgNode{
		ipkg: ipkg,
	}
	if ipkgsLinkList == nil {
		ipkgsLinkList = newNode
		tailOfLinkList = ipkgsLinkList
	} else {
		tailOfLinkList.next = newNode
		tailOfLinkList = tailOfLinkList.next
	}
}

func connToServer(port int, path string) error {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://localhost:%d%s", port, path)
	wc, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	websocketConn = wc
	isWsConnected = true

	go func() {
		var bs []byte
		for {
			if err := websocket.Message.Receive(wc, &bs); err != nil {
				if err == io.EOF {
					break
				}
				cmdface.Show(err.Error())
			}

			msg, err := m.NewMessageFromBytes(bs)
			if err != nil {
				cmdface.Show(err.Error())
			}

			if msg.Type() == m.MsgRandomUserID {
				uid = b.UserID(binary.BigEndian.Uint32(msg.Body()))
				continue
			}

			ipkg, err := m.NewInfoPkgFromMsg(msg)
			if err != nil {
				cmdface.Show(err.Error())
			}
			pushInfoPkg(ipkg)
		}
	}()

	return nil
}

func closeConnect() {
	if websocketConn != nil {
		websocketConn.Close()
		isWsConnected = false
	}
}

func cleanInfoPkgs() {
	ipkgsLinkList = nil
}

func sendMessage(ipkg m.InfoPkg) error {
	msg, err := m.NewMessageFromInfoPkg(ipkg)
	if err != nil {
		return err
	}

	bs, _ := msg.MarshalBinary()
	websocketConn.Write(bs)
	return nil
}

func sendConnectInfo(rid b.RoomID, name string) error {
	ci := &m.ConnectInfo{
		UID:      uid,
		Nickname: name,
		RID:      rid,
		Troop:    0,
	}

	return sendMessage(ci)
}

func sendDisconnectInfo(rid b.RoomID) error {
	di := &m.DisconnectInfo{
		UID: uid,
		RID: rid,
	}

	return sendMessage(di)
}

func sendPlaygroundInfo(cin, din, nin, dsin int) error {
	pi := tm.GenerateTestPlaygroundInfo(uid, nin, din, cin, dsin)

	return sendMessage(pi)
}

func showInfoPkgList() {
	ipkgNode := ipkgsLinkList
	count := 0
	for ipkgNode != nil {
		ipkg := ipkgNode.ipkg
		cmdface.Show(fmt.Sprintf("[%d] %s\n", count, infoTypeMap[ipkg.Type()]))
		count++
		ipkgNode = ipkgNode.next
	}
}

func showInfoPkgListFunc(params []string) {
	showInfoPkgList()
}

func showUidFunc(params []string) {
	cmdface.Show(fmt.Sprintf("uid: %d.\n", uid))
}

func cleanInfoPkgListFunc(params []string) {
	cleanInfoPkgs()
}

func sendConnectInfoFunc(params []string) {
	rid, err := strconv.Atoi(params[0])
	if err != nil {
		cmdface.Show(err.Error())
	}
	if err = sendConnectInfo(b.RoomID(rid), params[1]); err != nil {
		cmdface.Show(err.Error())
	}
}

func sendDisconnectInfoFunc(params []string) {
	rid, err := strconv.Atoi(params[0])
	if err != nil {
		cmdface.Show(err.Error())
	}
	if err = sendDisconnectInfo(b.RoomID(rid)); err != nil {
		cmdface.Show(err.Error())
	}
}

func sendPlaygroundInfoFunc(params []string) {
	nin, err := strconv.Atoi(params[2])
	if err != nil {
		cmdface.Show(err.Error())
	}
	din, err := strconv.Atoi(params[1])
	if err != nil {
		cmdface.Show(err.Error())
	}
	cin, err := strconv.Atoi(params[0])
	if err != nil {
		cmdface.Show(err.Error())
	}
	dsin, err := strconv.Atoi(params[0])
	if err != nil {
		cmdface.Show(err.Error())
	}

	if err = sendPlaygroundInfo(nin, din, cin, dsin); err != nil {
		cmdface.Show(err.Error())
	}
}

func main() {

	cmdface.AddCommand(
		"sci",
		"<rid> <string>, join a room",
		sendConnectInfoFunc)
	cmdface.AddCommand(
		"sdi",
		"<rid>, left a room",
		sendDisconnectInfoFunc)
	cmdface.AddCommand(
		"spi",
		"<nin> <din> <cin> <dsin>, send playground information",
		sendPlaygroundInfoFunc)
	cmdface.AddCommand(
		"uid",
		"show uid fron server.",
		showUidFunc)
	cmdface.AddCommand(
		"pkgs",
		"show all received info packages",
		showInfoPkgListFunc)
	cmdface.AddCommand(
		"clean",
		"clean all packages",
		cleanInfoPkgListFunc)

	connToServer(2334, "/ws")
	for {
		if err := cmdface.InputAndRunCommand(">>> "); err != nil {
			cmdface.Show(fmt.Sprintf("%s\n", err.Error()))
		}
	}
}
