package handlers

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os/exec"
	"switcher/config"
	"switcher/monitor"
)

func MqttConnect(settings *config.Settings) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", settings.Mqtt.Broker, settings.Mqtt.Port))
	opts.SetClientID("switcher_client")
	opts.SetDefaultPublishHandler(MessageHandler(settings))
	opts.OnConnect = ConnectionHandler(settings)
	opts.OnConnectionLost = ConnectionLostHandler()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func MessageHandler(settings *config.Settings) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		var decodedCommand []monitor.Command
		if err := json.Unmarshal(msg.Payload(), &decodedCommand); err != nil {
			errStr := fmt.Sprintf("Error decoding message: %s\n", err)
			panic(errStr)
		}

		for _, command := range decodedCommand {
			fmt.Printf("Command: %s\n", command)
			var monitorProfile config.Monitor
			ok := false
			if monitorProfile, ok = settings.Monitors[command.Monitor]; !ok {
				errorString := fmt.Sprintf("Monitor %s not found\n", command.Monitor)
				panic(errorString)
			}
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
	}
}

func ConnectionHandler(settings *config.Settings) func(client mqtt.Client) {
	return func(client mqtt.Client) {
		fmt.Println("Connected")

		if token := client.Subscribe(settings.Mqtt.Topic, 1, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Printf("Subscribed to topic: %s\n", settings.Mqtt.Topic)
	}
}

func ConnectionLostHandler() func(client mqtt.Client, err error) {
	return func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
}
