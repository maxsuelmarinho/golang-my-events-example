package main

import (
	"flag"
	"golang-my-events-example/booking-service/listener"
	"golang-my-events-example/lib/configuration"
	"golang-my-events-example/lib/msgqueue"
	msgqueue_amqp "golang-my-events-example/lib/msgqueue/amqp"
	"golang-my-events-example/lib/persistence/dblayer"
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
	config. _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	panicIfErr(err)

	eventListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
	panicIfErr(err)

	eventEmitter, err := msgqueue_amqp.NewAMQPeventEmitter(conn, "events")
	panicIfErr(err)

	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	panicIfErr(err)

	processor := listener.EventProcessr{eventListener, dbhandler}
	go processor.ProcessEvents()

	restServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
