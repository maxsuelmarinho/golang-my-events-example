package kafka

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/maxsuelmarinho/golang-my-events-example/lib/helper/kafka"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"

	"github.com/Shopify/sarama"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

func NewKafkaEventEmitterFromEnvironment() (msgqueue.EventEmitter, error) {
	brokers := []string{"locahost:9092"}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	client := <-kafka.RetryConnect(brokers, 5*time.Second)
	return NewKafkaEventEmitter(client)
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := kafkaEventEmitter{
		producer: producer,
	}

	return &emitter, nil
}

func (e *kafkaEventEmitter) Emit(event msgqueue.Event) error {
	envelope := messageEnvelope{
		event.EventName(),
		event,
	}
	jsonBody, err := json.Marshal(&envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: event.EventName(),
		Value: sarama.ByteEncoder(jsonBody),
	}

	log.Printf("published message with topic %s: %v", event.EventName(), jsonBody)

	//partitionNumber, offset, err = e.producer.SendMessage(msg)
	_, _, err = e.producer.SendMessage(msg)
	return err
}
