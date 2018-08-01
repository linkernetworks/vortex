# Vortex API

## Table of Contents
- [Vortex API](#vortex-api)
  - [Table of Contents](#table-of-contents)
  - [User](#user)
    - [Signup](#signup)
    - [Signin](#signin)
    - [Create User](#create-user)
    - [List User](#list-user)
    - [Get User](#get-user)
    - [Delete User](#delete-user)
  - [Network](#network)
    - [Create Network](#create-network)
    - [List Network](#list-network)
    - [Get Network](#get-network)
    - [Get Network Status](#get-network-status)
    - [Delete Network](#delete-network)
  - [Storage](#storage)
    - [Create Storage](#create-storage)
    - [List Storage](#list-storage)
    - [Remove Storage](#remove-storage)
  - [Volume](#volume)
    - [Create Volume](#create-volume)
    - [List Volume](#list-volume)
    - [Remove Volume](#remove-volume)
  - [Pod](#pod)
    - [Create Pod](#create-pod)
    - [List Pods](#list-pods)
    - [Get Pod](#get-pod)
    - [Delete Pod](#delete-pod)
  - [Resouce Monitoring](#resouce-monitoring)
    - [List Nodes](#list-nodes)
    - [Get Node](#get-node)
    - [List NICs of certain node](#list-nics-of-certain-node)
    - [List Pod](#list-pod)
    - [Get Pod](#get-pod)
    - [List Containers](#list-containers)
    - [Get Container](#get-container)
    - [List Services](#list-services)
    - [Get Service](#get-service)
    - [List Controllers](#list-controllers)
    - [Get Controller](#get-controller)
  - [Service](#service)
    - [Create Service](#create-service)
    - [List Services](#list-services)
    - [Get Service](#get-service)
    - [Delete Service](#delete-service)



## User

### Signup

**POST /v1/user/signup**

No need to give a role, server will assign a "user" role.

Example:

```json
{
  "loginCredential":{
    "email":"guest@linkernetworks.com",
    "password":"password"
  },
  "username":"John Doe",
  "firstName":"John",
  "lastName":"Doe",
  "phoneNumber":"0911111111"
}
```

Response Data:

```json
{
    "id": "5b5b418c760aab15e771bde2",
    "uuid": "44b4646a-d009-457c-9fdd-1cc0bf226543",
    "jwt": "",
    "loginCredential": {
        "email": "guest@linkernetworks.com",
        "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
    },
    "username": "John Doe",
    "role": "user",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "0911111111",
    "createdAt": "2018-07-28T00:00:12.632011379+08:00"
}
```

### Signin

**POST /v1/users/signin**

Example:

```json
{
    "email":"hello@linkernetworks.com",
    "password":"password"
}
```

Response Data:

```json
{
    "error": false,
    "message": "MY_JWT_TOKEN"
}
```


### Create User

**POST /v1/user**

Example:

role can only be "root", "user", "guest".
```json
{
  "loginCredential":{
    "email":"guest@linkernetworks.com",
    "password":"password"
  },
  "role": "guest",
  "username":"John Doe",
  "firstName":"John",
  "lastName":"Doe",
  "phoneNumber":"0911111111"
}
```

Response Data:

```json
{
    "id": "5b5b418c760aab15e771bde2",
    "uuid": "44b4646a-d009-457c-9fdd-1cc0bf226543",
    "jwt": "",
    "loginCredential": {
        "email": "guest@linkernetworks.com",
        "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
    },
    "username": "John Doe",
    "role": "guest",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "0911111111",
    "createdAt": "2018-07-28T00:00:12.632011379+08:00"
}
```

### List User

Request 

```
GET /v1/users
```


Response Data:

```json
[
    {
        "id": "5b5b4173760aab15e771bde0",
        "uuid": "52870ee9-4bfd-44ea-8cca-a9ce7826b1bd",
        "jwt": "",
        "loginCredential": {
            "email": "root@linkernetworks.com",
            "password": "$2a$14$CQasyFUsBuqwmmpk/i9t9.9j2BTyPzK3PyWATMgb/7g8do57c9oHe"
        },
        "username": "John Doe",
        "role": "root",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-27T23:59:47.564+08:00"
    },
    {
        "id": "5b5b4184760aab15e771bde1",
        "uuid": "a4604f7d-06a8-4226-9792-765e72b14f9c",
        "jwt": "",
        "loginCredential": {
            "email": "user@linkernetworks.com",
            "password": "$2a$14$SzULcUvWqsCy6XeelPdsRutCDJkdsrM4mi2HXpXPEaEugV.jJsMNC"
        },
        "username": "John Doe",
        "role": "user",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-28T00:00:04.261+08:00"
    },
    {
        "id": "5b5b418c760aab15e771bde2",
        "uuid": "44b4646a-d009-457c-9fdd-1cc0bf226543",
        "jwt": "",
        "loginCredential": {
            "email": "guest@linkernetworks.com",
            "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
        },
        "username": "John Doe",
        "role": "guest",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-28T00:00:12.632+08:00"
    }
]
```

### Get User

TODO 


### Delete User

Request

```
DELETE /v1/users/5b5aba2d7a3172bca6f1e280
```

Response Data

``` json
{
    "error": false,
    "message": "User Deleted Success"
}
```


## Network

### Create Network

**POST /v1/networks**

Example:

Request Data:

```json
{
  "type":"system",
  "isDPDKPort":false,
  "name":"my-net",
  "vlanTags":[],
  "nodes":[
    {
      "name":"vortex-dev",
      "physicalInterfaces":[
        {
          "name":"eth0"
        }
      ]
    }
  ]
}
```

Response Data:

```json

{
    "id": "5b5ed39484281d0001ac6735",
    "type": "system",
    "isDPDKPort": false,
    "name": "my-net",
    "vlanTags": [],
    "bridgeName": "system-62fc3f",
    "nodes": [
        {
            "name": "vortex-dev",
            "physicalInterfaces": [
                {
                    "name": "eth0",
                    "pciID": ""
                }
            ]
        }
    ],
    "createdAt": "2018-07-30T09:00:04.740082091Z"
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
[
    {
        "id": "5b47159c4807c50c741c579a",
        "type": "system",
        "isDPDKPort": false,
        "name": "my network-1",
        "vlanTags": [
            100,
            200
        ],
        "routes":[
           {
              "dstCIDR":"224.0.0.0/4",
              "gateway":"0.0.0.0"
           }
        ],
        "bridgeName": "",
        "nodes": [
            {
                "name": "vortex-dev",
                "physicalInterfaces": []
            }
        ],
        "createdAt": "2018-07-12T08:47:24.713Z"
    },
    {
        "id": "5b4716e94807c512d544f437",
        "type": "system",
        "isDPDKPort": false,
        "name": "my network-2",
        "vlanTags": [
            100,
            200
        ],
        "bridgeName": "ovsbr0",
        "nodes": [
            {
                "name": "vortex-dev",
                "physicalInterfaces": []
            }
        ],
        "createdAt": "2018-07-12T08:52:57.567Z"
    }
]
```

### Get Network

**GET /v1/networks/[id]**

Example:

```
curl http://localhost:7890/v1/networks/5b4716e94807c512d544f437
```

Response Data:

```json
{
    "id": "5b4716e94807c512d544f437",
    "type": "system",
    "isDPDKPort": false,
    "name": "my network-2",
    "vlanTags": [
        100,
        200
    ],
    "routes":[
       {
          "dstCIDR":"224.0.0.0/4",
          "gateway":"0.0.0.0"
       }
    ],
    "bridgeName": "ovsbr0",
    "nodes": [
        {
            "name": "vortex-dev",
            "physicalInterfaces": []
        }
    ],
    "createdAt": "2018-07-12T08:52:57.567Z"
}
```

### Get Network Status

This api will return the string array of Pod names and those Pod using the target network and still be running.

**GET /v1/networks/status/[id]**

Example:

```
curl http://localhost:7890/v1/networks/status/5b4716e94807c512d544f437
```

Response Data:

```json
[
    "mypod3",
    "mypod4",
    "mypod5"
]
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

Request Data:
```json
{
	"type": "nfs",
    "name": "My First Storage",
    "ip":"172.17.8.100",
    "path":"/nfs"
}
```
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

Request Data:
```json
{
	"storageName": "My First Storage",
	"name": "My Log",
	"accessMode":"ReadWriteMany",
	"capacity":"300Gi"
}
```

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

For each Pod, we have fileds need to handle.
1. name: the name of the Pod and it should follow the kubernetes yaml rules (Required)
2. labels: the map (string to stirng) for the kubernetes label
3. namespace: the namespace of the Pod.
4. containers: a array of a container (Required)
    - name: the name of the container, it also follow kubernetes naming rule.
    - image: the image of the contaienr.
    - command: a string array, the command of the container.
5. volumes: the array of the voluems that we want to mount to Pod. (Optional)
    - name: the name of the volume and it should be the volume we created before.
    - mountPath: the mountPath of the volume and the container can see files under this path.
6. networks: the array of the network that we want to create in the Pod (Optional)
    - name: the name of the network and it should be the network we created before.
    - ifName: the inteface name you want to create in your container.
    - vlanTag: the vlan tag for `ifName` interface.
    - ipADdress: the IPv4 address of the `ifName` interface.
    - netmask: the IPv4 netmask of the `ifName` interface.
7. capability: the power of the container, if it's ture, it will get almost all capability and act as a privileged=true.
8. restartPolicy: the attribute how the pod restart is container, it should be a string and only valid for those following strings.
    - Always,OnFailure,Never
9. hostNetwork: the bool option to run the Pod in the host network namespace, if it's true, all values in the networks will be ignored.


Example:

Request Data:

```json
{
  "name": "awesome",
  "containers": [{
    "name": "busybox",
    "image": "busybox",
    "command": ["sleep", "3600"]
  }],
  "networks":[
  {
      "name":"MyNetwork2",
      "ifName":"eth12",
      "vlan":0,
      "ipAddress":"1.2.3.4",
      "netmask":"255.255.255.0"
  },
  "volumes":[
  ],
  "capability":true
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
  "namespace": "default",
  "labels": null,
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
  "createdAt": "2018-07-11T06:01:24.637Z",
  "volumes": null,
  "networks": null
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

## Resouce Monitoring

### List Nodes
**GET /v1/monitoring/nodes**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/nodes
```

Response Data:
``` json
{
  "vortex-dev": {...},
  "node1": {...},
  ...
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
   "createAt": 1532573834,
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
   "allocatableMemory": 1891131400,
   "allocatablePods": 110,
   "capacityCPU": 2,
   "capacityMemory": 4143472600,
   "capacityPods": 110
  },
  "nics": {
   "cni0": {
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "ip": "10.244.0.1/24",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": [
      {
       "timestamp": 1532931326.997,
       "value": "1487.6274976818602"
      } ...
     ],
     "transmitBytesTotal": [
      {
       "timestamp": 1532931327.002,
       "value": "6528.226759513464"
      } ...
     ],
     "receivePacketsTotal": [
      {
       "timestamp": 1532931327.006,
       "value": "8.508936201159978"
      } ...
     ],
     "transmitPacketsTotal": [
      {
       "timestamp": 1532931327.011,
       "value": "10.690714714277922"
      } ...
     ]
    }
   },
   "docker0": { ... },
   "enp0s10": { ... },
   ...
  }
 }
```

### List NICs of certain node

**Get /v1/monitoring/nodes/{id}/nics**

Example:
```
curl -X GET  http://localhost:7890/v1/monitoring/nodes/vortex-dev/nics
```

Response Data:
``` json
{
  "nics": [
   {
    "name": "cni0",
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "docker0",
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "dpdk0",
    "default": false,
    "dpdk": true,
    "type": "virtual",
    "pciID": "0000:00:11.0"
   },
   {
    "name": "dpdk1",
    "default": false,
    "dpdk": true,
    "type": "virtual",
    "pciID": "0000:00:12.0"
   },
   {
    "name": "enp0s10",
    "default": false,
    "dpdk": false,
    "type": "physical",
    "pciID": "0000:00:0a.0"
   },
   {
    "name": "enp0s16",
    "default": false,
    "dpdk": false,
    "type": "physical",
    "pciID": "0000:00:10.0"
   },
   {
    "name": "enp0s8",
    "default": false,
    "dpdk": false,
    "type": "physical",
    "pciID": "0000:00:08.0"
   },
   {
    "name": "enp0s9",
    "default": false,
    "dpdk": false,
    "type": "physical",
    "pciID": "0000:00:09.0"
   },
   {
    "name": "flannel.1",
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "lo",
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth67bb7a60",
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "pciID": ""
   } ...
  ]
 }
```

### List Pod
**GET /v1/monitoring/pods?namespace=\.\*&node=\.\*&deployment=\.***

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/pods
```

Response Data:
``` json
{
  "busybox": { ... },
  "etcd-vortex-dev": { ... },
  ...
}
```

Example
```
curl -X GET http://localhost:7890/v1/monitoring/pods?namespace=vortex\&node\=vortex-dev\&controller\=prometheus
```

Response Data:
``` json
{
  "prometheus-7f759794cb-52t54": { ... }
}

```

### Get Pod
**Get /v1/monitoring/pods/{id}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/pods/cadvisor-qpsw7
```

Response Data:
``` json
{
  "podName": "cadvisor-pnpmn",
  "namespace": "vortex",
  "node": "vortex-dev",
  "status": "Running",
  "createAt": 1532931162,
  "createByKind": "DaemonSet",
  "createByName": "cadvisor",
  "ip": "10.244.0.16",
  "labels": {
   "controller_revision_hash": "1408846150",
   "name": "cadvisor",
   "pod_template_generation": "1"
  },
  "restartCount": 0,
  "containers": [
   "cadvisor"
  ],
  "nics": {
   "eth0": {
    "default": false,
    "dpdk": false,
    "type": "virtual",
    "ip": "10.244.0.1/24",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": [
      {
       "timestamp": 1532931969.382,
       "value": "291.60530191458025"
      } ...
     ],
     "transmitBytesTotal": [
      {
       "timestamp": 1532931969.384,
       "value": "55459.517445771744"
      } ...
     ],
     "receivePacketsTotal": [
      {
       "timestamp": 1532931969.386,
       "value": "3.76370479463263"
      } ...
     ],
     "transmitPacketsTotal": [
      {
       "timestamp": 1532931969.388,
       "value": "3.890979835997018"
      } ...
     ]
    }
   }
  }
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
{
  "node-exporter": { ... },
  "prometheus": { ... },
  "tiller": { ... },
  "vortex-server": { ... },
  ...
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers\?namespace\=vortex\&node\=vortex-dev\&pod\=vortex-server-6945b797bb-jbszk
```

Response Data:
``` json
{
  "vortex-server": { ... }
 }
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
   "pod": "prometheus-7f759794cb-52t54",
   "namespace": "vortex",
   "node": "vortex-dev",
   "image": "prom/prometheus:v2.2.1",
   "command": [
    "/bin/prometheus"
   ],
   "vNic": ""
  },
  "status": {
   "status": "running",
   "waitingReason": "",
   "terminatedReason": "",
   "restartTime": 0
  },
  "resource": {
   "cpuUsagePercentage": [
    {
     "timestamp": 1532932286.495,
     "value": "2.1569667381818194"
    } ...
   ],
   "memoryUsageBytes": [
    {
     "timestamp": 1532932286.493,
     "value": "258674688"
    } ...
   ]
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
{
  "kube-dns": { ... },
  "kube-state-metrics": { ... },
  "kubelet": { ... },
  ...
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/services\?namespace\=monitoring
```

Response Data:
``` json
{
  "kube-state-metrics": { ... },
  "mongo": { ... },
  "prometheus": { ... },
  "vortex-server": { ... }
 }
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
{
  "coredns": { ... },
  "kube-state-metrics": { ... },
  "prometheus": { ... },
  ...
 }
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/controllers\?namespace\=vortex
```

Response Data:
``` json
{
  "kube-state-metrics": { ... },
  "prometheus": { ... },
  "vortex-server": { ... }
 }
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

## Service

### Create Service

**POST /v1/services**

Example:

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"name":"awesome","namespace":"default","type":"NodePort","selector":{"podname":"awesome"},"ports":[{"name":"awesome","port":80,"targetPort":80,"nodePort":30000}]}' \
     http://localhost:7890/v1/services
```

Request Data:

```json
{
  "name": "awesome",
  "namespace": "default",
  "type": "NodePort",
  "selector": {
    "podname": "awesome"
  },
  "ports": [
    {
      "name": "awesome",
      "port": 80,
      "targetPort": 80,
      "nodePort": 30000
    }
  ]
}
```

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```

### List Services

**GET /v1/services/**

Example:

```
curl http://localhost:7890/v1/services/
```

Response Data:

```json
[
  {
   "id": "5b4edcbc4807c557d9feb69e",
   "name": "awesome",
   "namespace": "default",
   "type": "NodePort",
   "selector": {
    "podname": "awesome"
   },
   "ports": [
    {
     "name": "awesome",
     "port": 80,
     "targetPort": 80,
     "nodePort": 30000
    }
   ],
   "createdAt": "2018-07-18T06:22:52.403Z"
  }
]
```

### Get Service

**GET /v1/services/[id]**

Example:

```
curl http://localhost:7890/v1/services/5b4edcbc4807c557d9feb69e
```

Response Data:

```json
{
  "id": "5b4edcbc4807c557d9feb69e",
  "name": "awesome",
  "namespace": "default",
  "type": "NodePort",
  "selector": {
   "podname": "awesome"
  },
  "ports": [
   {
    "name": "awesome",
    "port": 80,
    "targetPort": 80,
    "nodePort": 30000
   }
  ],
  "createdAt": "2018-07-18T06:22:52.403Z"
}
```

### Delete Service

**DELETE /v1/services/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/services/5b4edcbc4807c557d9feb69e
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```
