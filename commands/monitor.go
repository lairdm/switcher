package commands

import "switcher/monitor"

type Monitor struct {
	Input   monitor.Input `json:"input,omitempty"`
	Power   monitor.Power `json:"power,omitempty"`
	Monitor string        `json:"monitor,omitempty"`
}

func _(input monitor.Input, power monitor.Power, monitor string) *Monitor {
	return &Monitor{
		Input:   input,
		Power:   power,
		Monitor: monitor,
	}
}
