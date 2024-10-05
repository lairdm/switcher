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
	opts.OnConnect = ConnectionHandler()
	opts.OnConnectionLost = ConnectionLostHandler()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Subscribe(settings.Mqtt.Topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", settings.Mqtt.Topic)

	return client
}

func MessageHandler(settings *config.Settings) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		var decodedCommand []monitor.Command
		json.Unmarshal(msg.Payload(), &decodedCommand)

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

			cmd := exec.Command(settings.Ddcutil.Bin, "--sn", monitorProfile.Serial, "setvcp", command.Input.String())
			_, err := cmd.Output()
			if err != nil {
				panic(err)
			}
		}
	}
}

func ConnectionHandler() func(client mqtt.Client) {
	return func(client mqtt.Client) {
		fmt.Println("Connected")
	}
}

func ConnectionLostHandler() func(client mqtt.Client, err error) {
	return func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
}
