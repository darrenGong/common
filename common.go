package common

import (
	"sync"
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"go_test/common/config"
	"io"
	"io/ioutil"
	"log"
	"encoding/json"
)

var (
	gZKConn  *zk.Conn
	gConnMap *ConnMap
	gConfig  *config.Config

	gSyncWaitGroup sync.WaitGroup
)

const (
	CONNECT_TIMEOUT = 5 * time.Second
)


func InitCommon(selfPath, selfValue string) {
	if err := ParseConfig("\\config/config.json", &gConfig);
		err != nil {
		log.Fatalf("Failed to parse config[err:%v]", err)
	}

	InitZookeeper(gConfig.Servers, selfPath, selfValue)
	InitConnMap()
}

func Release() {
	gSyncWaitGroup.Wait()
}

func ParseConfig(configPath string, config *config.Config) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read file[%s]", configPath)
		return err
	}

	if err := json.Unmarshal(bytes, config); err != nil {
		log.Fatalf("Failed to unmarshal json file[%s]", configPath)
		return err
	}

	return nil
}