# Build stage
FROM golang:1.10-alpine3.7
MAINTAINER David Chang <dchang@linkernetworks.com>
# TODO fix path after move out
WORKDIR /go/src/github.com/linkernetworks/vortex

RUN apk add --no-cache protobuf ca-certificates git

# TODO fix path after move out
COPY . /go/src/github.com/linkernetworks/vortex/src
COPY Makefile /go/src/github.com/linkernetworks/vortex

ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Shanghai"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata
# TODO test go install
RUN go install ./src/cmd/...

# the final image: vortex
FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /vortex

# copy the go binaries from the build image
COPY --from=0 /go/bin /go/bin

# copy the config files from the current working dir
COPY config /vortex/config

# select the config file for deployment
ARG CONFIG=config/k8s.json
COPY ${CONFIG} config/k8s.json

EXPOSE 7890
ENTRYPOINT ["/go/bin/vortex", "-port", "7890", "-config", "/vortex/config/k8s.json"]
