package clamav

import (
	"bytes"
	"fmt"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
)

var (
	//THREADS: live 1  idle 0 max 12
	threadsRegex = regexp.MustCompile("(live)+\\s*([0-9]+)\\s*(idle)\\s*([0-9]+)\\s*(max)\\s*([0-9]+)")
)

type Stats struct {
	Threads Threads
}

type Threads struct {
	Live float64
	Idle float64
	Max  float64
}

func ping(address string) float64 {
	if bytes.Compare(dial(address, commands.PING), []byte{'P', 'O', 'N', 'G', '\n'}) == 0 {
		return 1
	}
	return 0
}

func stats(address string) Stats {
	resp := dial(address, commands.STATS)

	toFloat := func(s string) float64 {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			_ = fmt.Errorf("couldn't parse string to float: %s", err)
		}
		return float
	}

	threadsLine := threadsRegex.FindAllStringSubmatch(string(resp), 1)
	threads := Threads{
		Live: toFloat(threadsLine[0][2]),
		Idle: toFloat(threadsLine[0][4]),
		Max:  toFloat(threadsLine[0][6]),
	}

	stats := Stats{
		Threads: threads,
	}
	return stats
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
