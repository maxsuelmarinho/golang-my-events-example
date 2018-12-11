package main

import (
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

	exchangeName := "events"
	exchangeType := "topic"
	durable := true     // the exchange will remain declared when the broker restarts
	autoDelete := false // the exchange will be deleted as soon as the channel that declared it is closed
	internal := false   // prevent publishers from publishing messages into this queue
	noWait := false     // don't wait for a successful response from the broker
	err = channel.ExchangeDeclare(exchangeName, exchangeType, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}

	routingKey := "some-routing-key"
	mandatory := false // make sure that the message is actually routed into at least one queue
	immediate := false // make sure that the message is actually delivered to at least one subscriber
	err = channel.Publish(exchangeName, routingKey, mandatory, immediate, message)
	if err != nil {
		panic("error while publishing message: " + err.Error())
	}

	defer connection.Close()
}
