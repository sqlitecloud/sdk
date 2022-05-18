# API Documentation

Get a JSON with all providers, regions and size parameters

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  value             = {
    address         = "127.0.0.1",               ; IPv[4,6] address or host name of this node
    port            = 9960,                      ; Port this node is listening on
    details         = "?/?/?",                   ; "SFO1/1GB/25GB disk or "i386/1/1MB/100MB",
    latidude        = float,
    longitude       = float,
    load            = [ float, float ],          ; some load info (for Details ask Andrea)
    id              = 1,                         ; NodeID - It is not good to have a simple int number!!!!!!
    name            = "Dev1 Server",             ; Name of this node
    node_id         = 1,
    provider        = "DigitalOcean",            ; Provider of this node
    raft            = [ int, int ],              ; array 8960, 8960 (for Details ask Andrea)
    region          = "Rome/Italy",              ; Regin data for this node
    size            = "small",                   ; Size info for this node
    stats           = [{}],                      ; Array with Time Sample Objects
    type            = "worker",                  ; Type fo this node, for example: Leader, Worker
    status          = "Unknown",                 ; Replicating
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
  cpu        = 0.543,                           ; CPU Load Value
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
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 08:13:15 GMT
Connection: close
Transfer-Encoding: chunked

{
  "message": "OK",
  "value": {
    "address": "64.227.11.116",
    "details": "i386/1/1MB/100MB",
    "id": 1,
    "latitude": 40.793,
    "load": [
      1,
      0.5
    ],
    "longitude": -74.0247,
    "name": "Dev1 Server",
    "node_id": 1,
    "port": 9960,
    "provider": "DigitalOcean",
    "raft": [
      0,
      0
    ],
    "region": "Rome/Italy",
    "size": "small",
    "stats": [
      {
        "bytes": {
          "reads": "864",
          "writes": "164262"
        },
        "clients": {
          "current": "4",
          "max": "4"
        },
        "commands": "27",
        "cpu": 0.5432,
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "4534280",
          "max": "4682488"
        },
        "sampletime": "2022-04-28 07:13:42"
      },
      ...
    } ],
    "status": "Unknown",
    "type": "worker"
  },
  "status": 200
}
```