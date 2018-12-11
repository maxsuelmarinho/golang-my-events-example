package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672"
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("could not establish AMQP connection: " + err.Error())
	}

	channel, err := connection.Channel()
	if err != nil {
		panic("could not open channel: " + err.Error())
	}

	queueName := "my_queue"
	durable := true     // the exchange will remain declared when the broker restarts
	autoDelete := false // the exchange will be deleted as soon as the channel that declared it is closed
	exclusive := false  // this consumer will be the only one allowed to consume this queue, otherwise other consumers might listen on the same queue.
	noWait := false     // don't wait for a successful response from the broker

	_, err = channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		panic("error while declaring the queue: " + err.Error())
	}

	exchangeName := "events"
	err = channel.QueueBind(queueName, "#", exchangeName, noWait, nil)
	if err != nil {
		panic("error while binding the queue: " + err.Error())
	}

	autoAck := false // received messages will be acknowledged automatically
	noLocal := false // indicate to the broker that this consumer should not be delivered messages that were published on the same channel
	consumerId := "" // when left blank, a unique identifier will be automatically generated
	msgs, err := channel.Consume(queueName, consumerId, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		panic("error while consuming the queue: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}

	defer connection.Close()
}
