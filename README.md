# Golang My Events Example
<p align="left">
    <a href="https://travis-ci.org/maxsuelmarinho/golang-my-events-example">
        <img src="https://travis-ci.org/maxsuelmarinho/golang-my-events-example.svg?branch=master" alt="Build Status"></img>
    </a>
    <a href="https://goreportcard.com/report/github.com/maxsuelmarinho/golang-my-events-example"><img src="https://goreportcard.com/badge/github.com/maxsuelmarinho/golang-my-events-example" /></a>
</p>

## Building applications

### Get dependencies

```
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2
go get github.com/streadway/amqp
go get github.com/Shopify/sarama
go get github.com/prometheus/client_golang
```

### Dependency Management (Dep)

**Useful commands:**
```
# Install
> curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Initializing
> dep init

# Add a dependency
> dep ensure -add github.com/gorilla/mux

# Update a dependency
> dep ensure -update github.com/gorilla/mux

# Update all dependencies
> dep ensure -update

# Display the dependencies
> dep status

# Visualizing dependencies
> sudo apt-get install -y graphviz
> dep status -dot | dot -T png | display
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

## Kubernetes

### Kube Control configuration

```
> cat ~/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority: ~/.minikube/ca.crt
    server: https://192.168.99.100:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: ~/.minikube/client.crt
    client-key: ~/.minikube/client.key
```

### Useful commands

**Cluster information:**
```
# Display cluster information
> kubectl cluster-info
Kubernetes master is running at https://192.168.99.100:8443
KubeDNS is running at https://192.168.99.100:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

**Pods:**
```
> kubectl get pods
NAME         READY   STATUS    RESTARTS   AGE
nginx-test   1/1     Running   0          36s

> kubectl apply -f k8s/examples/nginx-pod.yaml

> kubectl get pods
NAME         READY   STATUS              RESTARTS   AGE
nginx-test   0/1     ContainerCreating   0          22s

> kubectl get pods
NAME         READY   STATUS    RESTARTS   AGE
nginx-test   1/1     Running   0          36s

# Delete pod
> kubectl delete pod nginx-test
pod "nginx-test" deleted
```

**Deployment:**
```
> kubectl apply -f k8s/examples/nginx-deployment.yaml
deployment.apps/nginx-deployment created

kubectl get pods
NAME                              READY   STATUS    RESTARTS   AGE
nginx-deployment-7fc5849c-fjgrx   1/1     Running   0          31s
nginx-deployment-7fc5849c-jrvb7   1/1     Running   0          31s

> kubectl edit deployment nginx-deployment

# Add two more instances
> kubectl scale --replicas=4 deployment/nginx-deployment
deployment.extensions/nginx-deployment scaled

> kubectl get pods
NAME                              READY   STATUS    RESTARTS   AGE
nginx-deployment-7fc5849c-2sfgc   1/1     Running   0          19s
nginx-deployment-7fc5849c-fjgrx   1/1     Running   0          4m36s
nginx-deployment-7fc5849c-jrvb7   1/1     Running   0          4m36s
nginx-deployment-7fc5849c-rh2d9   1/1     Running   0          19s
```

**Services:**
```
> kubectl apply -f k8s/examples/nginx-service.yaml
service/nginx created

# List services
> kubectl get services
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP        25h
nginx        NodePort    10.106.36.63   <none>        80:30787/TCP   33s

> minikube service nginx
Opening kubernetes service default/nginx in default browser...
```

**Volumes:**
```
> kubectl apply -f k8s/examples/example-volume.yaml
persistentvolume/volume01 created

> kubectl get pv
NAME       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
volume01   1Gi        RWO,RWX        Retain           Available                                   2m22s

> kubectl apply -f k8s/examples/example-volume-claim.yaml
persistentvolumeclaim/my-data created

> kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM             STORAGECLASS   REASON   AGE
pvc-86970a30-148b-11e9-ba05-0800276224a1   1Gi        RWO            Delete           Bound       default/my-data   standard                46s
volume01                                   1Gi        RWO,RWX        Retain           Available                                             9m50s
```

**RabbitMQ:**
```
> kubectl apply -f k8s/rabbitmq-statefulset.yaml
statefulset.apps/rabbitmq created

> kubectl get pods
NAME                              READY   STATUS              RESTARTS   AGE
rabbitmq-0                        0/1     ContainerCreating   0          38s

> kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                     STORAGECLASS   REASON   AGE
pvc-a6759cef-154a-11e9-ae98-0800276224a1   1Gi        RWO            Delete           Bound       default/data-rabbitmq-0   standard                70s

> kubectl get pvc
NAME              STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-rabbitmq-0   Bound    pvc-a6759cef-154a-11e9-ae98-0800276224a1   1Gi        RWO            standard       102s
```

```
eval $(minikube docker-env)
docker image build -t myevents/eventservice .
```

**Secret:**
```
kubectl create secret docker-registry my-private-registry-credentials \
    --docker-server https://index.docker.io/v1/ \
    --docker-username <my-username>
    --docker-password <my-password>
    --docker-email <my-email>
```

**Ingress:**
```
> minikube addons enable ingress
ingress was successfully enabled

> minikube ip
192.168.99.100

> minikube dashboard
```

## Travis

**Useful commands:**
```
> gem install travis
> travis encrypt DOCKER_PASSWORD="my-secret" --add
```