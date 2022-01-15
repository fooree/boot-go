package container

import (
	"log"
	"strings"
)

type configuration struct {
	values      map[string]string
	commandArgs map[string]string
}

func (c *configuration) Get(key string) (value string, ok bool) {
	value, ok = c.values[key]
	return value, ok
}

func (c *configuration) ReadCommandArgs(args []string) {
	log.Println("origin command args:", args)
	for {
		if len(args) == 0 {
			break
		}

		arg := args[0]
		args = args[1:]

		if len(arg) < 3 { // too short, skip
			if len(arg) == 2 && arg[0] == '-' && arg[1] == '-' {
				break // "--" terminates the flags
			}
			continue
		}

		if arg[1] != 'D' || arg[2] == '=' { // -?*   -D=
			continue
		}

		name := arg[2:]
		idx := strings.IndexByte(name, '=')
		if idx > 0 { // -Dkey=value
			value := name[idx+1:]
			key := name[:idx]
			c.commandArgs[key] = value
		} else {
			c.commandArgs[name] = "" // -Dkey
		}
	}
	log.Println("command args:", c.commandArgs)
}

func newConfiguration() *configuration {
	return &configuration{
		values:      make(map[string]string, 16),
		commandArgs: make(map[string]string, 16),
	}
}
