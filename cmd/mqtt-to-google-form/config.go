package main

type Configuration struct {
	Form FormConfiguration
	MQTT MQTTConfiguration
}

type FormConfiguration struct {
	Key          string
	EventField   string
	MessageField string
}

type MQTTConfiguration struct {
	Broker string
}
