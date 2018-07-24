#!/bin/bash
#This script is used to prepare the environment for the travisCI

govendorsync() {
    echo "govendor syncing start `date`"
    govendor sync -v > /tmp/govendor
    echo "govendor syncing end`date`"
    return
}

dockerimages(){
    echo "image pull start `date`"
    sudo docker pull prom/prometheus:v2.2.1 > /tmp/dockerimage
    sudo docker pull hwchiu/ubuntu-nsenter:latest > /tmp/dockerimage
    sudo docker pull mongo:latest > /tmp/dockerimage
    sudo docker pull google/cadvisor:latest > /tmp/dockerimage
    sudo docker pull gcr.io/kubernetes-helm/tiller:v2.9.1 > /tmp/dockerimage
    echo "image pull start `date`"
}

aptget() {
  echo "apt-get start `date`"
  sudo apt-get install -y git build-essential openvswitch-switch socat jq httpie
  sudo add-apt-repository ppa:duggan/bats --yes
  sudo apt-get update -qq
  sudo apt-get install -qq bats
  echo "apt-get end`date`"
}

wgetfiles(){
    echo "download file start `date`"
 curl -s -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
 curl -s -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
 curl -s -L https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz > helm-v2.9.1-linux-amd64.tar.gz && tar -zxvf helm-v2.9.1-linux-amd64.tar.gz && chmod +x linux-amd64/helm && sudo mv linux-amd64/helm /usr/local/bin/helm
    echo "download file end `date`"
}

govendorsync &
dockerimages &
aptget &
wgetfiles &

wait 
