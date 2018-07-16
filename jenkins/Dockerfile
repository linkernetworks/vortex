FROM ubuntu:16.04

RUN apt-get update && apt-get install -y curl software-properties-common git socat
RUN curl -fsSL get.docker.com | sh

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
RUN go get -u github.com/t-yuki/gocover-cobertura

##
## kubectl
##
RUN curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl
RUN chmod +x kubectl
RUN mv kubectl /usr/local/bin/

##
## Helm
##
RUN curl -L https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz > helm-v2.9.1-linux-amd64.tar.gz
RUN tar -zxvf helm-v2.9.1-linux-amd64.tar.gz
RUN mv linux-amd64/helm /usr/local/bin/helm

##
## Make jenkins can do anything
##
RUN apt-get install -y sudo
RUN echo '%jenkins ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN chmod a+rw /etc/group
RUN chmod a+rw /etc/passwd
RUN echo 'jenkins:!::' >> /etc/gshadow
COPY home /home
RUN chmod -R a+rw /home
ENV HOME=/home/jenkins
COPY entry.sh /entry.sh
RUN chmod +x /entry.sh
ENTRYPOINT /entry.sh

