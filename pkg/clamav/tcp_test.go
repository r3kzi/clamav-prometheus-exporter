package clamav

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"net"
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
		if _, err = server.Write([]byte("PONG")); err != nil {
			t.Errorf("failed to write response: %s", err)
		}
	}()
	assert.Equal(t, float64(1), ping(listener.Addr().String()), "")
}
