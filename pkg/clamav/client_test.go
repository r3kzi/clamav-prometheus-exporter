/*
Copyright 2020 Christian Niehoff.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clamav

import (
	"bufio"
	"fmt"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"github.com/stretchr/testify/assert"
	"net"
	"regexp"
	"testing"
)

var tests = []struct {
	network, address string
}{
	{"tcp", "[::]:0"},
	{"unix", "/tmp/tmp.sock"},
}

func TestClient(t *testing.T) {
	for _, test := range tests {
		listener, _ := net.Listen(test.network, test.address)
		defer listener.Close()
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					continue
				}
				req, err := bufio.NewReader(conn).ReadBytes('\n')
				if err != nil {
					_ = fmt.Errorf("failed to read request: %s", err)
				}
				switch string(req) {
				case "PING\n":
					if _, err = conn.Write([]byte("PONG\n")); err != nil {
						t.Errorf("failed to write response: %s", err)
					}
				case "nSTATS\n":
					write := "POOLS: 1\n\n" +
						"STATE: VALID PRIMARY\n" +
						"THREADS: live 1  idle 0 max 12 idle-timeout 30\n" +
						"QUEUE: 0 items\n\t" +
						"        FILDES 41.249971 fd[11]\n" +
						"        STATS 0.000075\n" +
						"STATS 0.000146 \n\n" +
						"MEMSTATS: heap 3.656M mmap 0.129M used 3.236M free 0.420M releasable 0.127M pools 1 pools_used 1089.550M pools_total 1089.585M\n" +
						"END"
					if _, err = conn.Write([]byte(write)); err != nil {
						t.Errorf("failed to write response: %s", err)
					}
				}
				conn.Close()
			}
		}()
		client := New(listener.Addr().String(), test.network)
		assert.Equal(t, []byte{'P', 'O', 'N', 'G', '\n'}, client.Dial(commands.PING))

		stats := client.Dial(commands.STATS)
		regex := regexp.MustCompile("([0-9.]+)")
		matches := regex.FindAllStringSubmatch(string(stats), -1)
		assert.Equal(t, "1", matches[1][1])
		assert.Equal(t, "0", matches[2][1])
		assert.Equal(t, "12", matches[3][1])
		assert.Equal(t, "0", matches[5][1])
		assert.Equal(t, "3.656", matches[10][1])
		assert.Equal(t, "0.129", matches[11][1])
		assert.Equal(t, "3.236", matches[12][1])
	}
}
