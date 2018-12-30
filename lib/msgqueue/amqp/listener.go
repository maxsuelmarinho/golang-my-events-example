package amqp

import (
	"fmt"
	amqphelper "golang-my-events-example/lib/helper/amqp"
	"golang-my-events-example/lib/msgqueue"
	"os"
	"time"

	"github.com/streadway/amqp"
)

const eventNameHeader = "x-event-name"

type amqpEventListener struct {
	connection *amqp.Connection
	exchange   string
	queue      string
	mapper     msgqueue.EventMapper
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(a.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not declare queue %s", a.queue, err)
	}
	return nil
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()

	// create binding between queue and exchange for each listened event type
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, a.exchange, false, nil); err != nil {
			return nil, nil, fmt.Errorf("could not bind event %s to queue %s: %s", eventName, a.queue, err)
		}
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not consume queue: %s", err)
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			// Map message to actual event struct
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("message did not contain x-event-name header")
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name header is not string, but %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			event, err := a.mapper.MapEvent(eventName, msg.Body)
			if err != nil {
				errors <- fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
				msg.Nack(false, false)
				continue
			}

			events <- event
			msg.Ack(false)
		}
	}()

	return events, errors, nil
}

func NewAMQPListenerFromEnvironment() (msgqueue.EventListener, error) {
	var url string
	var exchange string
	var queue string

	if url = os.Getenv("AMQP_URL"); url == "" {
		url = "amqp://localhost:5672"
	}

	if exchange = os.Getenv("AMQP_EXCHANGE"); exchange == "" {
		exchange = "example"
	}

	if queue = os.Getenv("AMQP_QUEUE"); queue == "" {
		queue = "example"
	}

	conn := <-amqphelper.RetryConnect(url, 5*time.Second)
	return NewAMQPEventListener(conn, exchange, queue)
}

func NewAMQPEventListener(conn *amqp.Connection, exchange string, queue string) (msgqueue.EventListener, error) {
	listener := amqpEventListener{
		connection: conn,
		exchange:   exchange,
		queue:      queue,
		mapper:     msgqueue.NewEventMapper(),
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

func (a *amqpEventListener) Mapper() msgqueue.EventMapper {
	return a.mapper
}
