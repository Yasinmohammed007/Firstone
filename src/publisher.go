package src

import (
	"encoding/json"
	"log"

	"github.com/prnvkv/my-nats/pkg/req-qsub/pub"
	"github.com/prnvkv/my-nats/pkg/util"
)

var NatsTopicDelDNS string = util.GetEnv("NATS_TOPIC_DEL_DNS", "deldnsrecord") // Making this as a global variable
var NatsTopicDNS string = util.GetEnv("NATS_TOPIC_DNS", "dnsrecord")           // Making this as a global variable
var natsTopic string = util.GetEnv("NATS_TOPIC", "config")

type Record struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func PublishRecord(action string, message string) (string, error) {
	data := Record{
		Action:  action,
		Message: message,
	}

	json_val, err := json.Marshal(data)
	if err != nil {
		log.Println("Error Occured while preparing the message", err.Error())
		return "", err
	}

	response, err := PublishMessage(string(json_val))
	if err != nil {
		log.Println("Error Occured while publishing the message", err.Error())
		return "", err
	}

	return response, nil
}

func PublishMessage(message string) (string, error) {
	resp, err := pub.Publish(natsTopic, message)

	if err != nil {
		return "", err
	}
	return string(resp), err
}

// New signature for sending the messages
func PublishMessageToTopic(topic string, message string) (string, error) {
	resp, err := pub.Publish(topic, message)

	if err != nil {
		return "", err
	}
	return string(resp), err
}
