package handlers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"switcher/commands"
	"switcher/config"
	"switcher/monitor"
	"time"
)

func MonitorHandler(rawCommand json.RawMessage, settings *config.Settings) {
	var command commands.Monitor
	if err := json.Unmarshal(rawCommand, &command); err != nil {
		fmt.Printf("Error decoding monitor command: %s\n", err)
		return
	}

	fmt.Printf("Monitor: %s\n", command)
	var monitorProfile config.Monitor
	ok := false
	if monitorProfile, ok = settings.Monitors[command.Monitor]; !ok {
		fmt.Printf("Monitor %s not found\n", command.Monitor)
		return
	}

	if command.Input != 0 {
		InputHandler(command, monitorProfile, settings)
	}

	if command.Power != 0 {
		PowerHandler(command, monitorProfile, settings)
	}
}

func ddcutilArgs(monitorProfile config.Monitor) []string {
	if monitorProfile.Bus != 0 {
		return []string{"--bus", strconv.Itoa(monitorProfile.Bus)}
	}
	return []string{"--sn", monitorProfile.Serial}
}

// runDdcutil runs ddcutil, retrying once. DDC/i2c is occasionally flaky and
// ddcutil doesn't always succeed on the first attempt.
func runDdcutil(bin string, args []string) {
	var out []byte
	var err error
	for attempt := 1; attempt <= 2; attempt++ {
		out, err = exec.Command(bin, args...).Output()
		if err == nil {
			return
		}
		fmt.Printf("ddcutil attempt %d failed: %v\n", attempt, err)
		if len(out) > 0 {
			fmt.Println(string(out))
		}
		if attempt < 2 {
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Printf("ddcutil gave up: %v\n", err)
}

func InputHandler(command commands.Monitor, monitorProfile config.Monitor, settings *config.Settings) {

	if _, ok := monitorProfile.Inputs[command.Input.String()]; !ok {
		fmt.Printf("Input %s not found for monitor %s\n", command.Input.String(), command.Monitor)
		return
	}

	args := ddcutilArgs(monitorProfile)
	fmt.Printf("Attempting to change monitor %s to input %s\n", monitorProfile.Serial, monitorProfile.Inputs[command.Input.String()])
	fmt.Println(settings.Ddcutil.Bin, args[0], args[1], "setvcp", "60", monitorProfile.Inputs[command.Input.String()])

	runDdcutil(settings.Ddcutil.Bin, append(args, "setvcp", "60", monitorProfile.Inputs[command.Input.String()]))
}

func PowerHandler(command commands.Monitor, monitorProfile config.Monitor, settings *config.Settings) {
	if monitorProfile.Display == "" {
		fmt.Printf("Display not found for monitor %s\n", command.Monitor)
		return
	}

	fmt.Printf("Attempting to change monitor %s to power %s\n", monitorProfile.Serial, command.Power.String())

	args := ddcutilArgs(monitorProfile)
	switch command.Power {
	case monitor.On:
		runDdcutil(settings.Ddcutil.Bin, append(args, "setvcp", "0xD6", "0x01"))
	case monitor.Off:
		runDdcutil(settings.Ddcutil.Bin, append(args, "setvcp", "0xD6", "0x04"))
	case monitor.Wake:
		out, err := exec.Command(settings.Xset.Bin, monitorProfile.Display).Output()
		if err != nil {
			fmt.Printf("xset failed: %v\n", err)
			if len(out) > 0 {
				fmt.Println(string(out))
			}
		}
	default:
		fmt.Printf("Invalid power state: %s\n", command.Power.String())
	}
}
