package handlers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"switcher/commands"
	"switcher/config"
	"switcher/monitor"
)

func MonitorHandler(rawCommand json.RawMessage, settings *config.Settings) {
	var command commands.Monitor
	if err := json.Unmarshal(rawCommand, &command); err != nil {
		errStr := fmt.Sprintf("Error decoding message: %s\n", err)
		panic(errStr)
	}

	fmt.Printf("Monitor: %s\n", command)
	var monitorProfile config.Monitor
	ok := false
	if monitorProfile, ok = settings.Monitors[command.Monitor]; !ok {
		errorString := fmt.Sprintf("Monitor %s not found\n", command.Monitor)
		panic(errorString)
	}

	if command.Input != 0 {
		InputHandler(command, monitorProfile, settings)
	}

	if command.Power != 0 {
		PowerHandler(command, monitorProfile, settings)
	}
}

func InputHandler(command commands.Monitor, monitorProfile config.Monitor, settings *config.Settings) {

	if _, ok := monitorProfile.Inputs[command.Input.String()]; !ok {
		errorString := fmt.Sprintf("Input %s not found for monitor %s\n", command.Input.String(), command.Monitor)
		panic(errorString)
	}

	fmt.Printf("Attempting to change monitor %s to input %s\n", monitorProfile.Serial, monitorProfile.Inputs[command.Input.String()])
	fmt.Println(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", "60", monitorProfile.Inputs[command.Input.String()])

	cmd := exec.Command(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", "60", monitorProfile.Inputs[command.Input.String()])
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func PowerHandler(command commands.Monitor, monitorProfile config.Monitor, settings *config.Settings) {
	if monitorProfile.Display == "" {
		errorString := fmt.Sprintf("Display %s not found for monitor %s\n", command.Power.String(), command.Monitor)
		panic(errorString)
	}

	fmt.Printf("Attempting to change monitor %s to power %s\n", monitorProfile.Serial, command.Power.String())
	//	fmt.Println(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", "60", monitorProfile.Inputs[command.Input.String()])

	var cmd *exec.Cmd
	switch command.Power {
	case monitor.On:
		cmd = exec.Command(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", "0xD6", "0x01")
	case monitor.Off:
		cmd = exec.Command(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", "0xD6", "0x04")
	case monitor.Wake:
		cmd = exec.Command(settings.Xset.Bin, "-display", monitorProfile.Display, "dpms", "force", "on")
	default:
		errStr := fmt.Sprintf("Invalid power state: %s\n", command.Power.String())
		panic(errStr)
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}
