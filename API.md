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
