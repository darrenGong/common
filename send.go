package common

import (
	"log"
	"net"
)

const (
	UVERMANAGER  = "uvermanager"
	HELLOMANAGER = "hellomanager"
)

func InitConnMap() {
	gConnMap = NewConnMap()
}

func SendMsgToUVERManager(msg []byte) error {
	ifConn, err := getConn(UVERMANAGER, gConfig.NormalServerMap[UVERMANAGER])
	if err != nil {
		log.Printf("Failed to get conn[srv:%s, path:%s]", UVERMANAGER, gConfig.NormalServerMap[UVERMANAGER])
		return err
	}

	conn := ifConn.(net.Conn)
	log.Printf("Send data to uvermanager[%s]", conn.RemoteAddr().String())

	return nil
}
