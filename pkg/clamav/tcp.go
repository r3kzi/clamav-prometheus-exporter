package clamav

import (
	"fmt"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"io/ioutil"
	"net"
)

func ping(address string, port int) []byte {
	return dial(address, port, commands.PING)
}

func stats(address string, port int) []byte {
	return dial(address, port, commands.STATS)
}

func dial(address string, port int, command commands.Command) []byte {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		_ = fmt.Errorf("error creating tcp connection for command %s: %s", command, err)
	}
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("%s", command)))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		_ = fmt.Errorf("error reading tcp response for command %s: %s", command, err)
	}
	return resp
}
