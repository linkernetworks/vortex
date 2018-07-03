#!/bin/bash

UID=`id -u`
GID=`id -g`

echo 'jenkins:x:'$GID':' >> /etc/group
echo 'jenkins:x:'$UID':'$GID':,,,:/home/jenkins:/bin/bash' >> /etc/passwd

sudo service docker start
mkdir -p /home/jenkins/data/mongo
nohup mongod --dbpath=/home/jenkins/data/mongo 2>&1 > /dev/null &
docker run --name prometheus -d -p 9090:9090 prom/prometheus
for i in `seq 1 20`
do
    curl http://127.0.0.1:9090/api/v1/query?query=prometheus_build_info > /dev/null
    ret=$?
    if [ "$ret" = "0" ]; then
       break
    fi
    sleep 1
done
sudo chown -R jenkins:jenkins /home/jenkins

bash
