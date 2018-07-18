#!/bin/bash

UID=`id -u`
GID=`id -g`

echo 'jenkins:x:'$GID':' >> /etc/group
echo 'jenkins:x:'$UID':'$GID':,,,:/home/jenkins:/bin/bash' >> /etc/passwd

sudo service docker start
mkdir -p /home/jenkins/data/mongo
nohup mongod --dbpath=/home/jenkins/data/mongo 2>&1 > /dev/null &
sudo chown -R jenkins:jenkins /home/jenkins

# Make / mounted as rshared to support 
sudo mount --make-rshared /
# Make /sys mounted as rshared to support cadvisor
sudo mount --make-rshared /sys
# Download minikube.
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
sudo env CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none --bootstrapper=localkube --kubernetes-version=v1.9.0 --extra-config=apiserver.Authorization.Mode=RBAC
# Fix the kubectl context, as it's often stale.
minikube update-context
# Enable rbac
kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

# Check if entry.sh finish or not
touch src/github.com/linkernetworks/vortex/ready

bash
