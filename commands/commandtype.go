package commands

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Command struct {
	Command CommandType     `json:"command,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func _(command CommandType, data json.RawMessage) *Command {
	return &Command{
		Command: command,
		Data:    data,
	}
}

type CommandType uint8

const (
	MonitorCommand CommandType = iota + 1
	SoundCommand
	LockCommand
	SpeakerCommand
)

var (
	Command_name = map[uint8]string{
		1: "monitor",
		2: "sound",
		3: "lock",
		4: "speaker",
	}
	Command_value = map[string]uint8{
		"monitor": 1,
		"sound":   2,
		"lock":    3,
		"speaker": 4,
	}
)

func (c *CommandType) UnmarshalJSON(data []byte) (err error) {
	var command string
	if err := json.Unmarshal(data, &command); err != nil {
		return err
	}
	if *c, err = ParseCommand(command); err != nil {
		return err
	}
	return nil
}

func (c CommandType) String() string {
	return Command_name[uint8(c)]
}

func ParseCommand(s string) (CommandType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Command_value[s]
	if !ok {
		return CommandType(0), fmt.Errorf("%q is not a valid command", s)
	}
	return CommandType(value), nil
}
