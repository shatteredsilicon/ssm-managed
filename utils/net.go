package utils

import (
	"net"
)

// IsTCPAddressLocal checks if a tcp address local
func IsTCPAddressLocal(addr string) bool {
	ip, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return false
	}

	iAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}

	for _, iAddr := range iAddrs {
		if ipNet, ok := iAddr.(*net.IPNet); ok && ipNet.IP.Equal(ip.IP) {
			return true
		}
	}
	return false
}
