package common

import (
	"sync"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var (
	gZKConn  *zk.Conn
	gConnMap *ConnMap

	gSyncWaitGroup sync.WaitGroup
)

const (
	CONNECT_TIMEOUT = 5 * time.Second
)


func InitCommon(servers, selfPath, selfValue string) {
	InitZookeeper(servers, selfPath, selfValue)
	InitConnMap()
}

func Release() {
	gSyncWaitGroup.Wait()
}