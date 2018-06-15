FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y curl software-properties-common git
RUN curl -fsSL get.docker.com | sh

##
## Make jenkins can do anything
##
RUN apt-get install -y sudo
RUN echo '%jenkins ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN chmod a+rw /etc/group
RUN chmod a+rw /etc/passwd
RUN echo 'jenkins:!::' >> /etc/gshadow
COPY home /home
ENV HOME=/home/jenkins
COPY entry.sh /entry.sh
RUN chmod +x /entry.sh
ENTRYPOINT /entry.sh

##
## Go
##
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$HOME/go/bin

##
## Go tools
##
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/jstemmer/go-junit-report

RUN mkdir -p /data/db && chmod +rw /data/db
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 2930ADAE8CAF5059EE73BB4B58712A2291FA4AD5
RUN echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.6 multiverse" | \
    tee /etc/apt/sources.list.d/mongodb-org-3.6.list
RUN apt-get update
RUN apt-get install -y mongodb-org