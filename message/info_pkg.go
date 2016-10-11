package message

import (
	"barrage-server/base"
)

type infoType msgType

// Info is a interfase used as InfoPkg body.
type Info interface {
	base.CommunicationData
}

// InfoPkg is used to transfer data among major module(user, room and playground).
type InfoPkg interface {
	Type() infoType
	Body() Info
}
