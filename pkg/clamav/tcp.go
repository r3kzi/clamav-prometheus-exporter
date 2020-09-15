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
	"fmt"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
)

//Client corresponds to a ClamAV client
type Client struct {
	address string
}

//New create a new TCP Client for ClamAV
func New(address string) *Client {
	return &Client{
		address: address,
	}
}

// Dial connects to a tcp socket based on address. Sends commands.Command.
func (c Client) Dial(command commands.Command) []byte {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		log.Errorf("error creating tcp connection for command %s: %s", command, err)
		return nil
	}
	defer conn.Close()

	_, _ = conn.Write([]byte(fmt.Sprintf("%s", command)))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Errorf("error reading tcp response for command %s: %s", command, err)
		return nil
	}
	return resp
}
