package commands

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Speaker struct {
	Command SpeakerCommandType `json:"command,omitempty"`
}

func _(command SpeakerCommandType) *Speaker {
	return &Speaker{
		Command: command,
	}
}

type SpeakerCommandType uint8

const (
	On SpeakerCommandType = iota + 1
	Off
	Up
	Down
	Mute
	Bluetooth
	Aux
	PC
	Opt
	Coax
)

var (
	Speaker_name = map[uint8]string{
		1:  "on",
		2:  "off",
		3:  "up",
		4:  "down",
		5:  "mute",
		6:  "bluetooth",
		7:  "aux",
		8:  "pc",
		9:  "opt",
		10: "coax",
	}
	Speaker_value = map[string]uint8{
		"on":        1,
		"off":       2,
		"up":        3,
		"down":      4,
		"mute":      5,
		"bluetooth": 6,
		"aux":       7,
		"pc":        8,
		"opt":       9,
		"coax":      10,
	}
)

func (s *SpeakerCommandType) UnmarshalJSON(data []byte) (err error) {
	var speaker string
	if err := json.Unmarshal(data, &speaker); err != nil {
		return err
	}
	if *s, err = ParseSpeaker(speaker); err != nil {
		return err
	}
	return nil
}

func (s SpeakerCommandType) String() string {
	return Speaker_name[uint8(s)]
}

func ParseSpeaker(s string) (SpeakerCommandType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Speaker_value[s]
	if !ok {
		return SpeakerCommandType(0), fmt.Errorf("%q is not a valid speaker command", s)
	}
	return SpeakerCommandType(value), nil
}
