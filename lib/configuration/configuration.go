package configuration

import (
	"encoding/json"
	"fmt"
	"golang-my-events-example/lib/persistence/dblayer"
	"os"
)

var (
	DBTypeDefault            = dblayer.DBTYPE("mongodb")
	DBConnectionDefault      = "mongodb://127.0.0.1"
	RestfulEPDefault         = "localhost:8181"
	RestfulTLSEPDefault      = "localhost:9191"
	AMQPMessageBrokerDefault = "amqp://guest:guest@localhost:5672"
)

type ServiceConfig struct {
	Databasetype       dblayer.DBTYPE `json:"databasetype"`
	DBConnection       string         `json:"dbconnection"`
	RestfulEndpoint    string         `json:"restfulapi_endpoint"`
	RestfulTLSEndpoint string         `json:"restfulapi_tlsendpoint"`
	AMQPMessageBroker  string         `json:"amqp_message_broker"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		AMQPMessageBrokerDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
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
		conf.AMQPMessageBroker = broker
	}

	return conf, err
}
