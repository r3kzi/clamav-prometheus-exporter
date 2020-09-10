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
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"github.com/stretchr/testify/assert"
	"net"
	"regexp"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	listener, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Errorf("couldn't create tcp listener: %s", err)
	}
	defer listener.Close()

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
	assert.Equal(t, []byte{'P', 'O', 'N', 'G', '\n'}, dial(listener.Addr().String(), commands.PING))
}

func TestStats(t *testing.T) {
	listener, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Errorf("couldn't create tcp listener: %s", err)
	}
	defer listener.Close()

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

		write := "THREADS: live 1  idle 0 max 12\n" +
			"MEMSTATS: heap 3.656M mmap 0.129M used 3.236M free 0.420M releasable 0.127M pools 1 pools_used 1089.550M pools_total 1089.585M\n"
		if _, err = server.Write([]byte(write)); err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}()
	stats := dial(listener.Addr().String(), commands.STATS)
	regex = regexp.MustCompile("(live|idle|max|heap|mmap|\\bused)\\s([0-9.]+)[MG]*")
	matches := regex.FindAllStringSubmatch(string(stats), -1)
	assert.Equal(t, "1", matches[0][2])
	assert.Equal(t, "0", matches[1][2])
	assert.Equal(t, "12", matches[2][2])
	assert.Equal(t, "3.656", matches[3][2])
	assert.Equal(t, "0.129", matches[4][2])
	assert.Equal(t, "3.236", matches[5][2])
}
