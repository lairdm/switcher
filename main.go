package main

import (
	"fmt"
	"os"
	"os/signal"
	"switcher/config"
	"switcher/handlers"
	"syscall"
	_ "time"
)

func main() {
	settings := config.ParseSettings()
	fmt.Println(settings)

	client := handlers.MqttConnect(settings)
	//	testJson := []byte(`[{"Input":"dp","Monitor":"left"},{"Input":"usbc","Monitor":"right"}]`)

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	//client.Unsubscribe(topic)
	client.Disconnect(250)
}
