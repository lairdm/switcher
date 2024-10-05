package monitor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Input uint8

const (
	HDMI Input = iota + 1
	DP
	USB_C
	DVI
)

var (
	Input_name = map[uint8]string{
		1: "hdmi",
		2: "dp",
		3: "usbc",
		4: "dvi",
	}
	Input_value = map[string]uint8{
		"hdmi": 1,
		"dp":   2,
		"usbc": 3,
		"dvi":  4,
	}
)

func (i *Input) UnmarshalJSON(data []byte) (err error) {
	var inputs string
	if err := json.Unmarshal(data, &inputs); err != nil {
		return err
	}
	if *i, err = ParseInput(inputs); err != nil {
		return err
	}
	return nil
}

func (i Input) String() string {
	return Input_name[uint8(i)]
}

func ParseInput(s string) (Input, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Input_value[s]
	if !ok {
		return Input(0), fmt.Errorf("%q is not a valid input", s)
	}
	return Input(value), nil
}
