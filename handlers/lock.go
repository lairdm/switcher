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
		errStr := fmt.Sprintf("Error decoding message: %s\n", err)
		panic(errStr)
	}

	fmt.Printf("Lock: %s\n", command)

	if !command.Lock {
		fmt.Println("Unlocking")
		cmd := exec.Command(settings.XLock.Bin)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(string(out))
			panic(err)
		}
	}
}
