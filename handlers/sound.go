package handlers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"switcher/commands"
	"switcher/config"
	"switcher/sound"
)

func SoundHandler(rawCommand json.RawMessage, settings *config.Settings) {
	var command commands.Sound
	if err := json.Unmarshal(rawCommand, &command); err != nil {
		fmt.Printf("Error decoding sound command: %s\n", err)
		return
	}

	fmt.Printf("Sound: %s\n", command)

	if command.Volume != 0 {
		VolumeHandler(command, settings)
	}
}

func VolumeHandler(command commands.Sound, settings *config.Settings) {
	fmt.Printf("Attempting to change volume to %d\n", command.Volume)

	var args []string
	switch command.Volume {
	case sound.Mute:
		fmt.Println("Muting sound")
		args = []string{"sset", "Master", "0"}
	case sound.Up:
		fmt.Println("Increasing volume")
		args = []string{"sset", "Master", "5%+"}
	case sound.Down:
		fmt.Println("Decreasing volume")
		args = []string{"sset", "Master", "5%-"}
	default:
		fmt.Printf("Unknown volume command: %s\n", command.Volume.String())
		return
	}

	out, err := exec.Command(settings.Amixer.Bin, args...).Output()
	if err != nil {
		fmt.Printf("amixer failed: %v\n", err)
		if len(out) > 0 {
			fmt.Println(string(out))
		}
	}
}
