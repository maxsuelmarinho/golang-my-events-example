package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	brokerList := os.Getenv("KAFKA_BROKERS")
	if brokerList == "" {
		brokerList = "localhost:9092"
	}

	brokers := strings.Split(brokerList, ",")
	config := sarama.NewConfig()

	client, err := sarama.NewClient(brokers, config)

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	jsonBody, _ := json.Marshal(&Message{
		Content: "Hello, kafka!",
	})

	msg := &sarama.ProducerMessage{
		Topic: "greeting.topic",
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = producer.SendMessage(msg)
}
