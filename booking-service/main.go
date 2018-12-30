package main

import (
	"flag"
	"golang-my-events-example/booking-service/listener"
	"golang-my-events-example/booking-service/rest"
	"golang-my-events-example/lib/configuration"
	"golang-my-events-example/lib/msgqueue"
	msgqueue_amqp "golang-my-events-example/lib/msgqueue/amqp"
	"golang-my-events-example/lib/persistence/dblayer"

	"github.com/streadway/amqp"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("config", "./configuration/config.json", "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	panicIfErr(err)

	eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
	panicIfErr(err)

	eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
	panicIfErr(err)

	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	panicIfErr(err)

	processor := listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
