#!/bin/bash

UID=`id -u`
GID=`id -g`

echo 'jenkins:x:'$GID':' >> /etc/group
echo 'jenkins:x:'$UID':'$GID':,,,:/home/jenkins:/bin/bash' >> /etc/passwd

sudo service docker start
mkdir -p /home/jenkins/data/mongo
nohup mongod --dbpath=/home/jenkins/data/mongo 2>&1 > /dev/null &
docker run --name prometheus -d -p 9090:9090 prom/prometheus

sudo chown -R jenkins:jenkins /home/jenkins

bash
