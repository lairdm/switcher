package monitor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Power uint8

const (
	On Power = iota + 1
	Off
	Wake
	Sleep
)

var (
	Power_name = map[uint8]string{
		1: "on",
		2: "off",
		3: "wake",
		4: "sleep",
	}
	Power_value = map[string]uint8{
		"on":   1,
		"off":  2,
		"wake": 3,
		"sleep": 4,
	}
)

func (p *Power) UnmarshalJSON(data []byte) (err error) {
	var power string

	if err := json.Unmarshal(data, &power); err != nil {
		return err
	}

	if *p, err = ParsePower(power); err != nil {
		return err
	}

	return nil
}

func (p Power) String() string {
	return Power_name[uint8(p)]
}

func ParsePower(s string) (Power, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Power_value[s]
	if !ok {
		return Power(0), fmt.Errorf("%q is not a valid power state", s)
	}
	return Power(value), nil
}
