package sound

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Volume uint8

const (
	Mute Volume = iota + 1
	Up
	Down
)

var (
	Volume_name = map[uint8]string{
		1: "mute",
		2: "up",
		3: "down",
	}
	Volume_value = map[string]uint8{
		"mute": 1,
		"up":   2,
		"down": 3,
	}
)

func (v *Volume) UnmarshalJSON(data []byte) (err error) {
	var volume string

	if err := json.Unmarshal(data, &volume); err != nil {
		return err
	}

	if *v, err = ParseVolume(volume); err != nil {
		return err
	}

	return nil
}

func (v Volume) String() string {
	return Volume_name[uint8(v)]
}

func ParseVolume(s string) (Volume, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Volume_value[s]
	if !ok {
		return Volume(0), fmt.Errorf("%q is not a valid volume state", s)
	}
	return Volume(value), nil
}
