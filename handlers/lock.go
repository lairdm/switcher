package handlers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"switcher/commands"
	"switcher/config"
)

func LockHandler(rawCommand json.RawMessage, settings *config.Settings) {
	var command commands.Lock
	if err := json.Unmarshal(rawCommand, &command); err != nil {
		fmt.Printf("Error decoding lock command: %s\n", err)
		return
	}

	if command.Unlock {
		fmt.Println("Unlocking")
		cmd := exec.Command(settings.XLock.Bin)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("The screen might already be unlocked")
		}
	}
}
