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
)

func (c Command) String() string {
	if c.Prefix == "n" {
		return fmt.Sprintf("%s%s\n", c.Prefix, c.Name)
	}
	return c.Prefix + c.Name + "\n"
}
