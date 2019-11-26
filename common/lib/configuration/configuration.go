package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maxsuelmarinho/golang-my-events-example/common/lib/persistence/dblayer"
)

var (
	DBTypeDefault              = dblayer.DBTYPE("mongodb")
	DBConnectionDefault        = "mongodb://127.0.0.1"
	RestfulEPDefault           = "localhost:8181"
	RestfulTLSEPDefault        = "localhost:9191"
	MessageBrokerTypeDefault   = "amqp"
	AMQPMessageBrokerDefault   = "amqp://guest:guest@localhost:5672"
	KafkaMessageBrokersDefault = []string{"localhost:9092"}
)

type ServiceConfig struct {
	Databasetype        dblayer.DBTYPE `json:"databasetype"`
	DBConnection        string         `json:"dbconnection"`
	RestfulEndpoint     string         `json:"restfulapi_endpoint"`
	RestfulTLSEndpoint  string         `json:"restfulapi_tlsendpoint"`
	MessageBrokerType   string         `json:"message_broker_type"`
	AMQPMessageBroker   string         `json:"amqp_message_broker"`
	KafkaMessageBrokers []string       `json:"kafka_message_brokers"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.", err)
	}

	err = json.NewDecoder(file).Decode(&conf)

	if v := os.Getenv("LISTEN_URL"); v != "" {
		conf.RestfulEndpoint = v
	}

	if v := os.Getenv("MONGO_URL"); v != "" {
		conf.Databasetype = "mongodb"
		conf.DBConnection = v
	}

	if broker := os.Getenv("AMQP_BROKER_URL"); broker != "" {
		conf.MessageBrokerType = "amqp"
		conf.AMQPMessageBroker = broker
	} else if brokers := os.Getenv("KAFKA_BROKER_URLS"); brokers != "" {
		conf.MessageBrokerType = "kafka"
		conf.KafkaMessageBrokers = strings.Split("brokers", ",")
	}

	prettyJSON, err := json.MarshalIndent(conf, "", "\t")
	fmt.Printf("Using Configuration:\n%s\n", prettyJSON)

	return conf, err
}
