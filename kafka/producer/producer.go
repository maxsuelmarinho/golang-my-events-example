package main

import "os"
import "strings"
import "github.com/Shopify/sarama"

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

	
}
