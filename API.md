# Vortex API

## Network

### Create Network

**POST /v1/networks**

Example:

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"type":"ovs","name":"awesomeNetworks","nodeName":"vortex-dev","ovs":{"bridgeName":"awesomeBridge", "physicalPorts":[]}}' \
     http://localhost:7890/v1/networks
```

Request Data:

```json
{
  "type": "ovs",
  "name": "awesomeNetworks",
  "nodeName": "vortex-dev",
  "ovs": {
    "bridgeName": "awesomeBridge",
    "physicalPorts":[]
  }
}
```

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```

### List Network

**GET /v1/networks/**

Example:

```
curl http://localhost:7890/v1/networks/
```

Response Data:

```json
[{
  "id": "5b3475f94807c5199773910a",
  "type": "ovs",
  "name": "awesomeNetworks",
  "nodeName": "vortex-dev",
  "createdAt": "2018-06-28T05:45:29.828Z",
  "ovs": {
   "bridgeName": "awesomeBridge",
   "physicalPorts": []
  },
  "fake": {
   "bridgeName": "",
   "iWantFail": false
  }
}]
```

### Get Network

**GET /v1/networks/[id]**

Example:

```
curl http://localhost:7890/v1/networks/5b3475f94807c5199773910a
```

Response Data:

```json
{
  "id": "5b3475f94807c5199773910a",
  "type": "ovs",
  "name": "awesomeNetwork",
  "nodeName": "vortex-dev",
  "createdAt": "2018-06-28T05:45:29.828Z",
  "ovs": {
   "bridgeName": "awesomeBridge",
   "physicalPorts": []
  },
  "fake": {
   "bridgeName": "",
   "iWantFail": false
  }
}
```

### Delete Network

**DELETE /v1/networks/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/networks/5b3475f94807c5199773910a
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```

## Storage
### Create Storage

**POST /v1/storage**

Request file:
Type: The storage type we want to connect, it only supoorts `nfs` now.
Name: The name of your storage and it will be used when we want to create the volume.
NFS Parameter:
In the NFS server, there're two parametes we need to provide, the `server IP address` and `exporting path`


Example:

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"type":"nfs","name":"My First Storage","ip":"172.17.8.100","path":"/nfs"}' \
     http://localhost:7890/v1/storage
```

Request Data:
```json
{
	"type": "nfs",
    "name": "My First Storage",
    "ip":"172.17.8.100",
    "path":"/nfs"
}

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```

### List Storage
**GET /v1/storage/**

List all the storages we created before and adding new files.

storageClassName: the storage class name we will used for volume


Example:
```
curl http://localhost:7890/v1/storage/
```

Response Data:

```json
[
    {
        "id": "5b42d9944807c52e1c804fbb",
        "type": "nfs",
        "name": "My First Storage",
        "createdAt": "2018-07-09T03:42:12.708Z",
        "storageClassName": "nfs-storageclass-5b42d9944807c52e1c804fbb",
        "ip": "172.17.8.100",
        "path": "/nfs"
    }
]
```

### Remove Storage
**DELETE /v1/storage/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/storage/5b3475f94807c5199773910a
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```

## Volume
### Create Volume

**POST /v1/volume**

Request file:
storageName: The Storage Name you created before, the system will allocate a space for the volume to use.
accessMode: The accessMode of the Volume including the following options.
- ReadWriteOnce
- ReadWriteMany
- ReeaOneMany
But those options won't work for NFS storage since the permission is controled by the linux permission system.
capacity: The capacity of the volume,

Example:

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"storageName":"My First Storage","name":"My Log","accessMode":"ReadWriteMany","capacity":"300Gi"}' \
     http://localhost:7890/v1/storage
```

Request Data:
```json
{
	"storageName": "My First Storage",
	"name": "My Log",
	"accessMode":"ReadWriteMany",
	"capacity":"300Gi"
}

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```


### List Volume

**GET /v1/volume/**

List all the volumes we created.

storageClassName: the storage class name we will used for volume


Example:
```
curl http://localhost:7890/v1/storage/
```

Response Data:

```json
[
    {
        "id": "5b42f25c4807c52e1c804fbc",
        "name": "My Log",
        "storageName": "My First Storage2",
        "accessMode": "ReadWriteMany",
        "capacity": "300",
        "createdAt": "2018-07-09T05:27:56.244Z"
    }
]
```


### Remove Volume

**DELETE /v1/volume/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/volume/5b3475f94807c5199773910a
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```

## Pod

### Create Pod

**POST /v1/pods**

Example:

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"name":"awesome","containers":[{"name":"busybox","image":"busybox","command":["sleep","3600"]}]}' \
     http://localhost:7890/v1/pods
```

Request Data:

```json
{
  "name": "awesome",
  "containers": [{
    "name": "busybox",
    "image": "busybox",
    "command": ["sleep", "3600"]
  }]
}
```

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```

### List Pods

**GET /v1/pods/**

Example:

```
curl http://localhost:7890/v1/pods/
```

Response Data:

```json
[{
  "id": "5b459d344807c5707ddad740",
  "name": "awesome",
  "containers": [
   {
    "name": "busybox",
    "image": "busybox",
    "command": [
     "sleep",
     "3600"
    ]
   }
  ],
  "createdAt": "2018-07-11T06:01:24.637Z"
}]
```

### Get Pod

**GET /v1/pods/[id]**

Example:

```
curl http://localhost:7890/v1/pods/5b459d344807c5707ddad740
```

Response Data:

```json
{
  "id": "5b459d344807c5707ddad740",
  "name": "awesome",
  "containers": [
   {
    "name": "busybox",
    "image": "busybox",
    "command": [
     "sleep",
     "3600"
    ]
   }
  ],
  "createdAt": "2018-07-11T06:01:24.637Z"
}
```

### Delete Pod

**DELETE /v1/pods/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/pods/5b459d344807c5707ddad740
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```

## Resouce Monitoring

### List Nodes
**GET /v1/monitoring/nodes**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/nodes
```

Response Data:
``` json
[
  "vortex-dev"
]
```

### Get Node
**Get /v1/monitoring/nodes/{id}**

Example
```
curl -X GET http://localhost:7890/v1/monitoring/nodes/vortex-dev
```

Response Data:
``` json
{
  "detail": {
   "hostname": "vortex-dev",
   "createAt": 1530584195,
   "status": "Ready",
   "os": "Ubuntu 16.04.4 LTS",
   "kernelVersion": "4.4.0-128-generic",
   "kubeproxyVersion": "v1.11.0",
   "kubernetesVersion": "v1.11.0",
   "labels": {
    "beta_kubernetes_io_arch": "amd64",
    "beta_kubernetes_io_os": "linux",
    "kubernetes_io_hostname": "vortex-dev"
   }
  },
  "resource": {
   "cpuRequests": 1.05,
   "cpuLimits": 0.6,
   "memoryRequests": 283115520,
   "memoryLimits": 3009413000,
   "allocatableCPU": 2,
   "allocatableMemory": 4038615000,
   "allocatablePods": 110,
   "allocatableEphemeralStorage": 9306748000,
   "capacityCPU": 2,
   "capacityMemory": 4143472600,
   "capacityPods": 110,
   "capacityEphemeralStorage": 10340831000
  },
  "nics": {
   "cni0": {
    "default": "false",
    "type": "virtual",
    "ip": "10.244.0.1/24",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 226772166,
     "transmitBytesTotal": 1383390910,
     "receivePacketsTotal": 2073736,
     "transmitPacketsTotal": 2394274
    }
   },
   "docker0": {
    "default": "false",
    "type": "virtual",
    "ip": "172.18.0.1/16",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1441241,
     "transmitBytesTotal": 1648414,
     "receivePacketsTotal": 10605,
     "transmitPacketsTotal": 18840
    }
   },
   "enp0s3": {
    "default": "true",
    "type": "physical",
    "ip": "10.0.2.15/24",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1279973603,
     "transmitBytesTotal": 17027222,
     "receivePacketsTotal": 1071124,
     "transmitPacketsTotal": 264981
    }
   },
   "enp0s8": {
    "default": "false",
    "type": "physical",
    "ip": "172.17.8.100/24",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 7886367,
     "transmitBytesTotal": 116409318,
     "receivePacketsTotal": 69027,
     "transmitPacketsTotal": 84938
    }
   },
   "flannel.1": {
    "default": "false",
    "type": "virtual",
    "ip": "10.244.0.0/32",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 0,
     "transmitBytesTotal": 0,
     "receivePacketsTotal": 0,
     "transmitPacketsTotal": 0
    }
   },
   "lo": {
    "default": "false",
    "type": "virtual",
    "ip": "127.0.0.1/8",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 7206320651,
     "transmitBytesTotal": 7206320651,
     "receivePacketsTotal": 30220996,
     "transmitPacketsTotal": 30220996
    }
   },
   "veth0369f171": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::f4bf:f1ff:feb0:51e0/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1927350790,
     "transmitBytesTotal": 10996445,
     "receivePacketsTotal": 139639,
     "transmitPacketsTotal": 138994
    }
   },
   "veth09ed0568": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::c430:9cff:fe31:a21a/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 129245700,
     "transmitBytesTotal": 171833140,
     "receivePacketsTotal": 217329,
     "transmitPacketsTotal": 237110
    }
   },
   "veth1542be78": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::d047:ebff:fece:fb3f/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 3237703,
     "transmitBytesTotal": 116003495,
     "receivePacketsTotal": 58301,
     "transmitPacketsTotal": 184394
    }
   },
   "veth17a03665": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::f493:3bff:fe1c:7bfe/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 356941,
     "transmitBytesTotal": 422027,
     "receivePacketsTotal": 4073,
     "transmitPacketsTotal": 2917
    }
   },
   "veth340cca37": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::e4ac:1cff:fe02:e06b/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 5789767,
     "transmitBytesTotal": 355262377,
     "receivePacketsTotal": 55486,
     "transmitPacketsTotal": 50945
    }
   },
   "veth50f06084": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::5882:6dff:fec4:526c/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 16031926,
     "transmitBytesTotal": 20077451,
     "receivePacketsTotal": 182376,
     "transmitPacketsTotal": 223934
    }
   },
   "veth5e506aa3": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::c8fa:5fff:fec8:c7de/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 96955868,
     "transmitBytesTotal": 186560406,
     "receivePacketsTotal": 515276,
     "transmitPacketsTotal": 601733
    }
   },
   "vethb90c3da6": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::8c36:f2ff:fe60:a9e7/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 223475,
     "transmitBytesTotal": 234836,
     "receivePacketsTotal": 1246,
     "transmitPacketsTotal": 2529
    }
   },
   "vethe17ec2c5": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::f8fe:6aff:fee4:b2a6/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 96946925,
     "transmitBytesTotal": 186348463,
     "receivePacketsTotal": 516114,
     "transmitPacketsTotal": 598524
    }
   },
   "vethf0a4680": {
    "default": "false",
    "type": "virtual",
    "ip": "fe80::a00b:b1ff:fe97:d00c/64",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1589711,
     "transmitBytesTotal": 1649062,
     "receivePacketsTotal": 10605,
     "transmitPacketsTotal": 18848
    }
   }
  }
 }
```

### List Pod
**GET /v1/monitoring/pods?namespace=\.\*&node=\.\*&deployment=\.***

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/pods
curl -X GET http://localhost:7890/v1/monitoring/pods?namespace=monitoring\&node\=vortex-dev\&controller\=prometheus
```

Response Data:
``` json
[
  "etcd-vortex-dev",
  "kube-apiserver-vortex-dev",
  "kube-controller-manager-vortex-dev",
  "kube-scheduler-vortex-dev",
  "cadvisor-mktsc",
  "kube-flannel-ds-wrqhd",
  "kube-proxy-5knh8",
  "node-exporter-q2ckj",
  "coredns-78fcdf6894-hxvw2",
  "coredns-78fcdf6894-lbfnd",
  "develop-66855589b7-tzwxd",
  "kube-state-metrics-849d66bcc4-9csb7",
  "prometheus-69dfbf887b-n2zf7",
  "tiller-deploy-759cb9df9-mnkj6",
  "youngling-echidna-vortex-server-6c6dbd8bc8-bb4g2",
  "mongo-0"
 ]
```

Example
```
curl -X GET http://localhost:7890/v1/monitoring/pods?namespace=monitoring\&node\=vortex-dev\&controller\=prometheus
```

Response Data:
``` json
[
  "prometheus-69dfbf887b-n2zf7"
 ]
```

### Get Pod
**Get /v1/monitoring/pods/{id}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/pods/cadvisor-mktsc
```

Response Data:
``` json
{
  "podName": "cadvisor-mktsc",
  "namespace": "monitoring",
  "node": "vortex-dev",
  "status": "Running",
  "createAt": 1531124080,
  "createByKind": "DaemonSet",
  "createByName": "cadvisor",
  "ip": "10.244.0.25",
  "labels": {
   "controller_revision_hash": "3793291166",
   "name": "cadvisor",
   "pod_template_generation": "1"
  },
  "restartCount": 0,
  "containers": [
   "cadvisor"
  ]
 }
```

### List Containers
**GET /v1/monitoring/container?namespace=\.\*&node=\.\*&podpo=\.***

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers
```

Response Data:
``` json
[
  "addon-resizer",
  "cadvisor",
  "coredns",
  "coredns",
  "develop",
  "etcd",
  "kube-apiserver",
  "kube-controller-manager",
  "kube-flannel",
  "kube-proxy",
  "kube-scheduler",
  "kube-state-metrics",
  "mongo",
  "node-exporter",
  "prometheus",
  "tiller",
  "vortex-server"
 ]
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers\?namespace\=monitoring\&node\=vortex-dev\&pod\=kube-state-metrics-849d66bcc4-9csb7
```

Response Data:
``` json
[
  "addon-resizer",
  "kube-state-metrics"
 ]
```

### Get Container
**Get /v1/monitoring/container/{id}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers/prometheus
```

Response Data:
``` json
{
  "detail": {
   "containerName": "prometheus",
   "createAt": 0,
   "pod": "prometheus-69dfbf887b-n2zf7",
   "namespace": "monitoring",
   "node": "vortex-dev",
   "image": "prom/prometheus:v2.2.1",
   "command": [
    "/bin/prometheus"
   ]
  },
  "status": {
   "status": "running",
   "waitingReason": "",
   "terminatedReason": "",
   "restartTime": 0
  },
  "resource": {
   "cpuUsagePercentage": 2.3317826,
   "memoryUsageBytes": 423919600
  },
  "nicNetworkTraffic": {
   "receiveBytesTotal": 477653430,
   "transmitBytesTotal": 7745781,
   "receivePacketsTotal": 68539,
   "transmitPacketsTotal": 74386
  }
 }
```

### List Services
**GET /v1/monitoring/service?namespace=\.\***

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/services
```

Response Data:
``` json
[
  "prometheus-external",
  "kube-state-metrics",
  "mongo",
  "tiller-deploy",
  "youngling-echidna-vortex-server",
  "kubernetes",
  "kube-dns",
  "prometheus"
 ]
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/services\?namespace\=monitoring
```

Response Data:
``` json
[
  "prometheus-external",
  "kube-state-metrics",
  "mongo",
  "youngling-echidna-vortex-server",
  "prometheus"
 ]
```

### Get Service
**Get /v1/monitoring/service/{id}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/services/mongo
```

Response Data:
``` json
{
  "serviceName": "mongo",
  "namespace": "monitoring",
  "type": "ClusterIP",
  "createAt": 1531196180,
  "clusterIP": "10.107.88.103",
  "Ports": [
   {
    "protocol": "TCP",
    "port": 27017,
    "targetPort": 27017
   }
  ],
  "labels": {
   "name": "mongo",
   "service": "mongo"
  }
 }
```

### List Controllers
**GET /v1/monitoring/controller?namespace=\.\***

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/controllers
```

Response Data:
``` json
[
  "coredns",
  "develop",
  "kube-state-metrics",
  "prometheus",
  "tiller-deploy",
  "youngling-echidna-vortex-server"
 ]
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/controllers\?namespace\=monitoring
```

Response Data:
``` json
[
  "kube-state-metrics",
  "prometheus",
  "youngling-echidna-vortex-server"
 ]
```

### Get Controller
**Get /v1/monitoring/controller/{id}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/controllers/prometheus
```

Response Data:
``` json
{
  "controllerName": "prometheus",
  "type": "deployment",
  "namespace": "monitoring",
  "strategy": "",
  "createAt": 1531211728,
  "desiredPod": 1,
  "currentPod": 1,
  "availablePod": 1,
  "labels": {
   "name": "prometheus-deployment"
  }
 }
```