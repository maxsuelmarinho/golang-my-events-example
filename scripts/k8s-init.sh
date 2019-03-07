#!/usr/bin/env bash

kubectl create serviceaccount gitlab-ci
kubectl describe serviceaccount gitlab-ci
serviceAccountToken=$(kubectl describe serviceaccount gitlab 2>&1 | awk '/Tokens:/{print $2}')
tokenEncoded=$(kubectl get secret $serviceAccountToken -o=yaml | awk '/token: /{print $2}')
caCertificateEncoded=$(kubectl get secret $serviceAccountToken -o=yaml | awk '/ca.crt: /{print $2}')
tokenDecoded=$(echo "$tokenEncoded" | base64 --decode)
caCertificateDecoded=$(echo "$caCertificateEncoded" | base64 --decode)

kubectl create clusterrolebinding gitlab-cluster-admin --clusterrole=cluster-admin --serviceaccount=default:gitlab

# DOCKER_USERNAME
# DOCKER_PASSWORD
# KUBE_TOKEN
# KUBE_CA_CERT
# KUBE_SERVER
# KUBE_CLUSTER

kubectl create secret docker-registry my-private-registry \
    --docker-server "https://index.docker.io/v1/" \
    --docker-username "$DOCKER_USERNAME" \
    --docker-password "$DOCKER_PASSWORD" \
    --docker-email "$DOCKER_EMAIL"

kubectl apply -f k8s/prometheus-configmap.yaml
kubectl apply -f k8s/prometheus-statefulset.yaml
kubectl apply -f k8s/prometheus-service.yaml

kubectl apply -f k8s/rabbitmq-statefulset.yaml
kubectl apply -f k8s/rabbitmq-service.yaml

kubectl apply -f k8s/ingress.yaml

kubectl apply -f k8s/bookings-db-statefulset.yaml
kubectl apply -f k8s/bookings-db-service.yaml

kubectl apply -f k8s/events-db-statefulset.yaml
kubectl apply -f k8s/events-db-service.yaml

kubectl apply -f k8s/bookings-deployment.yaml
kubectl apply -f k8s/bookings-service.yaml

kubectl apply -f k8s/events-deployment.yaml
kubectl apply -f k8s/events-service.yaml

