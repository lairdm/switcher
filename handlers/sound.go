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
		errStr := fmt.Sprintf("Error decoding message: %s\n", err)
		panic(errStr)
	}

	fmt.Printf("Sound: %s\n", command)

	if command.Volume != 0 {
		VolumeHandler(command, settings)
	}
}

func VolumeHandler(command commands.Sound, settings *config.Settings) {
	fmt.Printf("Attempting to change volume to %d\n", command.Volume)

	var cmd *exec.Cmd
	switch command.Volume {
	case sound.Mute:
		fmt.Println("Muting sound")
		cmd = exec.Command(settings.Amixer.Bin, "sset", "Master", "0")
	case sound.Up:
		fmt.Println("Increasing volume")
		cmd = exec.Command(settings.Amixer.Bin, "sset", "Master", "5%+")
	case sound.Down:
		fmt.Println("Decreasing volume")
		cmd = exec.Command(settings.Amixer.Bin, "sset", "Master", "5%-")
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}
