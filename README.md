Vortex [![Build Status](https://travis-ci.org/linkernetworks/vortex.svg?branch=develop)](https://travis-ci.org/linkernetworks/vortex) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/vortex)](https://goreportcard.com/report/github.com/linkernetworks/vortex) [![codecov](https://codecov.io/gh/linkernetworks/vortex/branch/develop/graph/badge.svg)](https://codecov.io/gh/linkernetworks/vortex) [![Docker Build Status](https://img.shields.io/docker/build/sdnvortex/vortex.svg)](https://hub.docker.com/r/sdnvortex/vortex/)
===

# Vortex server

![overview](./images/overview.png)

## Features

### Kubernetes

- Kubernetes resource data visualization including network, cpu, memory etc.
    - Nodes
    - Pods
- Kubernetes resources deployment & management
    - Namespaces
    - Deployments, autoscale features provided
    - Services
    - PVs and PVCs with dynamic volume provisioning through NFS
- Getting a shell to a container from web UI
    - Open containers terminal
- Debug and monitor applications on cluster from web UI 
    - Open containers logs
    - Download containers logs
- Users management
    - Create different roles for users
- Private Registry
    - View and pull docker images from private registry

### Network 

- By default, the vortex cluster already installed [flannel](https://github.com/coreos/flannel) for all pods on cluster, however network 
  features can add extra multiple network interfaces for pods
- Custom underlay networking including [Open vSwitch](https://www.openvswitch.org/) and [DPDK](https://www.dpdk.org/) integrations
- Pods multiple network interfaces with static ip and custom route

## Frontend

- [UI Portal](https://github.com/linkernetworks/vortex-portal)

## Backend services

- MongoDB
- InfluxDB
- [Prometheus](https://prometheus.io/)
- [Metrics Server](https://github.com/kubernetes-incubator/metrics-server)
- [CNI Network Controller](https://github.com/linkernetworks/network-controller)
  - Use Open vSwitch as a second bridge for underlay networking
  - Enable Kubernetes Pods have multiple network interfaces and add default routes

## Deploy to bare metal servers (using helm)

```shell
$ make apps.init-helm

# configure private registry url
$ vim config/k8s.json

# configure production yaml 
$ vim deploy/helm/config/production.yaml

$ make apps.launch-prod
```

## Access web UI

```
http://<Kubernetes-Nodes-IP>:32767
```

Default account: admin@vortex.com

Default password: password

## Development and RESTful API endpoint

```
http://<Kubernetes-Nodes-IP>:32326
```

## Upgrade

```shell
$ apps.upgrade-prod
```

## Teardown all system

```shell
$ make apps.teardown-prod
```
