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
ENV PATH=$HOME/go/bin

##
## Go tools
##
RUN go get -u github.com/kardianos/govendor