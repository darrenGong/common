package common

import (
	"strings"
	"uframework/message/proto/uns"
	"log"
	"errors"
	"math/rand"

	"github.com/samuel/go-zookeeper/zk"
	"time"
	"uframework/message/protobuf/proto"
	"fmt"
)

const (
	TIMEOUT = 5 * time.Second
	REGISTER_TIMEOUT = 20 * time.Second
)

func InitZookeeper(servers, selfPath, selfValue string) {
	zkConn, _, err := zk.Connect(strings.Split(servers, ","), TIMEOUT)
	if err != nil {
		log.Fatalf("Failed to connect zookeeper[%v]", err)
	}
	gZKConn = zkConn

	go CreateMySelf(selfPath, selfValue)
}

func CreateMySelf(path, value string) {
	gSyncWaitGroup.Add(1)
	defer gSyncWaitGroup.Done()

	for {
		if _, err := createNode(path, []byte(value)); err != nil {
			log.Printf("Failed to create myself[path:%s, value:%s, err:%v]", path, value, err)
		}

		time.Sleep(REGISTER_TIMEOUT)
	}
}

func GetServerNode(path string) (string, error) {
	valueByte, err := getNode(path)
	if err != nil {
		log.Printf("Failed to get node[%s]", path)
		return "", err
	}

	ipaddr, port, err := parseZKProtoContent(valueByte)
	if err != nil {
		log.Printf("Failed to parse proto content[%v]", valueByte)
		return "", err
	}

	log.Printf("ipaddr: %s, port:%d", ipaddr, port)
	return strings.Join([]string{ipaddr, fmt.Sprintf("%d", port)}, ":"), nil
}

func getNode(path string) ([]byte, error) {
	nodePath, err := getNodePath(path)
	if err != nil {
		log.Printf("Failed to get node path[%s]", path)
		return nil, err
	}

	valueByte, _, err := gZKConn.Get(nodePath)
	if err != nil {
		log.Printf("Failed to get node value[path:%s]", nodePath)
		return nil, err
	}

	return valueByte, nil
}

func createNode(path string, data []byte) (string, error) {
	paths := strings.Split(path, "/")

	var nodePath, reallyPath string
	for index, node := range paths[1:] {
		nodePath = strings.Join([]string{nodePath, node}, "/")
		bValid, _, err := gZKConn.Exists(nodePath)
		if err != nil {
			log.Fatalf("Failed to check valid node path[%s]", nodePath)
			return "", err
		}

		if bValid {
			continue
		}

		flags := int32(0)
		value := []byte(nil)
		if index == len(paths) - 2 {
			flags = zk.FlagEphemeral
			value = data
		}

		reallyPath, err = gZKConn.Create(nodePath, value, flags, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalf("Failed to create node path[%s]", nodePath)
			return "", err
		}

	}

	return reallyPath, nil
}

func getNodePath(path string) (string, error) {
	children, stat, err := gZKConn.Children(path)
	if err != nil {
		log.Printf("Failed to get children[%v]", err)
		return "", err
	}

	if stat.NumChildren < 1 {
		log.Printf("There are not have child node[%+v]", stat)
		return "", errors.New("Children num is zero")
	}

	lenChildren := len(children)
	if 0 == lenChildren {
		log.Println("children length is zero")
		return "", errors.New("children length is zero")
	}

	return strings.Join([]string{path, children[rand.Intn(lenChildren)]}, "/"), nil
}

func parseZKProtoContent(values []byte) (string, uint32, error) {
	nodeContent := &uns.NameNodeContent{}
	if err := proto.Unmarshal(values, nodeContent); err != nil {
		log.Printf("Failed to unmarshal bytes[%v]", values)
		return "", 0, err
	}

	return nodeContent.GetIp(), nodeContent.GetPort(), nil
}