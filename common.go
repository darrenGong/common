package common

import (
	"encoding/json"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"go_test/common/config"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
	"google.golang.org/grpc"
)

var (
	gZKConn  *zk.Conn
	gConnMap *ConnMap

	gConfig        config.Config
	gSyncWaitGroup sync.WaitGroup
)

const (
	CONNECT_TIMEOUT = 5 * time.Second
)

func InitCommon(selfPath, selfValue string) {
	if err := ParseConfig("F:/go-dev/src/go_test/common/config/config.json", &gConfig); err != nil {
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
		log.Printf("Failed to read file[%s]", configPath)
		return err
	}

	if err := json.Unmarshal(bytes, config); err != nil {
		log.Printf("Failed to unmarshal json file[%s]", configPath)
		return err
	}

	return nil
}

func CheckZKValue(value string) (bool, error) {
	addrs := strings.Split(value, ":")
	if len(addrs) != 2 {
		log.Printf("Addrs length is not equal at two[len:%d]", len(addrs))
		return false, errors.New("Invalid length")
	}

	// IP
	ipRegexExpression := "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	if err := RegexpCheck(addrs[0], ipRegexExpression); err != nil {
		log.Printf("Failed to regex check[%s]", addrs[0])
		return false, err
	}

	// port
	portRegexExpression := "[0-9]+"
	if err := RegexpCheck(addrs[1], portRegexExpression); err != nil {
		log.Printf("Failed to regex check[%s]", addrs[1])
		return false, err
	}

	return true, nil
}

func RegexpCheck(src, regexExpression string) error {
	re, _ := regexp.Compile(regexExpression)
	if !re.MatchString(src) {
		return errors.New("regular expression mismatch")
	}

	return nil
}

func getConn(serverType string, path string) (interface{}, error) {
	conn := gConnMap.GetConn(serverType)
	if conn == nil {
		log.Printf("Invalid connect so that reconnect target[%s]", serverType)

		laddr, err := GetServerNode(path)
		if err != nil {
			log.Printf("Failed to get node value[%s]", path)
			return nil, err
		}

		if bValid, err := CheckZKValue(laddr); !bValid || err != nil {
			log.Printf("Failed to check value[%s]", laddr)
			return nil, err
		}

		conn, err = getReallyConn(serverType, laddr)
		if err != nil {
			log.Printf("Failed to dial connect[%s]", laddr)
			return nil, err
		}

		gConnMap.AddConn(serverType, conn)
	}

	return conn, nil
}

func getReallyConn(serverType, laddr string) (interface{}, error) {
	if _, ok := gConfig.NormalServerMap[serverType]; ok {
		return net.DialTimeout("tcp", laddr, CONNECT_TIMEOUT)
	}

	if _, ok := gConfig.GrpcServerMap[serverType]; ok {
		return grpc.Dial(laddr, grpc.WithInsecure())
	}

	return nil, errors.New("Invalid server type")
}
