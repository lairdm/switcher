package handlers

import (
	"encoding/json"
	"fmt"
	"switcher/commands"
	"switcher/config"
)

func SpeakerHandler(rawCommand json.RawMessage, settings *config.Settings) {
	var command commands.Speaker
	if err := json.Unmarshal(rawCommand, &command); err != nil {
		errStr := fmt.Sprintf("Error decoding message: %s\n", err)
		panic(errStr)
	}

	fmt.Printf("Speaker: %s\n", command)

	switch command.Command {
	case commands.Mute:
		fmt.Println("Muting sound")
	case commands.Up:
		fmt.Println("Increasing volume")
	case commands.Down:
		fmt.Println("Decreasing volume")
	case commands.On:
		fmt.Println("Turning on")
	case commands.Off:
		fmt.Println("Turning off")
	case commands.Bluetooth:
		fmt.Println("Toggling bluetooth")
	case commands.Aux:
		fmt.Println("Toggling aux")
	case commands.PC:
		fmt.Println("Toggling PC")
	case commands.Opt:
		fmt.Println("Toggling optical")
	case commands.Coax:
		fmt.Println("Toggling coax")
	}

}
