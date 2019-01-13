package builder

import (
	"errors"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue/amqp"
	"log"
	"os"
)

func NewEventListenerFromEnvironment() (msgqueue.EventListener, error) {
	var listener msgqueue.EventListener
	var err error

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting to AMQP broker at %s", url)
		listener, err = amqp.NewEventListenerFromEnvironment()
		if err != nil {
			return nil, err
		}
	} else if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)

		listener, err = kafka.NewKafkaEventListenerFromEnvironment()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Neither AMQP_URL nor KAFKA_BROKERS specified")
	}

	return listener, nil
}
