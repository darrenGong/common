package common

import (
	"regexp"
	"errors"
	"strings"
	"log"
	"net"
)


const (
	UVERMANAGER = "uvermanager"
)

func InitConnMap() {
	gConnMap = NewConnMap()
}

func SendMsgToUVERManager(msg []byte) error {
	ifConn, err := getConn(UVERMANAGER, gConfig.UVERManager)
	if err != nil {
		log.Printf("Failed to get conn[srv:%s, path:%s]", UVERMANAGER, gConfig.UVERManager)
		return err
	}

	conn := ifConn.(net.Conn)
	log.Printf("Send data to uvermanager[%s]", conn.RemoteAddr().String())

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
		log.Printf("Failed to regex check[src:%s]", addrs[0])
		return false, err
	}

	// port
	portRegexExpression := "\b[0-9]+\b"
	if err := RegexpCheck(addrs[1], portRegexExpression); err != nil {
		log.Printf("Failed to regex check[src:%s]", addrs[0])
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

		conn, err := net.DialTimeout("tcp", laddr, CONNECT_TIMEOUT)
		if err != nil {
			log.Printf("Failed to dial connect[%s]", laddr)
			return nil, err
		}

		gConnMap.AddConn(serverType, conn)
	}

	return conn
}
