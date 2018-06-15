#!/bin/bash

UID=`id -u`
GID=`id -g`

echo 'jenkins:x:'$GID':' >> /etc/group
echo 'jenkins:x:'$UID':'$GID':,,,:/home/jenkins:/bin/bash' >> /etc/passwd

sudo service docker start
nohup mongod 2>&1 > /dev/null &

sudo chown -R jenkins:jenkins /home/jenkins
mkdir /home/jenkins/.cache/
chmod -R +rw /home/jenkins/.cache/

bash
