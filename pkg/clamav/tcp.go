package clamav

import (
	"fmt"
	"net"
)

func NewTCPClient(address string, port int) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		_ = fmt.Errorf("error creating tcp connection: %s", err)
	}
	return conn
}
