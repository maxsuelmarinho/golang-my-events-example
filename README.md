# Golang My Events Example

```
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2
GOOS=linux GOARCH=amd64 go build
```

## Mongo uri connection

```
[mongodb://[user:pass@]host1[:port1][,host2[:port2],...][/database][?options]]
mongodb://127.0.0.1:27017
```

**Certificate**

```
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365
# Or
go run $GOROOT/src/crypto/tls/generate_cert.go --host=localhost [--duration :hours] [--rsa-bits 2048]
```