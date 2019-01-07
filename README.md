# Golang My Events Example

## Building applications

**Get dependencies:**
```
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2
go get github.com/streadway/amqp
go get github.com/Shopify/sarama
go get github.com/prometheus/client_golang
```

**Build:**
```
# Dynamically linked
> GOOS=linux GOARCH=amd64 go build
> file event-service
events-service: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, not stripped
> ldd ./events-service
linux-vdso.so.1 (0x00007ffffb558000)
libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f6e62190000)
libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f6e61d90000)
/lib64/ld-linux-x86-64.so.2 (0x00007f6e62400000)

# Static compilation
> CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
> ldd ./events-service
    not a dynamic executable

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
openssl req -x509 -nodes -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -subj "/C=BR/ST=RJ/localityName=Rio de Janeiro/O=MMarinho/organizationalUnitName=MM/commonName=localhost:8181"
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

## React project

**Project initialize:**
```
npm init
```

**Add dependencies:**
```
npm install --save react@16 react-dom@16 @types/react@16 @types/react-dom@16 react-router-dom
npm install --save whatwg-fetch promise-polyfill
npm install --save-dev typescript awesome-typescript-loader source-map-loader webpack webpack-cli http-server @types/react-router-dom
npm install --save bootstrap@^3.3.7
npm install --save whatwg-fetch promise-polyfill
```

**Running Webpack continuously:**
```
webpack --watch
```

**Running web server:**
```
http-server
```

## Docker

**Test Docker installation:**
```
docker container run --rm hello-world
```

**Useful commands:**
```
# List running containers:
docker container ls

# Build a image:
docker image build -t my-image .

# Create a network:
docker network create my-network

# List networks:
docker network ls

# Network connection example
docker container run --rm -d --network=my-network --name=webserver nginx
docker container run --rm -d --network=my-network appropriate/curl http://webserver/

# Create a volume
docker volume create my-volume

# List volumes
docker volume ls

# Volume usage example
docker container run --rm -v my-volume:/my-volume debian:jessie /bin/bash -c "echo Hello > /my-volume/test.txt"
docker container run --rm -v my-volume:/my-volume debian:jessie cat /my-volume/test.txt
```

## Prometheus

**[PromQL:](https://prometheus.io/docs/prometheus/latest/querying/basics/)**
```
go_memstats_alloc_bytes
rate(process_cpu_seconds_total[1m])
sum(myevents_bookings_count) by (eventName)
```