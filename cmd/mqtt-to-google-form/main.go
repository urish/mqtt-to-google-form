package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

var config Configuration

func postForm(event string, msg string) {
	formUrl := fmt.Sprintf("https://docs.google.com/forms/d/%s/formResponse", config.Form.Key)

	resp, err := http.PostForm(
		formUrl,
		url.Values{
			config.Form.EventField:   {event},
			config.Form.MessageField: {msg},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error posting to Google Form: %s", resp.Status)
	}
}

var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	payloadStr := string(msg.Payload())
	log.Default().Printf("Incoming message: %s", payloadStr)
	if strings.HasPrefix(payloadStr, "msg=") {
		postForm("msg", payloadStr[4:])
	}
	if strings.HasPrefix(payloadStr, "status=") {
		postForm("status", payloadStr[7:])
	}
}

func main() {
	fmt.Println("Starting mqtt-to-google-form")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/mqtt-to-google-form/")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Panicf("Error loading config file, %s", err)
	}

	fmt.Println("Configuration:")
	fmt.Println(config)

	mqtt.WARN = log.New(os.Stdout, "WARN ", 0)
	mqtt.ERROR = log.New(os.Stdout, "ERROR ", 0)
	opts := mqtt.NewClientOptions().AddBroker(config.MQTT.Broker).SetClientID("ICOM-client-" + time.Now().Format("2006-01-02 15:04:05"))
	opts.SetDefaultPublishHandler(onMessage)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("ICOM", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	select {}
}
