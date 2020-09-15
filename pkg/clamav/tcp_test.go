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
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	listener net.Listener
	err      error
)

func TestMain(m *testing.M) {
	listener, err = net.Listen("tcp", "[::]:0")
	if err != nil {
		_ = fmt.Errorf("couldn't create tcp listener: %s", err)
	}
	defer listener.Close()
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	go func() {
		server, err := listener.Accept()
		defer server.Close()
		if err != nil {
			t.Errorf("failed to accept connect: %s", err)
		}
		resp, err := bufio.NewReader(server).ReadBytes('\n')
		if err != nil {
			t.Errorf("failed to read request: %s", err)
		}
		assert.Equal(t, "PING", strings.TrimSpace(string(resp)), "unexpected command")
		if _, err = server.Write([]byte("PONG\n")); err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}()
	client := New(listener.Addr().String())
	assert.Equal(t, []byte{'P', 'O', 'N', 'G', '\n'}, client.Dial(commands.PING))
}

func TestStats(t *testing.T) {
	go func() {
		server, err := listener.Accept()
		defer server.Close()
		if err != nil {
			t.Errorf("failed to accept connect: %s", err)
		}
		resp, err := bufio.NewReader(server).ReadBytes('\n')
		if err != nil {
			t.Errorf("failed to read request: %s", err)
		}
		assert.Equal(t, "nSTATS", strings.TrimSpace(string(resp)), "unexpected command")

		write := "POOLS: 1\n\n" +
			"STATE: VALID PRIMARY\n" +
			"THREADS: live 1  idle 0 max 12 idle-timeout 30\n" +
			"QUEUE: 0 items\n\t" +
			"STATS 0.000146 \n\n" +
			"MEMSTATS: heap 3.656M mmap 0.129M used 3.236M free 0.420M releasable 0.127M pools 1 pools_used 1089.550M pools_total 1089.585M\n" +
			"END"
		if _, err = server.Write([]byte(write)); err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}()
	client := New(listener.Addr().String())
	stats := client.Dial(commands.STATS)

	regex := regexp.MustCompile("([0-9.]+)")
	matches := regex.FindAllStringSubmatch(string(stats), -1)

	assert.Equal(t, "1", matches[1][1])
	assert.Equal(t, "0", matches[2][1])
	assert.Equal(t, "12", matches[3][1])
	assert.Equal(t, "0", matches[5][1])
	assert.Equal(t, "3.656", matches[7][1])
	assert.Equal(t, "0.129", matches[8][1])
	assert.Equal(t, "3.236", matches[9][1])

}

func TestVersion(t *testing.T) {
	go func() {
		server, err := listener.Accept()
		defer server.Close()
		if err != nil {
			t.Errorf("failed to accept connect: %s", err)
		}
		resp, err := bufio.NewReader(server).ReadBytes('\n')
		if err != nil {
			t.Errorf("failed to read request: %s", err)
		}
		assert.Equal(t, "VERSION", strings.TrimSpace(string(resp)), "unexpected command")

		write := "ClamAV 0.102.4/25913/Fri Aug 28 13:19:15 2020\n"
		if _, err = server.Write([]byte(write)); err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}()
	client := New(listener.Addr().String())
	stats := client.Dial(commands.VERSION)

	regex := regexp.MustCompile("((ClamAV)+\\s([0-9.]*)/([0-9.]*))")
	matches := regex.FindAllStringSubmatch(string(stats), -1)

	assert.Equal(t, "0.102.4", matches[0][3])
	assert.Equal(t, "25913", matches[0][4])
}
