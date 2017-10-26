package rrserver

import (
	"fmt"
	"net"
)

const (
	IP_PROTOCOL       = "ipv4"
	DEFAULT_POOL_SIZE = 10
)

var (
	cp *TCPConnectionPool
)

func init() {
	cp = CreateTCPConnectionPool(DEFAULT_POOL_SIZE)
}

func SendTCPRequest(addr string, msg []byte) (error, []byte) {
	err, c := cp.Get(addr)
	if err != nil {
		return err, nil
	}
	if err := c.SetKeepAlive(true); err != nil {
		return err, nil
	}
	if err := c.Write(msg); err != nil {
		return err, nil
	}
	err, b := c.Read()
	cp.Add(addr, c)
	return err, b
}

func getIpAddrByInterface(inf string) (error, string) {
	if len(inf) < 1 {
		return fmt.Errorf("Interface name is an empty string!"), inf
	}
	if net.ParseIP(inf) != nil {
		return nil, inf
	}
	infHandle, err := net.InterfaceByName(inf)
	if err != nil {
		return err, inf
	}
	infAddrs, err := infHandle.Addrs()
	if err != nil {
		return err, inf
	}
	for _, addr := range infAddrs {
		ipHandle, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return err, inf
		}
		if IP_PROTOCOL == "ipv4" {
			if ipHandle.To4() != nil {
				return nil, ipHandle.To4().String()
			}
		} else if IP_PROTOCOL == "ipv6" {
			if ipHandle.To16() != nil {
				return nil, ipHandle.To16().String()
			}
		} else {
			return fmt.Errorf("Ip protocol [%s] not support", IP_PROTOCOL), inf
		}
	}
	return fmt.Errorf("Failed when try to get ip address, [%s]", inf), inf
}
