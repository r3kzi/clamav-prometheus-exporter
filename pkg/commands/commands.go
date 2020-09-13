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

package commands

import "fmt"

//Command corresponds to a ClamAV command that is accepted by `clamd` over the tcp socket. See `man clamd`.
type Command struct {
	Name   string
	Prefix string
}

var (
	//PING - Check the server's state. It should reply with "PONG".
	PING = Command{Name: "PING", Prefix: ""}

	//STATS - It is mandatory to newline terminate this command, or prefix with n or z.
	//Replies with statistics about the scan queue, contents of scan queue, and memory usage.
	STATS = Command{Name: "STATS", Prefix: "n"}

	VERSION = Command{Name: "VERSION", Prefix: ""}
)

func (c Command) String() string {
	if c.Prefix == "n" {
		return fmt.Sprintf("%s%s\n", c.Prefix, c.Name)
	}
	return c.Prefix + c.Name + "\n"
}
