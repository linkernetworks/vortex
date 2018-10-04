# Building stage
FROM golang:1.11-alpine3.7

WORKDIR /go/src/github.com/linkernetworks/vortex

RUN apk add --no-cache protobuf ca-certificates make git

# Source code, building tools and dependences
COPY src /go/src/github.com/linkernetworks/vortex/src
COPY Makefile /go/src/github.com/linkernetworks/vortex
COPY vendor /go/src/github.com/linkernetworks/vortex/vendor

ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Taipei"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata

RUN go get -x -u github.com/kardianos/govendor
RUN govendor sync
RUN make src.build
RUN mv build/src/cmd/vortex/vortex /go/bin

# Production stage
FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /vortex

# copy the go binaries from the building stage
COPY --from=0 /go/bin /go/bin

# copy the config files from the current working dir
COPY config /vortex/config

# select the config file for deployment
ARG CONFIG=config/k8s.json
COPY ${CONFIG} config/k8s.json

EXPOSE 7890
ENTRYPOINT ["/go/bin/vortex", "-port", "7890", "-config", "/vortex/config/k8s.json"]
