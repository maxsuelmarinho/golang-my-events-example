package main

import (
	"flag"
	"fmt"
	"golang-my-events-example/events-service/rest"
	"golang-my-events-example/lib/configuration"
	"golang-my-events-example/lib/msgqueue"
	msgqueue_amqp "golang-my-events-example/lib/msgqueue/amqp"
	"golang-my-events-example/lib/persistence/dblayer"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler, eventEmitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
