# Golang My Events Example

```
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2
go get -u github.com/streadway/amqp
GOOS=linux GOARCH=amd64 go build
```

## MongoDB

**uri connection:** 
```
[mongodb://[user:pass@]host1[:port1][,host2[:port2],...][/database][?options]]
mongodb://127.0.0.1:27017
```

## HTTPS

**Certificate**

```
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365
# Or
go run $GOROOT/src/crypto/tls/generate_cert.go --host=localhost [--duration :hours] [--rsa-bits 2048]
```

## RabbitMQ

**Management UI:**

* **url:** http://[docker-ip]:15672 (**docker-ip:** docker virtual machine ip or localhost)
* **user:** guest
* **pass:** guest

**Other informations:**

* [How to set up a cluster](http://www.rabbitmq.com/clustering.html)
