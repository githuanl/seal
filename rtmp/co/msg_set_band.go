package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) MsgSetBand(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("MsgSetBand")

	p := pt.SetPeerBandWidthPacket{}
	err = p.Decode(msg.Payload)

	return
}
