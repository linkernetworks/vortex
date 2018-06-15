FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y curl
RUN curl -fsSL get.docker.com | sh

RUN apt-get install -y sudo
RUN echo '%jenkins ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

## Make /etc/group & /etc/passwd writable for anyuser.
## So, in entry.sh, any user can modify them.
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
RUN apt-get install golang-go

