# API Documentation

Get Stat of Node for quick GUI update

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/stat" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c'
```

### **GET** - /dashboard/v1/{projectID}/node/{nodeID}/stat

### Request object

```
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message
  value             = {
    id              = 6,                         ; Unique Node ID 
    type            = "Leader",                  ; Type fo this node, for example: Leader, Follower, Worker
    status          = "Replicate",               ; progress status of the node, for example: "Replicate", "Probe", "Snapshot" (cluster) or "Running" (nocluster)
    raft            = [ 0, 0 ],                  ; array, index of the last raft entry matched by the node and by the leader, respectively
    load            = [12,0.5,36.52],            ; Array with machine's info: num_clients, server_load, disk_usage_perc
    stats           = [{}]                       ; Array with Time Sample Objects
  },
}
```

#### Time Sample object:

```json
{
  sampletime = "2022-04-28 07:35:12",           ; Date and Time of last activity
  bytes      = {},                              ; Byte Info Object
  clients    = {},                              ; Client Info Object
  commands   = 3,                               ; Number of commands executed
  cpu        = 0.543,                           ; CPU Load value
  io         = {},                              ; IO Info Object
  memory     = {}                               ; Memory Info Object  
}
```

#### Byte Info object:

```json
{
  read       = "123",                           ; Number of bytes read
  writes     = "3186"                           ; Number of bytes written
}
```

#### Clients Info object:

```json
{
  current    = "4",                             ; Number of currently connected clients
  max        = "4"                              ; Maximum number of connected clients
}
```

#### IO Info object:

```json
{
  read       = "123",                           ; Number of bytes read
  writes     = "3186"                           ; Number of bytes written
}
```

#### Memory Info object:

```json
{
  current    = "2652944",                       ; Number of currently used bytes (for Details ask Andrea)
  max        = "2772472"                        ; Maximum number of bytes (for Details ask Andrea)
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/stat HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 07:41:16 GMT
Connection: close
Transfer-Encoding: chunked

{
{
  "message": "OK",
  "value": {
    "id": 1,
    "type": "Leader"
    "status": "Replicate",
    "raft":[
      449,
      449
    ],
    "load": [
        10,
        0.5,
        36.6
    ],
    "stats": [
      {
        "bytes": {
          "reads": "174",
          "writes": "3186"
        },
        "clients": {
          "current": "4",
          "max": "4"
        },
        "commands": "3",
        "cpu": 0.5432,
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "2652944",
          "max": "2772472"
        },
        "sampletime": "2022-04-28 06:41:41"
      },
      ...
  ] },
  "status": 200
}
```