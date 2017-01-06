package common

import (
	"log"
	"net"
	"google.golang.org/grpc"
)

const (
	UVERMANAGER  = "uvermanager"
	HELLOMANAGER = "hellomanager"
)

func InitConnMap() {
	gConnMap = NewConnMap()
}

func ReleaseConn() {
	for key, ifConn := range gConnMap.m {
		if _, ok := gConfig.NormalServerMap[key]; ok {
			conn := ifConn.(net.Conn)
			conn.Close()
		}

		if _, ok := gConfig.GrpcServerMap[key]; ok {
			conn := ifConn.(grpc.ClientConn)
			conn.Close()
		}

		log.Printf("Invalid server type when release conn[%s]", key)
	}
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

func GetConnToHelloManager() (*grpc.ClientConn, error) {
	ifConn, err := getConn(HELLOMANAGER, gConfig.NormalServerMap[HELLOMANAGER])
	if err != nil {
		log.Printf("Failed to get conn[srv:%s, path:%s]", HELLOMANAGER, gConfig.NormalServerMap[HELLOMANAGER])
		return nil, err
	}

	return ifConn.(*grpc.ClientConn), nil
}
