package handlers

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"switcher/commands"
	"switcher/config"
)

func MqttConnect(settings *config.Settings) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", settings.Mqtt.Broker, settings.Mqtt.Port))
	opts.SetClientID(settings.Mqtt.ClientId)
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
		var decodedCommand []commands.Command
		if err := json.Unmarshal(msg.Payload(), &decodedCommand); err != nil {
			errStr := fmt.Sprintf("Error decoding message: %s\n", err)
			panic(errStr)
		}

		for _, command := range decodedCommand {
			switch command.Command {
			case commands.MonitorCommand:
				MonitorHandler(command.Data, settings)
			case commands.SoundCommand:
				SoundHandler(command.Data, settings)
			case commands.LockCommand:
				LockHandler(command.Data, settings)
			case commands.SpeakerCommand:
				SpeakerHandler(command.Data, settings)
			default:
				fmt.Printf("Unknown command: %s\n", command.Command)
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
