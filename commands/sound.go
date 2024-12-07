package commands

import "switcher/sound"

type Sound struct {
	Volume sound.Volume `json:"volume,omitempty"`
}

func _(volume sound.Volume) *Sound {
	return &Sound{
		Volume: volume,
	}
}
