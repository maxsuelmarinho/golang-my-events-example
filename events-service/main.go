package main

import (
	"flag"
	"fmt"
	"golang-my-events-example/events-service/configuration"
	msgqueue_amqp "golang-my-events-example/events-service/msgqueue/amqp"
	"golang-my-events-example/events-service/persistence/dblayer"
	"golang-my-events-example/events-service/rest"
	"log"

	"github.com/streadway/amqp"
)

func main() {

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")

	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*confPath)
	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
