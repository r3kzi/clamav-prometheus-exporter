package commands

import "fmt"

type Command struct {
	Name   string
	Prefix string
}

var (
	VERSION = Command{Name: "VERSION", Prefix: ""}
	PING    = Command{Name: "PING", Prefix: ""}
	STATS   = Command{Name: "STATS", Prefix: "n"}
)

func (c Command) String() string {
	if c.Prefix == "n" {
		return fmt.Sprintf("%s%s\n", c.Prefix, c.Name)
	}
	return c.Prefix + c.Name + "\n"
}
