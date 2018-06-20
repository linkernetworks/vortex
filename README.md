Vortex [![Build Status](https://travis-ci.org/linkernetworks/vortex.svg?branch=develop)](https://travis-ci.org/linkernetworks/vortex) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/vortex)](https://goreportcard.com/report/github.com/linkernetworks/vortex) [![codecov](https://codecov.io/gh/linkernetworks/vortex/branch/develop/graph/badge.svg)](https://codecov.io/gh/linkernetworks/vortex)
===

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
