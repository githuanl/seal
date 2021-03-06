package co

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgUserCtrl(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("MsgUserCtrl")

	if nil == msg {
		return
	}

	p := pt.UserControlPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	log.Println("MsgUserCtrl event type=", p.EventType)

	switch p.EventType {
	case pt.SrcPCUCStreamBegin:
	case pt.SrcPCUCStreamEOF:
	case pt.SrcPCUCStreamDry:
	case pt.SrcPCUCSetBufferLength:
	case pt.SrcPCUCStreamIsRecorded:
	case pt.SrcPCUCPingRequest:
		err = rc.ctrlPingRequest(&p)
	case pt.SrcPCUCPingResponse:
	default:
		log.Println("msg user ctrl unknown event type.type=", p.EventType)
	}

	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) ctrlPingRequest(p *pt.UserControlPacket) (err error) {

	log.Println("ctrl ping request.")

	if pt.SrcPCUCSetBufferLength == p.EventType {

	} else if pt.SrcPCUCPingRequest == p.EventType {
		var pp pt.UserControlPacket
		pp.EventType = pt.SrcPCUCPingResponse
		pp.EventData = p.EventData
		if err = rc.sendPacket(&pp, 0); err != nil {
			return
		}

		log.Println("send ping response success.")

	}

	return
}
