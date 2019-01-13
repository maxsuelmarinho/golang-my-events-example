package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/maxsuelmarinho/golang-my-events-example/events-service/rest"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/configuration"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"
	msgqueue_amqp "github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue/amqp"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue/kafka"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/persistence/dblayer"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

func main() {
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", `../lib/configuration/config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*confPath)

	fmt.Println("Connecting to Message Broker")
	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(fmt.Errorf("Could not connect to the RabbitMQ Broker using url %s: %s", config.AMQPMessageBroker, err))
		}

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		if err != nil {
			panic(err)
		}
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		if err != nil {
			panic(fmt.Errorf("Could not connect to the Kafka Broker using url %s: %s", config.KafkaMessageBrokers, err))
		}

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		if err != nil {
			panic(err)
		}
	default:
		panic("Could not connect to broker type: " + config.MessageBrokerType)
	}

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	go func() {
		fmt.Println("Serving metrics API")

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe(":18181", h)
	}()

	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler, eventEmitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
