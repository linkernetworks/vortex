# Vortex API

## Table of Contents
* [User](#user)
   + [Sign Up](#signup)
   + [Sign In](#signin)
   + [Create User](#create-user)
   + [List User](#list-user)
   + [Get User](#get-user)
   + [Delete User](#delete-user)
* [Network](#network)
   + [Create Network](#create-network)
   + [List Network](#list-network)
   + [Get Network](#get-network)
   + [Get Network Status](#get-network-status)
   + [Delete Network](#delete-network)
* [Storage](#storage)
   + [Create Storage](#create-storage)
   + [List Storage](#list-storage)
   + [Remove Storage](#remove-storage)
* [Volume](#volume)
   + [Create Volume](#create-volume)
   + [List Volume](#list-volume)
   + [Remove Volume](#remove-volume)
* [Pod](#pod)
   + [Create Pod](#create-pod)
   + [List Pods](#list-pods)
   + [Get Pod](#get-pod)
   + [Delete Pod](#delete-pod)
* [Resouce Monitoring](#resouce-monitoring)
   + [List Nodes](#list-nodes)
   + [Get Node](#get-node)
   + [List NICs of certain node](#list-nics-of-certain-node)
   + [List Pod](#list-pod)
   + [Get Pod](#get-pod-1)
   + [List Containers](#list-containers)
   + [Get Container](#get-container)
   + [List Services](#list-services)
   + [Get Service](#get-service)
   + [List Controllers](#list-controllers)
   + [Get Controller](#get-controller)
* [Service](#service)
   + [Create Service](#create-service)
   + [List Services](#list-services-1)
   + [Get Service](#get-service-1)
   + [Delete Service](#delete-service)



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
  "vortex-dev": {
   "detail": {
    "hostname": "vortex-dev",
    "createAt": 1531720236,
    "status": "Ready",
    "os": "Ubuntu 16.04.4 LTS",
    "kernelVersion": "4.4.0-130-generic",
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
    "allocatableMemory": 4038606800,
    "allocatablePods": 110,
    "capacityCPU": 2,
    "capacityMemory": 4143464400,
    "capacityPods": 110
   },
   "nics": {
    "cni0": {
     "default": false,
     "type": "virtual",
     "ip": "10.244.0.1/24",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 84888991,
      "transmitBytesTotal": 722809054,
      "receivePacketsTotal": 1009035,
      "transmitPacketsTotal": 1126535
     }
    },
    "docker0": {
     "default": false,
     "type": "virtual",
     "ip": "172.18.0.1/16",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "enp0s3": {
     "default": true,
     "type": "physical",
     "ip": "10.0.2.15/24",
     "pciID": "0000:00:03.0",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 141660005,
      "transmitBytesTotal": 2874858,
      "receivePacketsTotal": 138451,
      "transmitPacketsTotal": 44464
     }
    },
    "enp0s8": {
     "default": false,
     "type": "physical",
     "ip": "172.17.8.100/24",
     "pciID": "0000:00:08.0",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 3290201,
      "transmitBytesTotal": 8253168,
      "receivePacketsTotal": 22636,
      "transmitPacketsTotal": 9259
     }
    },
    "flannel.1": {
     "default": false,
     "type": "virtual",
     "ip": "10.244.0.0/32",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "lo": {
     "default": false,
     "type": "virtual",
     "ip": "127.0.0.1/8",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 3976173174,
      "transmitBytesTotal": 3976173174,
      "receivePacketsTotal": 16152616,
      "transmitPacketsTotal": 16152616
     }
    },
    "veth0e4d10be": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::50c7:43ff:fe27:8ac6/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 70170,
      "transmitBytesTotal": 92158,
      "receivePacketsTotal": 415,
      "transmitPacketsTotal": 757
     }
    },
    "veth6568ca3a": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::cc63:6cff:fe6e:ffae/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 88360079,
      "transmitBytesTotal": 107555728,
      "receivePacketsTotal": 334917,
      "transmitPacketsTotal": 396363
     }
    },
    "veth8423549b": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::d014:f3ff:fe1e:56f2/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 87897057,
      "transmitBytesTotal": 107740789,
      "receivePacketsTotal": 334293,
      "transmitPacketsTotal": 399204
     }
    },
    "veth92e16bdd": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::84c3:baff:fe74:44f2/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 20757600,
      "transmitBytesTotal": 122265,
      "receivePacketsTotal": 1609,
      "transmitPacketsTotal": 1578
     }
    },
    "vethb904323b": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::8f1:cff:fe16:531f/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 22020,
      "transmitBytesTotal": 34320,
      "receivePacketsTotal": 232,
      "transmitPacketsTotal": 217
     }
    },
    "vethbdc03226": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::c46f:42ff:fea2:92af/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 1322545,
      "transmitBytesTotal": 1411885,
      "receivePacketsTotal": 1392,
      "transmitPacketsTotal": 1543
     }
    },
    "vethd5140c96": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::947a:a6ff:fef9:f0d9/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 12678361,
      "transmitBytesTotal": 18212913,
      "receivePacketsTotal": 140293,
      "transmitPacketsTotal": 169129
     }
    },
    "vethf96c51b2": {
     "default": false,
     "type": "virtual",
     "ip": "fe80::d0d5:a3ff:fe26:e143/64",
     "pciID": "",
     "nicNetworkTraffic": {
      "receiveBytesTotal": 1050456,
      "transmitBytesTotal": 23931661,
      "receivePacketsTotal": 4594,
      "transmitPacketsTotal": 5045
     }
    }
   }
  }
 }
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
   "createAt": 1531720236,
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
   "capacityCPU": 2,
   "capacityMemory": 4143472600,
   "capacityPods": 110,
  },
  "nics": {
   "cni0": {
    "default": false,
    "type": "virtual",
    "ip": "10.244.0.1/24",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 60275571,
     "transmitBytesTotal": 329098735,
     "receivePacketsTotal": 842874,
     "transmitPacketsTotal": 933364
    }
   },
   "docker0": {
    "default": false,
    "type": "virtual",
    "ip": "172.18.0.1/16",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 4027243,
     "transmitBytesTotal": 381978119,
     "receivePacketsTotal": 99172,
     "transmitPacketsTotal": 128489
    }
   },
   "enp0s3": {
    "default": true,
    "type": "physical",
    "ip": "10.0.2.15/24",
    "pciID": "0000:00:03.0",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1567333520,
     "transmitBytesTotal": 26003945,
     "receivePacketsTotal": 1333865,
     "transmitPacketsTotal": 414729
    }
   },
   "enp0s8": {
    "default": false,
    "type": "physical",
    "ip": "172.17.8.100/24",
    "pciID": "0000:00:08.0",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 3549805,
     "transmitBytesTotal": 14717952,
     "receivePacketsTotal": 30557,
     "transmitPacketsTotal": 27994
    }
   },
   "flannel.1": {
    "default": false,
    "type": "virtual",
    "ip": "10.244.0.0/32",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 0,
     "transmitBytesTotal": 0,
     "receivePacketsTotal": 0,
     "transmitPacketsTotal": 0
    }
   },
   "lo": {
    "default": false,
    "type": "virtual",
    "ip": "127.0.0.1/8",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 4973678098,
     "transmitBytesTotal": 4973678098,
     "receivePacketsTotal": 21051659,
     "transmitPacketsTotal": 21051659
    }
   },
   "veth0756b817": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::d402:ddff:fee1:924f/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 15795365,
     "transmitBytesTotal": 20571146,
     "receivePacketsTotal": 176548,
     "transmitPacketsTotal": 205045
    }
   },
   "veth0ee29e7": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::4082:94ff:fec0:39f4/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 1559314,
     "transmitBytesTotal": 102459217,
     "receivePacketsTotal": 28724,
     "transmitPacketsTotal": 34721
    }
   },
   "veth1fd22c92": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::98b0:f7ff:fead:c93d/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 122850,
     "transmitBytesTotal": 168478,
     "receivePacketsTotal": 760,
     "transmitPacketsTotal": 1346
    }
   },
   "veth22ed2ac7": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::e0cb:69ff:fe1e:edab/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 33869723,
     "transmitBytesTotal": 120688011,
     "receivePacketsTotal": 309937,
     "transmitPacketsTotal": 345322
    }
   },
   "veth256ca549": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::48f3:9eff:fe5d:9e15/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 33932785,
     "transmitBytesTotal": 120753785,
     "receivePacketsTotal": 309171,
     "transmitPacketsTotal": 346335
    }
   },
   "veth7da58df2": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::3477:fff:fecb:5128/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 100588516,
     "transmitBytesTotal": 595158,
     "receivePacketsTotal": 7769,
     "transmitPacketsTotal": 7651
    }
   },
   "vethbd37bcbc": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::6858:a0ff:fe78:7f6b/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 2000970,
     "transmitBytesTotal": 114194785,
     "receivePacketsTotal": 18975,
     "transmitPacketsTotal": 17595
    }
   },
   "vethddeea13c": {
    "default": false,
    "type": "virtual",
    "ip": "fe80::a8b4:b0ff:fe86:a25f/64",
    "pciID": "",
    "nicNetworkTraffic": {
     "receiveBytesTotal": 6382391,
     "transmitBytesTotal": 5280085,
     "receivePacketsTotal": 6794,
     "transmitPacketsTotal": 7497
    }
   }
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
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "docker0",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "enp0s8",
    "default": false,
    "type": "physical",
    "pciID": "0000:00:08.0"
   },
   {
    "name": "flannel.1",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "lo",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth0756b817",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth0ee29e7",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth1fd22c92",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth22ed2ac7",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth256ca549",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "veth7da58df2",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "vethbd37bcbc",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "vethddeea13c",
    "default": false,
    "type": "virtual",
    "pciID": ""
   },
   {
    "name": "enp0s3",
    "default": true,
    "type": "physical",
    "pciID": "0000:00:03.0"
   }
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
  "busybox": {
   "podName": "busybox",
   "namespace": "default",
   "node": "vortex-dev",
   "status": "Succeeded",
   "createAt": 1532512253,
   "createByKind": "\u003cnone\u003e",
   "createByName": "\u003cnone\u003e",
   "ip": "10.244.0.68",
   "labels": {},
   "restartCount": 0,
   "containers": [
    "busybox"
   ],
   "nics": {}
  },
  "cadvisor-qpsw7": {
   "podName": "cadvisor-qpsw7",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532512322,
   "createByKind": "DaemonSet",
   "createByName": "cadvisor",
   "ip": "10.244.0.71",
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
     "nicNetworkTraffic": {
      "receiveBytesTotal": 756267,
      "transmitBytesTotal": 142368899,
      "receivePacketsTotal": 9735,
      "transmitPacketsTotal": 10008
     }
    }
   }
  },
  "etcd-vortex-dev": {
   "podName": "etcd-vortex-dev",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720301,
   "createByKind": "\u003cnone\u003e",
   "createByName": "\u003cnone\u003e",
   "ip": "10.0.2.15",
   "labels": {
    "component": "etcd",
    "tier": "control-plane"
   },
   "restartCount": 1,
   "containers": [
    "etcd"
   ],
   "nics": {
    "cni0": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 108500943,
      "transmitBytesTotal": 881251596,
      "receivePacketsTotal": 1235729,
      "transmitPacketsTotal": 1383430
     }
    },
    "enp0s3": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 177553664,
      "transmitBytesTotal": 3957292,
      "receivePacketsTotal": 178076,
      "transmitPacketsTotal": 59349
     }
    },
    "enp0s8": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 9779895,
      "transmitBytesTotal": 29102108,
      "receivePacketsTotal": 69328,
      "transmitPacketsTotal": 39141
     }
    },
    "flannel.1": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "ovs-system": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-027715": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-4a7929": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-5ca534": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-63c3d0": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-88964e": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    },
    "system-98f025": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 0,
      "transmitBytesTotal": 0,
      "receivePacketsTotal": 0,
      "transmitPacketsTotal": 0
     }
    }
   }
  }......
```

Example
```
curl -X GET http://localhost:7890/v1/monitoring/pods?namespace=vortex\&node\=vortex-dev\&controller\=prometheus
```

Response Data:
``` json
{
{
  "prometheus-7f759794cb-9h7gr": {
   "podName": "prometheus-7f759794cb-9h7gr",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532512322,
   "createByKind": "ReplicaSet",
   "createByName": "prometheus-7f759794cb",
   "ip": "10.244.0.72",
   "labels": {
    "app": "prometheus",
    "pod_template_hash": "3931535076"
   },
   "restartCount": 0,
   "containers": [
    "prometheus"
   ],
   "nics": {
    "eth0": {
     "nicNetworkTraffic": {
      "receiveBytesTotal": 168663621,
      "transmitBytesTotal": 3313150,
      "receivePacketsTotal": 24828,
      "transmitPacketsTotal": 25897
     }
    }
   }
  }
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
  "podName": "cadvisor-qpsw7",
  "namespace": "vortex",
  "node": "vortex-dev",
  "status": "Running",
  "createAt": 1532512322,
  "createByKind": "DaemonSet",
  "createByName": "cadvisor",
  "ip": "10.244.0.71",
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
    "nicNetworkTraffic": {
     "receiveBytesTotal": 795057,
     "transmitBytesTotal": 149620330,
     "receivePacketsTotal": 10233,
     "transmitPacketsTotal": 10522
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
  "node-exporter": {
   "detail": {
    "containerName": "node-exporter",
    "createAt": 0,
    "pod": "node-exporter-bnmlj",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/node-exporter:develop",
    "command": null,
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 0
   },
   "resource": {
    "cpuUsagePercentage": 0.31914312,
    "memoryUsageBytes": 10280960
   }
  },
  "prometheus": {
   "detail": {
    "containerName": "prometheus",
    "createAt": 0,
    "pod": "prometheus-7f759794cb-9h7gr",
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
    "cpuUsagePercentage": 1.8198096,
    "memoryUsageBytes": 295559170
   }
  },
  "tiller": {
   "detail": {
    "containerName": "tiller",
    "createAt": 0,
    "pod": "tiller-deploy-759cb9df9-wp229",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "gcr.io/kubernetes-helm/tiller:v2.9.1",
    "command": null,
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 0.024541074,
    "memoryUsageBytes": 40525824
   }
  },
  "vortex-server": {
   "detail": {
    "containerName": "vortex-server",
    "createAt": 0,
    "pod": "vortex-server-6945b797bb-jbszk",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/vortex:v0.1.4",
    "command": null,
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 0
   },
   "resource": {
    "cpuUsagePercentage": 0.018635757,
    "memoryUsageBytes": 5439488
   }
  }
 }
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers\?namespace\=vortex\&node\=vortex-dev\&pod\=vortex-server-6945b797bb-jbszk
```

Response Data:
``` json
{
  "vortex-server": {
   "detail": {
    "containerName": "vortex-server",
    "createAt": 0,
    "pod": "vortex-server-6945b797bb-jbszk",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/vortex:v0.1.4",
    "command": null,
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 0
   },
   "resource": {
    "cpuUsagePercentage": 0.0146377785,
    "memoryUsageBytes": 5570560
   }
  }
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
  "kube-dns": {
   "serviceName": "kube-dns",
   "namespace": "kube-system",
   "type": "ClusterIP",
   "createAt": 1531720240,
   "clusterIP": "10.96.0.10",
   "ports": [
    {
     "name": "dns",
     "protocol": "UDP",
     "port": 53,
     "targetPort": 53
    },
    {
     "name": "dns-tcp",
     "protocol": "TCP",
     "port": 53,
     "targetPort": 53
    }
   ],
   "labels": {
    "k8s_app": "kube-dns",
    "kubernetes_io_cluster_service": "true",
    "kubernetes_io_name": "KubeDNS"
   }
  },
  "kube-state-metrics": {
   "serviceName": "kube-state-metrics",
   "namespace": "vortex",
   "type": "ClusterIP",
   "createAt": 1532488491,
   "clusterIP": "10.111.1.210",
   "ports": [
    {
     "name": "http-metrics",
     "protocol": "TCP",
     "port": 8080,
     "targetPort": "http-metrics"
    }
   ],
   "labels": {
    "app": "kube-state-metrics"
   }
  },
  "kubelet": {
   "serviceName": "kubelet",
   "namespace": "kube-system",
   "type": "ClusterIP",
   "createAt": 1531811567,
   "clusterIP": "None",
   "ports": [
    {
     "name": "https-metrics",
     "protocol": "TCP",
     "port": 10250,
     "targetPort": 10250
    }
   ],
   "labels": {
    "k8s_app": "kubelet"
   }
  }......
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/services\?namespace\=monitoring
```

Response Data:
``` json
{
  "kube-state-metrics": {
   "serviceName": "kube-state-metrics",
   "namespace": "vortex",
   "type": "ClusterIP",
   "createAt": 1532488491,
   "clusterIP": "10.111.1.210",
   "ports": [
    {
     "name": "http-metrics",
     "protocol": "TCP",
     "port": 8080,
     "targetPort": "http-metrics"
    }
   ],
   "labels": {
    "app": "kube-state-metrics"
   }
  },
  "mongo": {
   "serviceName": "mongo",
   "namespace": "vortex",
   "type": "NodePort",
   "createAt": 1532488490,
   "clusterIP": "10.110.173.107",
   "ports": [
    {
     "protocol": "TCP",
     "port": 27017,
     "targetPort": 27017,
     "nodePort": 31717
    }
   ],
   "labels": {}
  },
  "prometheus": {
   "serviceName": "prometheus",
   "namespace": "vortex",
   "type": "NodePort",
   "createAt": 1532488491,
   "clusterIP": "10.98.223.167",
   "ports": [
    {
     "protocol": "TCP",
     "port": 9090,
     "targetPort": 9090,
     "nodePort": 30003
    }
   ],
   "labels": {
    "app": "prometheus"
   }
  },
  "vortex-server": {
   "serviceName": "vortex-server",
   "namespace": "vortex",
   "type": "NodePort",
   "createAt": 1532488491,
   "clusterIP": "10.97.86.71",
   "ports": [
    {
     "protocol": "TCP",
     "port": 7890,
     "targetPort": 7890,
     "nodePort": 32326
    }
   ],
   "labels": {
    "app": "vortex-server"
   }
  }
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
  "coredns": {
   "controllerName": "coredns",
   "type": "deployment",
   "namespace": "kube-system",
   "strategy": "",
   "createAt": 1531720240,
   "desiredPod": 2,
   "currentPod": 2,
   "availablePod": 2,
   "labels": {
    "k8s_app": "kube-dns"
   }
  },
  "kube-state-metrics": {
   "controllerName": "kube-state-metrics",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {}
  },
  "prometheus": {
   "controllerName": "prometheus",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {
    "name": "prometheus-deployment"
   }
  },
  "tiller-deploy": {
   "controllerName": "tiller-deploy",
   "type": "deployment",
   "namespace": "kube-system",
   "strategy": "",
   "createAt": 1531735621,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {
    "app": "helm",
    "name": "tiller"
   }
  },
  "vortex-server": {
   "controllerName": "vortex-server",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {
    "app": "vortex-server"
   }
  }
 }
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/controllers\?namespace\=vortex
```

Response Data:
``` json
{
  "kube-state-metrics": {
   "controllerName": "kube-state-metrics",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {}
  },
  "prometheus": {
   "controllerName": "prometheus",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {
    "name": "prometheus-deployment"
   }
  },
  "vortex-server": {
   "controllerName": "vortex-server",
   "type": "deployment",
   "namespace": "vortex",
   "strategy": "",
   "createAt": 1532488491,
   "desiredPod": 1,
   "currentPod": 1,
   "availablePod": 1,
   "labels": {
    "app": "vortex-server"
   }
  }
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
