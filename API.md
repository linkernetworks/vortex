# Vortex API

## Network

### Create Network

**POST /v1/networks**

Example:

```
  curl -X POST \
  http://172.17.8.100:7890/v1/networks \
  -H 'Content-Type: application/json' \
  -d '{
   "type":"system",
   "isDPDKPort":false,
   "name":"my network-3",
   "vlanTags":[
      100,
      200
   ],
   "bridgeName":"ovsbr0",
   "nodes":[
      {
         "name":"vortex-dev",
         "physicalInterfaces":[
            {
               "name":"eth0",
               "pciID":""
            }
         ]
      }
   ]
}
'
```

Request Data:

```json
{
   "type":"system",
   "isDPDKPort":false,
   "name":"my network-3",
   "vlanTags":[
      100,
      200
   ],
   "bridgeName":"ovsbr0",
   "nodes":[
      {
         "name":"vortex-dev",
         "physicalInterfaces":[
            {
               "name":"eth0",
               "pciID":""
            }
         ]
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
  "cadvisor-mnq4m": {
   "podName": "cadvisor-mnq4m",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "DaemonSet",
   "createByName": "cadvisor",
   "ip": "10.244.0.61",
   "labels": {
    "controller_revision_hash": "1408846150",
    "name": "cadvisor",
    "pod_template_generation": "1"
   },
   "restartCount": 0,
   "containers": [
    "cadvisor"
   ]
  },
  "coredns-78fcdf6894-fzr58": {
   "podName": "coredns-78fcdf6894-fzr58",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720256,
   "createByKind": "ReplicaSet",
   "createByName": "coredns-78fcdf6894",
   "ip": "10.244.0.32",
   "labels": {
    "k8s_app": "kube-dns",
    "pod_template_hash": "3497892450"
   },
   "restartCount": 1,
   "containers": [
    "coredns"
   ]
  },
  "coredns-78fcdf6894-lwp85": {
   "podName": "coredns-78fcdf6894-lwp85",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720256,
   "createByKind": "ReplicaSet",
   "createByName": "coredns-78fcdf6894",
   "ip": "10.244.0.36",
   "labels": {
    "k8s_app": "kube-dns",
    "pod_template_hash": "3497892450"
   },
   "restartCount": 1,
   "containers": [
    "coredns"
   ]
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
   ]
  },
  "kube-apiserver-vortex-dev": {
   "podName": "kube-apiserver-vortex-dev",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720298,
   "createByKind": "\u003cnone\u003e",
   "createByName": "\u003cnone\u003e",
   "ip": "10.0.2.15",
   "labels": {
    "component": "kube-apiserver",
    "tier": "control-plane"
   },
   "restartCount": 1,
   "containers": [
    "kube-apiserver"
   ]
  },
  "kube-controller-manager-vortex-dev": {
   "podName": "kube-controller-manager-vortex-dev",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720295,
   "createByKind": "\u003cnone\u003e",
   "createByName": "\u003cnone\u003e",
   "ip": "10.0.2.15",
   "labels": {
    "component": "kube-controller-manager",
    "tier": "control-plane"
   },
   "restartCount": 1,
   "containers": [
    "kube-controller-manager"
   ]
  },
  "kube-flannel-ds-t9sv2": {
   "podName": "kube-flannel-ds-t9sv2",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720256,
   "createByKind": "DaemonSet",
   "createByName": "kube-flannel-ds",
   "ip": "10.0.2.15",
   "labels": {
    "app": "flannel",
    "controller_revision_hash": "2856285119",
    "pod_template_generation": "1",
    "tier": "node"
   },
   "restartCount": 2,
   "containers": [
    "kube-flannel"
   ]
  },
  "kube-proxy-sw74x": {
   "podName": "kube-proxy-sw74x",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720256,
   "createByKind": "DaemonSet",
   "createByName": "kube-proxy",
   "ip": "10.0.2.15",
   "labels": {
    "controller_revision_hash": "1151982146",
    "k8s_app": "kube-proxy",
    "pod_template_generation": "1"
   },
   "restartCount": 1,
   "containers": [
    "kube-proxy"
   ]
  },
  "kube-scheduler-vortex-dev": {
   "podName": "kube-scheduler-vortex-dev",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531720313,
   "createByKind": "\u003cnone\u003e",
   "createByName": "\u003cnone\u003e",
   "ip": "10.0.2.15",
   "labels": {
    "component": "kube-scheduler",
    "tier": "control-plane"
   },
   "restartCount": 1,
   "containers": [
    "kube-scheduler"
   ]
  },
  "kube-state-metrics-5fd47f6b7c-9w4pq": {
   "podName": "kube-state-metrics-5fd47f6b7c-9w4pq",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "ReplicaSet",
   "createByName": "kube-state-metrics-5fd47f6b7c",
   "ip": "10.244.0.63",
   "labels": {
    "app": "kube-state-metrics",
    "pod_template_hash": "1980392637"
   },
   "restartCount": 0,
   "containers": [
    "addon-resizer",
    "kube-state-metrics"
   ]
  },
  "mongo-0": {
   "podName": "mongo-0",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488490,
   "createByKind": "StatefulSet",
   "createByName": "mongo",
   "ip": "10.244.0.59",
   "labels": {
    "controller_revision_hash": "mongo-ccff94585",
    "role": "db",
    "service": "mongo",
    "statefulset_kubernetes_io_pod_name": "mongo-0"
   },
   "restartCount": 0,
   "containers": [
    "mongo"
   ]
  },
  "network-controller-server-tcp-bmgm9": {
   "podName": "network-controller-server-tcp-bmgm9",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "DaemonSet",
   "createByName": "network-controller-server-tcp",
   "ip": "10.0.2.15",
   "labels": {
    "controller_revision_hash": "3079531314",
    "name": "network-controller-server-tcp",
    "pod_template_generation": "1"
   },
   "restartCount": 0,
   "containers": [
    "network-controller-server-tcp"
   ]
  },
  "network-controller-server-unix-kcl9s": {
   "podName": "network-controller-server-unix-kcl9s",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "DaemonSet",
   "createByName": "network-controller-server-unix",
   "ip": "10.0.2.15",
   "labels": {
    "controller_revision_hash": "4094250391",
    "name": "network-controller-server-unix",
    "pod_template_generation": "1"
   },
   "restartCount": 0,
   "containers": [
    "network-controller-server-unix"
   ]
  },
  "node-exporter-mjv8p": {
   "podName": "node-exporter-mjv8p",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "DaemonSet",
   "createByName": "node-exporter",
   "ip": "10.0.2.15",
   "labels": {
    "controller_revision_hash": "3813803286",
    "name": "node-exporter",
    "pod_template_generation": "1"
   },
   "restartCount": 0,
   "containers": [
    "node-exporter"
   ]
  },
  "prometheus-7f759794cb-dnbt5": {
   "podName": "prometheus-7f759794cb-dnbt5",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "ReplicaSet",
   "createByName": "prometheus-7f759794cb",
   "ip": "10.244.0.62",
   "labels": {
    "app": "prometheus",
    "pod_template_hash": "3931535076"
   },
   "restartCount": 0,
   "containers": [
    "prometheus"
   ]
  },
  "tiller-deploy-759cb9df9-wp229": {
   "podName": "tiller-deploy-759cb9df9-wp229",
   "namespace": "kube-system",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1531735720,
   "createByKind": "ReplicaSet",
   "createByName": "tiller-deploy-759cb9df9",
   "ip": "10.244.0.34",
   "labels": {
    "app": "helm",
    "name": "tiller",
    "pod_template_hash": "315765895"
   },
   "restartCount": 1,
   "containers": [
    "tiller"
   ]
  },
  "vortex-server-6945b797bb-klx7f": {
   "podName": "vortex-server-6945b797bb-klx7f",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "ReplicaSet",
   "createByName": "vortex-server-6945b797bb",
   "ip": "10.244.0.60",
   "labels": {
    "app": "vortex-server",
    "pod_template_hash": "2501635366"
   },
   "restartCount": 0,
   "containers": [
    "vortex-server"
   ]
  }
 }
```

Example
```
curl -X GET http://localhost:7890/v1/monitoring/pods?namespace=vortex\&node\=vortex-dev\&controller\=prometheus
```

Response Data:
``` json
{
  "prometheus-7f759794cb-dnbt5": {
   "podName": "prometheus-7f759794cb-dnbt5",
   "namespace": "vortex",
   "node": "vortex-dev",
   "status": "Running",
   "createAt": 1532488491,
   "createByKind": "ReplicaSet",
   "createByName": "prometheus-7f759794cb",
   "ip": "10.244.0.62",
   "labels": {
    "app": "prometheus",
    "pod_template_hash": "3931535076"
   },
   "restartCount": 0,
   "containers": [
    "prometheus"
   ]
  }
 }
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
{
  "addon-resizer": {
   "detail": {
    "containerName": "addon-resizer",
    "createAt": 0,
    "pod": "kube-state-metrics-5fd47f6b7c-9w4pq",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/addon-resizer:1.7",
    "command": [
     "/pod_nanny",
     "--container=kube-state-metrics",
     "--cpu=100m",
     "--extra-cpu=1m",
     "--memory=100Mi",
     "--extra-memory=2Mi",
     "--threshold=5",
     "--deployment=kube-state-metrics"
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
    "cpuUsagePercentage": 0.013569055,
    "memoryUsageBytes": 5459968
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "cadvisor": {
   "detail": {
    "containerName": "cadvisor",
    "createAt": 0,
    "pod": "cadvisor-mnq4m",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "google/cadvisor:latest",
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
    "cpuUsagePercentage": 14.916666,
    "memoryUsageBytes": 64532480
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 159197,
    "transmitBytesTotal": 27554169,
    "receivePacketsTotal": 2054,
    "transmitPacketsTotal": 2107
   }
  },
  "coredns": {
   "detail": {
    "containerName": "coredns",
    "createAt": 0,
    "pod": "coredns-78fcdf6894-fzr58",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/coredns:1.1.3",
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
    "cpuUsagePercentage": 0.6946377,
    "memoryUsageBytes": 12357632
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 107847759,
    "transmitBytesTotal": 87997698,
    "receivePacketsTotal": 399620,
    "transmitPacketsTotal": 334664
   }
  },
  "etcd": {
   "detail": {
    "containerName": "etcd",
    "createAt": 0,
    "pod": "etcd-vortex-dev",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/etcd-amd64:3.2.18",
    "command": [
     "etcd",
     "--advertise-client-urls=https://127.0.0.1:2379",
     "--cert-file=/etc/kubernetes/pki/etcd/server.crt",
     "--client-cert-auth=true",
     "--data-dir=/var/lib/etcd",
     "--initial-advertise-peer-urls=https://127.0.0.1:2380",
     "--initial-cluster=vortex-dev=https://127.0.0.1:2380",
     "--key-file=/etc/kubernetes/pki/etcd/server.key",
     "--listen-client-urls=https://127.0.0.1:2379",
     "--listen-peer-urls=https://127.0.0.1:2380",
     "--name=vortex-dev",
     "--peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt",
     "--peer-client-cert-auth=true",
     "--peer-key-file=/etc/kubernetes/pki/etcd/peer.key",
     "--peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt",
     "--snapshot-count=10000",
     "--trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 1.5130409,
    "memoryUsageBytes": 306180100
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-apiserver": {
   "detail": {
    "containerName": "kube-apiserver",
    "createAt": 0,
    "pod": "kube-apiserver-vortex-dev",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/kube-apiserver-amd64:v1.11.0",
    "command": [
     "kube-apiserver",
     "--authorization-mode=Node,RBAC",
     "--advertise-address=172.17.8.100",
     "--allow-privileged=true",
     "--client-ca-file=/etc/kubernetes/pki/ca.crt",
     "--disable-admission-plugins=PersistentVolumeLabel",
     "--enable-admission-plugins=NodeRestriction",
     "--enable-bootstrap-token-auth=true",
     "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
     "--etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt",
     "--etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key",
     "--etcd-servers=https://127.0.0.1:2379",
     "--insecure-port=0",
     "--kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt",
     "--kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key",
     "--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
     "--proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt",
     "--proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key",
     "--requestheader-allowed-names=front-proxy-client",
     "--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt",
     "--requestheader-extra-headers-prefix=X-Remote-Extra-",
     "--requestheader-group-headers=X-Remote-Group",
     "--requestheader-username-headers=X-Remote-User",
     "--secure-port=6443",
     "--service-account-key-file=/etc/kubernetes/pki/sa.pub",
     "--service-cluster-ip-range=10.96.0.0/12",
     "--tls-cert-file=/etc/kubernetes/pki/apiserver.crt",
     "--tls-private-key-file=/etc/kubernetes/pki/apiserver.key"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 2.7560353,
    "memoryUsageBytes": 385847300
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-controller-manager": {
   "detail": {
    "containerName": "kube-controller-manager",
    "createAt": 0,
    "pod": "kube-controller-manager-vortex-dev",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/kube-controller-manager-amd64:v1.11.0",
    "command": [
     "kube-controller-manager",
     "--address=127.0.0.1",
     "--allocate-node-cidrs=true",
     "--cluster-cidr=10.244.0.0/16",
     "--cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt",
     "--cluster-signing-key-file=/etc/kubernetes/pki/ca.key",
     "--controllers=*,bootstrapsigner,tokencleaner",
     "--kubeconfig=/etc/kubernetes/controller-manager.conf",
     "--leader-elect=true",
     "--node-cidr-mask-size=24",
     "--root-ca-file=/etc/kubernetes/pki/ca.crt",
     "--service-account-private-key-file=/etc/kubernetes/pki/sa.key",
     "--use-service-account-credentials=true"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 2.6655617,
    "memoryUsageBytes": 114413570
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-flannel": {
   "detail": {
    "containerName": "kube-flannel",
    "createAt": 0,
    "pod": "kube-flannel-ds-t9sv2",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "quay.io/coreos/flannel:v0.9.1-amd64",
    "command": [
     "/opt/bin/flanneld",
     "--ip-masq",
     "--kube-subnet-mgr",
     "--iface=enp0s8"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 2
   },
   "resource": {
    "cpuUsagePercentage": 0.22413358,
    "memoryUsageBytes": 13639680
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-proxy": {
   "detail": {
    "containerName": "kube-proxy",
    "createAt": 0,
    "pod": "kube-proxy-sw74x",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/kube-proxy-amd64:v1.11.0",
    "command": [
     "/usr/local/bin/kube-proxy",
     "--config=/var/lib/kube-proxy/config.conf"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 0.4472663,
    "memoryUsageBytes": 43425790
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-scheduler": {
   "detail": {
    "containerName": "kube-scheduler",
    "createAt": 0,
    "pod": "kube-scheduler-vortex-dev",
    "namespace": "kube-system",
    "node": "vortex-dev",
    "image": "k8s.gcr.io/kube-scheduler-amd64:v1.11.0",
    "command": [
     "kube-scheduler",
     "--address=127.0.0.1",
     "--kubeconfig=/etc/kubernetes/scheduler.conf",
     "--leader-elect=true"
    ],
    "vNic": ""
   },
   "status": {
    "status": "running",
    "waitingReason": "",
    "terminatedReason": "",
    "restartTime": 1
   },
   "resource": {
    "cpuUsagePercentage": 1.2145014,
    "memoryUsageBytes": 40239104
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "kube-state-metrics": {
   "detail": {
    "containerName": "kube-state-metrics",
    "createAt": 0,
    "pod": "kube-state-metrics-5fd47f6b7c-9w4pq",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/kube-state-metrics:develop",
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
    "cpuUsagePercentage": 0.43355533,
    "memoryUsageBytes": 19472384
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 1731350,
    "transmitBytesTotal": 1758701,
    "receivePacketsTotal": 2010,
    "transmitPacketsTotal": 1817
   }
  },
  "mongo": {
   "detail": {
    "containerName": "mongo",
    "createAt": 0,
    "pod": "mongo-0",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "mongo:4.1-xenial",
    "command": [
     "mongod",
     "--bind_ip",
     "0.0.0.0"
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
    "cpuUsagePercentage": 0.36435637,
    "memoryUsageBytes": 149557250
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 97384,
    "transmitBytesTotal": 75387,
    "receivePacketsTotal": 815,
    "transmitPacketsTotal": 446
   }
  },
  "network-controller-server-tcp": {
   "detail": {
    "containerName": "network-controller-server-tcp",
    "createAt": 0,
    "pod": "network-controller-server-tcp-bmgm9",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/network-controller:v0.3.0",
    "command": [
     "/go/bin/server"
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
    "cpuUsagePercentage": 0,
    "memoryUsageBytes": 3829760
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "network-controller-server-unix": {
   "detail": {
    "containerName": "network-controller-server-unix",
    "createAt": 0,
    "pod": "network-controller-server-unix-kcl9s",
    "namespace": "vortex",
    "node": "vortex-dev",
    "image": "sdnvortex/network-controller:v0.3.0",
    "command": [
     "/go/bin/server"
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
    "cpuUsagePercentage": 0.00016977146,
    "memoryUsageBytes": 4251648
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "node-exporter": {
   "detail": {
    "containerName": "node-exporter",
    "createAt": 0,
    "pod": "node-exporter-mjv8p",
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
    "cpuUsagePercentage": 0.32407588,
    "memoryUsageBytes": 8974336
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 0,
    "transmitBytesTotal": 0,
    "receivePacketsTotal": 0,
    "transmitPacketsTotal": 0
   }
  },
  "prometheus": {
   "detail": {
    "containerName": "prometheus",
    "createAt": 0,
    "pod": "prometheus-7f759794cb-dnbt5",
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
    "cpuUsagePercentage": 2.8768983,
    "memoryUsageBytes": 233467900
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 30859285,
    "transmitBytesTotal": 1263710,
    "receivePacketsTotal": 6266,
    "transmitPacketsTotal": 5795
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
    "cpuUsagePercentage": 0.019468896,
    "memoryUsageBytes": 41377790
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 18227523,
    "transmitBytesTotal": 12689599,
    "receivePacketsTotal": 169295,
    "transmitPacketsTotal": 140420
   }
  },
  "vortex-server": {
   "detail": {
    "containerName": "vortex-server",
    "createAt": 0,
    "pod": "vortex-server-6945b797bb-klx7f",
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
    "cpuUsagePercentage": 0.0129164,
    "memoryUsageBytes": 5525504
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 43017,
    "transmitBytesTotal": 28000,
    "receivePacketsTotal": 267,
    "transmitPacketsTotal": 295
   }
  }
 }
```

Example:
```
curl -X GET http://localhost:7890/v1/monitoring/containers\?namespace\=vortex\&node\=vortex-dev\&pod\=vortex-server-6945b797bb-klx7f
```

Response Data:
``` json
{
  "vortex-server": {
   "detail": {
    "containerName": "vortex-server",
    "createAt": 0,
    "pod": "vortex-server-6945b797bb-klx7f",
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
    "cpuUsagePercentage": 0.015330612,
    "memoryUsageBytes": 5390336
   },
   "nicNetworkTraffic": {
    "receiveBytesTotal": 47243,
    "transmitBytesTotal": 30834,
    "receivePacketsTotal": 290,
    "transmitPacketsTotal": 324
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
  },
  "kubernetes": {
   "serviceName": "kubernetes",
   "namespace": "default",
   "type": "ClusterIP",
   "createAt": 1531720236,
   "clusterIP": "10.96.0.1",
   "ports": [
    {
     "name": "https",
     "protocol": "TCP",
     "port": 443,
     "targetPort": 6443
    }
   ],
   "labels": {
    "component": "apiserver",
    "provider": "kubernetes"
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
  "tiller-deploy": {
   "serviceName": "tiller-deploy",
   "namespace": "kube-system",
   "type": "ClusterIP",
   "createAt": 1531735621,
   "clusterIP": "10.104.53.128",
   "ports": [
    {
     "name": "tiller",
     "protocol": "TCP",
     "port": 44134,
     "targetPort": "tiller"
    }
   ],
   "labels": {
    "app": "helm",
    "name": "tiller"
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
