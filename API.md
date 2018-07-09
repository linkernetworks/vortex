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
     -d '{"type":"nfs","name":"My First Storage","nfs":{"ip":"172.17.8.100","path":"/nfs"}}' \
     http://localhost:7890/v1/storage
```

Request Data:
```json
{
	"type": "nfs",
	"name": "My First Storage",
	"nfs": {
		"ip":"172.17.8.100",
		"path":"/nfs"
	}
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
        "nfs": {
            "ip": "172.17.8.100",
            "path": "/nfs"
        }
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
### List Volume
### Remove Volume

