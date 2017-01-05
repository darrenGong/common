package common

import (
	"net"
	"log"
	"errors"
)

func GetIPV4Addr(iface string) (string, error) {
	inface, err := net.InterfaceByName(iface)
	if err != nil {
		log.Printf("Failed to get interface by name[%s]", iface)
		return "", err
	}

	addrs, err := inface.Addrs()
	if err != nil {
		log.Println("Failed to get addrs")
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.To4().String(), nil
			}
		}
	}

	return "", errors.New("no valid addr")
}
