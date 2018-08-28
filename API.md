# Vortex API

## Table of Contents
- [Vortex API](#vortex-api)
  - [Table of Contents](#table-of-contents)
  - [User](#user)
    - [Signup](#signup)
    - [Verify Token](#verify-token)
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
  - [Deployment](#deployment)
    - [Create Deployment](#create-deployment)
    - [List Deployments](#list-deployments)
    - [Get Deployment](#get-deployment)
    - [Delete Deployment](#delete-deployment)
  - [Resouce Monitoring](#resouce-monitoring)
    - [List Nodes](#list-nodes)
    - [Get Node](#get-node)
    - [List NICs of certain node](#list-nics-of-certain-node)
    - [List Pod](#list-pod)
    - [Get Pod](#get-pod)
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
  - [Namespace](#namespace)
    - [Create Namespace](#create-namespace)
    - [List Namespaces](#list-namespaces)
    - [Get Namespace](#get-namespace)
    - [Delete Namespace](#delete-namespace)
  - [OVS](#ovs)
    - [Get PortStats](#get-portstats)
   


## User

### Signup

**POST /v1/user/signup**

No need to give a role, server will assign a "user" role.

Example:

```json
{
  "loginCredential":{
    "username":"guest@linkernetworks.com",
    "password":"password"
  },
  "displayName":"John Doe",
  "firstName":"John",
  "lastName":"Doe",
  "phoneNumber":"0911111111"
}
```

Response Data:

```json
{
    "id": "5b5b418c760aab15e771bde2",
    "loginCredential": {
        "username": "guest@linkernetworks.com",
        "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
    },
    "displayName": "John Doe",
    "role": "user",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "0911111111",
    "createdAt": "2018-07-28T00:00:12.632011379+08:00"
}
```

### Verify Token

**GET /v1/user/verify/auth**

with a authorization JWT key

```
Authorization: Bearer <MY_TOKEN>
```

Response status code:

when jwt is valid return status code 303 redirect to /v1/user/<id>
when jwt is invalid return status code 401 


### Signin

**POST /v1/users/signin**

Example:

```json
{
    "username":"hello@linkernetworks.com",
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

**POST /v1/users**

Example:

role can only be "root", "user", "guest".
```json
{
  "loginCredential":{
    "username":"guest@linkernetworks.com",
    "password":"password"
  },
  "role": "guest",
  "displayName":"John Doe",
  "firstName":"John",
  "lastName":"Doe",
  "phoneNumber":"0911111111"
}
```

Response Data:

```json
{
    "id": "5b5b418c760aab15e771bde2",
    "loginCredential": {
        "username": "guest@linkernetworks.com",
        "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
    },
    "displayName": "John Doe",
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
        "loginCredential": {
            "username": "root@linkernetworks.com",
            "password": "$2a$14$CQasyFUsBuqwmmpk/i9t9.9j2BTyPzK3PyWATMgb/7g8do57c9oHe"
        },
        "displayName": "John Doe",
        "role": "root",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-27T23:59:47.564+08:00"
    },
    {
        "id": "5b5b4184760aab15e771bde1",
        "loginCredential": {
            "username": "user@linkernetworks.com",
            "password": "$2a$14$SzULcUvWqsCy6XeelPdsRutCDJkdsrM4mi2HXpXPEaEugV.jJsMNC"
        },
        "displayName": "John Doe",
        "role": "user",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-28T00:00:04.261+08:00"
    },
    {
        "id": "5b5b418c760aab15e771bde2",
        "loginCredential": {
            "username": "guest@linkernetworks.com",
            "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
        },
        "displayName": "John Doe",
        "role": "guest",
        "firstName": "John",
        "lastName": "Doe",
        "phoneNumber": "0911111111",
        "createdAt": "2018-07-28T00:00:12.632+08:00"
    }
]
```

### Get User

**GET /v1/users/5b5b418c760aab15e771bde2**

Response

```json
{
    "id": "5b5b418c760aab15e771bde2",
    "loginCredential": {
        "username": "guest@linkernetworks.com",
        "password": "$2a$14$XO4OOUCaiTNQHm.ZTzHU5..WwtP2ec2Q2HPPQuMHP1WoXCjXiRrxa"
    },
    "displayName": "John Doe",
    "role": "guest",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "0911111111",
    "createdAt": "2018-07-28T00:00:12.632011379+08:00"
}
```

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
    - routesGw: a array of route with gateway (Optional)
        - dstCIDR(required): destination network cidr for add IP routing table
        - gateway(required): the gateway of the interface subnet
    - routeIntf: a array of route without gateway (Optional)
        - dstCIDR(required): destination network cidr for add IP routing table
7. capability: the power of the container, if it's ture, it will get almost all capability and act as a privileged=true.
8. restartPolicy: the attribute how the pod restart is container, it should be a string and only valid for those following strings.
    - Always,OnFailure,Never
9. networkType: the string options for network type, support "host", "custom" and "cluster".
10. nodeAffinity: the string array to indicate whchi nodes I want my Pod can run in.
11. envVars: the environment variables for containers and it's map (string to stirng) form.

Example:

Request Data:

```json
{  
  "name":"awesome",
  "labels":{},
  "envVars":{},
  "containers":[  
    {  
      "name":"busybox",
      "image":"busybox",
      "command":[  
        "sleep",
        "3600"
      ]
    }
  ],
  "networks":[  
    {  
      "name":"MyNetwork2",
      "ifName":"eth12",
      "vlanTag":0,
      "ipAddress":"1.2.3.4",
      "netmask":"255.255.255.0",
      "routesGw":[  
        {  
          "dstCIDR":"192.168.2.0/24",
          "gateway":"192.168.2.254"
        }
      ],
      "routeIntf":[  
        {  
          "dstCIDR":"224.0.0.0/4"
        }
      ]
    }
  ],
  "volumes":[],
  "capability":true,
  "networkType":"host",
  "nodeAffinity":[]
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

## Deployment

### Create Deployment

**POST /v1/deployments**

For each Deployment, we have fileds need to handle.
1. name: the name of the Deployment and it should follow the kubernetes yaml rules (Required)
2. labels: the map (string to stirng) for the kubernetes label
3. namespace: the namespace of the Deployment.
4. containers: a array of a container (Required)
    - name: the name of the container, it also follow kubernetes naming rule.
    - image: the image of the contaienr.
    - command: a string array, the command of the container.
5. volumes: the array of the voluems that we want to mount to Deployment. (Optional)
    - name: the name of the volume and it should be the volume we created before.
    - mountPath: the mountPath of the volume and the container can see files under this path.
6. networks: the array of the network that we want to create in the Deployment (Optional)
    - name: the name of the network and it should be the network we created before.
    - ifName: the inteface name you want to create in your container.
    - vlanTag: the vlan tag for `ifName` interface.
    - ipADdress: the IPv4 address of the `ifName` interface.
    - netmask: the IPv4 netmask of the `ifName` interface.
    - routesGw: a array of route with gateway (Optional)
        - dstCIDR(required): destination network cidr for add IP routing table
        - gateway(required): the gateway of the interface subnet
    - routeIntf: a array of route without gateway (Optional)
        - dstCIDR(required): destination network cidr for add IP routing table
7. capability: the power of the container, if it's ture, it will get almost all capability and act as a privileged=true.
8.
9. networkType: the string options for network type, support "host", "custom" and "cluster".
10. nodeAffinity: the string array to indicate whchi nodes I want my Deployment can run in.
11. envVars: the environment variables for containers and it's map (string to stirng) form.
12. replicas: the number of the Pods

Example:

Request Data:

```json
{
  "name": "awesome",
  "labels": {},
  "envVars":{},
  "containers": [{
    "name": "busybox",
    "image": "busybox",
    "command": ["sleep", "3600"]
  }],
  "networks":[
  {
      "name":"MyNetwork2",
      "ifName":"eth12",
      "vlanTag":0,
      "ipAddress":"1.2.3.4",
      "netmask":"255.255.255.0",
      "routesGw": [
         {
            "dstCIDR":"192.168.2.0/24",
            "gateway":"192.168.2.254"
         }
      ],
      "routeIntf": [
         {
            "dstCIDR":"224.0.0.0/4",
         }
      ]
  }],
  "volumes":[
  ],
  "capability":true,
  "networkType":"host",
  "nodeAffinity":[],
  "replicas":1
}
```

Response Data:

```json
{
  "error": false,
  "message": "Create success"
}
```

### List Deployments

**GET /v1/deployments/**

Example:

```
curl http://localhost:7890/v1/deployments/
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

### Get Deployment

**GET /v1/deployments/[id]**

Example:

```
curl http://localhost:7890/v1/deployments/5b459d344807c5707ddad740
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

### Delete Deployment

**DELETE /v1/deployments/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/deployments/5b459d344807c5707ddad740
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
   "kernelVersion": "4.4.0-133-generic",
   "dockerVersion": "17.6.2",
   "kubeproxyVersion": "v1.11.0",
   "kubernetesVersion": "v1.11.0",
   "labels": {
    "beta_kubernetes_io_arch": "amd64",
    "beta_kubernetes_io_os": "linux",
    "kubernetes_io_hostname": "vortex-dev"
   }
  },
  "resource": {
   "cpuRequests": 1.3,
   "cpuLimits": 0,
   "memoryRequests": 146800640,
   "memoryLimits": 356515840,
   "memoryTotalHugepages": 1024,
   "memoryFreeHugepages": 512,
   "memoryHugepagesSize": 2097152,
   "allocatableCPU": 2,
   "allocatableMemory": 2948079600,
   "allocatablePods": 110,
   "capacityCPU": 2,
   "capacityMemory": 5200421000,
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
    "type": "physical",
    "pciID": "0000:00:11.0"
   },
   {
    "name": "dpdk1",
    "default": false,
    "dpdk": true,
    "type": "physical",
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

### Get Container
**Get /v1/monitoring/pods/{pod}/{container}**

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers/prometheus
```

Response Data:
``` json
{
  "detail": {
   "containerName": "test1",
   "createAt": 1535361241,
   "status": "running",
   "restartTime": 4,
   "pod": "atest",
   "namespace": "default",
   "node": "vortex-dev",
   "image": "busybox:latest",
   "command": [
    "sleep",
    "3600"
   ]
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
  "namespace": "vortex",
  "strategy": "",
  "createAt": 1535031877,
  "desiredPod": 1,
  "currentPod": 1,
  "availablePod": 1,
  "pods": [
   "prometheus-85bc764c94-kj4wg"
  ],
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

## Namespace
### Create Namespace

**POST /v1/namespaces**

Example:

```
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"awesome"}' \
  http://localhost:7890/v1/namespaces
```

Request Data:

```json
{
  "name": "awesome",
}
```

Response Data:

```json
{
  "id": "5b4edcbc4807c557d9feb69e",
  "name": "awesome",
  "createdAt": "2018-07-18T06:22:52.403Z"
}
```

### List Namespaces

**GET /v1/namespaces/**

Example:

```
curl http://localhost:7890/v1/namespaces/
```

Response Data:

```json
[
  {
   "id": "5b4edcbc4807c557d9feb69e",
   "name": "awesome",
   "createdAt": "2018-07-18T06:22:52.403Z"
  }
]
```

### Get Namespace

**GET /v1/namespaces/[id]**

Example:

```
curl http://localhost:7890/v1/namespaces/5b4edcbc4807c557d9feb69e
```

Response Data:

```json
{
  "id": "5b4edcbc4807c557d9feb69e",
  "name": "awesome",
  "createdAt": "2018-07-18T06:22:52.403Z"
}
```

### Delete Namespace

**DELETE /v1/namespaces/[id]**

Example:

```
curl -X DELETE http://localhost:7890/v1/namespaces/5b4edcbc4807c557d9feb69e
```

Response Data:

```json
{
  "error": false,
  "message": "Delete success"
}
```

## OVS
In the ovs api, we should use two parameter to indicate what OVS we want to operate in.
1. NodeName: the node name in the kubernetes cluster
2. BridgeName: the bridge name when admin create the network in the network page.
the portal can use the list network to fetch the actual bridge name of each network.

### Get PortStats

**GET /v1/ovs/portstats/?nodeName=xxx&bridge=xxx**

Example:

```
curl http://localhost:7890/v1/ovs/portstats?nodeName=vortex-dev&bridgeName=system-47f8ce
```

Response Data:

```json=
{
 {
     "PortID": 2,
         "received": {
             "packets": 0,
             "bytes": 0,
             "dropped": 0,
             "errors": 0
         },
         "traansmitted": {
             "packets": 8,
             "bytes": 648,
             "dropped": 0,
             "errors": 0,
             "collisions": 0
         }
 }
}
```
