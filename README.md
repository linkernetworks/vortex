Vortex [![Build Status](https://travis-ci.org/linkernetworks/vortex.svg?branch=master)](https://travis-ci.org/linkernetworks/vortex) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/vortex)](https://goreportcard.com/report/github.com/linkernetworks/vortex)
===

# Package sharing between Aurora

Define what package can and can't be shared between aurora, vortex, and 5g-vortex. Since the ownership of vortex source code will be transfered to ITRI or open source.

Stay private:
1. Core aurora service packages like jobserver, jobupdater, ...
2. Kubernetes yamls
3. All aurora API (handler files) contain business logic

Shared packages (will go public)
1. Interface package of public tools like DBs (mongo, influxdb, redis...), logger, json

# Vortex server

Vortex share the same config and dependent services with aurora. Make sure dependent services are available before start vortex server.

- MongoDB
- InfluxDB
- Redis
- Gearmand

### GoBuild

Build
```
make deps vortex
```

Run
```
make run
```

### Docker build

```
make image
```

### Test vortex image

1. Start dependent services like mongo or influxdb
2. Use docker run with host network

```
docker run -it --network=host asia.gcr.io/linker-aurora/vortex:<git-branch> bash
// example
docker run -it --network=host asia.gcr.io/linker-aurora/vortex:develop bash
```
