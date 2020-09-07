package clamav

import (
	"fmt"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"io/ioutil"
	"net"
	"strings"
)

func ping(address string) float64 {
	bytes := dial(address, commands.PING)
	if strings.TrimSpace(string(bytes)) == "PONG" {
		return 1
	}
	return 0
}

func stats(address string) []byte {
	return dial(address, commands.STATS)
}

func dial(address string, command commands.Command) []byte {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		_ = fmt.Errorf("error creating tcp connection for command %s: %s", command, err)
		return nil
	}
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("%s", command)))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		_ = fmt.Errorf("error reading tcp response for command %s: %s", command, err)
		return nil
	}
	return resp
}
