package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/maxsuelmarinho/golang-my-events-example/booking-service/listener"
	"github.com/maxsuelmarinho/golang-my-events-example/booking-service/rest"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/configuration"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"
	msgqueue_amqp "github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue/amqp"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/persistence/dblayer"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

func panicIfErr(err error, message string) {
	if err != nil {
		panic(fmt.Errorf("%s: %s", message, err))
	}
}

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("config", "../lib/configuration/config.json", "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	panicIfErr(err, "Could not connect to the RabbitMQ Broker using url "+config.AMQPMessageBroker)

	eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
	panicIfErr(err, "Could not create new AMQP Event Listener")

	eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
	panicIfErr(err, "Could not create new AMQP Events Emitter")

	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	panicIfErr(err, "Could not connect to database")

	processor := listener.EventProcessor{
		EventListener: eventListener,
		Database:      dbhandler,
	}
	go processor.ProcessEvents()

	go func() {
		fmt.Println("Serving metrics API")

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe(":18282", h)
	}()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
