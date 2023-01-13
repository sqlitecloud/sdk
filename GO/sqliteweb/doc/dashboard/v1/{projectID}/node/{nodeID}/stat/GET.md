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
    stats           = [{}],                      ; Array with Time Sample Objects
    memory          = 1028956160                 ; Total physical memory of the node
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
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with, X-SQLiteCloud-Api-Key
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 12 Jan 2023 15:48:54 GMT
Connection: close
Transfer-Encoding: chunked

{
  "message": "OK",
  "status": 200,
  "value": {
    "id": 6,
    "load": [
      8,
      0.01,
      44.75
    ],
    "memory": 1028956160,
    "raft": [
      12080,
      12080
    ],
    "stats": [
      {
        "bytes": {
          "reads": "3058",
          "writes": "280823"
        },
        "clients": {
          "current": "4",
          "max": "7"
        },
        "commands": "79",
        "cpu": "1.45205823293173",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "6943440",
          "max": "7101592"
        },
        "sampletime": "2023-01-12 14:49:29"
      },
      {
        "bytes": {
          "reads": "3877",
          "writes": "336543"
        },
        "clients": {
          "current": "6",
          "max": "7"
        },
        "commands": "97",
        "cpu": "1.48950904048808",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "6992656",
          "max": "7101592"
        },
        "sampletime": "2023-01-12 14:50:29"
      },
      {
        "bytes": {
          "reads": "3969",
          "writes": "337377"
        },
        "clients": {
          "current": "6",
          "max": "7"
        },
        "commands": "99",
        "cpu": "1.47648303811057",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "6992656",
          "max": "7101592"
        },
        "sampletime": "2023-01-12 14:51:29"
      },
      {
        "bytes": {
          "reads": "4637",
          "writes": "420339"
        },
        "clients": {
          "current": "7",
          "max": "8"
        },
        "commands": "117",
        "cpu": "1.55682776911076",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7002928",
          "max": "7102816"
        },
        "sampletime": "2023-01-12 14:52:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.62983308119012",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:53:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.61353225648556",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:54:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.59858197812649",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:55:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.58456398520573",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:56:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.57107386414755",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:57:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.55786278458844",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:58:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.54597853179684",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 14:59:29"
      },
      {
        "bytes": {
          "reads": "4810",
          "writes": "475289"
        },
        "clients": {
          "current": "8",
          "max": "9"
        },
        "commands": "122",
        "cpu": "1.53377129783694",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7043016",
          "max": "7142904"
        },
        "sampletime": "2023-01-12 15:00:30"
      },
      {
        "bytes": {
          "reads": "5228",
          "writes": "530753"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "133",
        "cpu": "1.55910491071429",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055392",
          "max": "7154056"
        },
        "sampletime": "2023-01-12 15:01:30"
      },
      {
        "bytes": {
          "reads": "5228",
          "writes": "530753"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "133",
        "cpu": "1.54760059429477",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055392",
          "max": "7154056"
        },
        "sampletime": "2023-01-12 15:02:30"
      },
      {
        "bytes": {
          "reads": "5228",
          "writes": "530753"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "133",
        "cpu": "1.53686969814242",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055392",
          "max": "7154056"
        },
        "sampletime": "2023-01-12 15:03:30"
      },
      {
        "bytes": {
          "reads": "5228",
          "writes": "530753"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "133",
        "cpu": "1.52645446293495",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055392",
          "max": "7154056"
        },
        "sampletime": "2023-01-12 15:04:30"
      },
      {
        "bytes": {
          "reads": "5228",
          "writes": "530753"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "133",
        "cpu": "1.51656564349112",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055392",
          "max": "7154056"
        },
        "sampletime": "2023-01-12 15:05:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.54024952966715",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:06:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.53059058073654",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:07:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.52099715672677",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:08:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.51208097826087",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:09:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.50342390146471",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:10:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.49510982375979",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:11:30"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.486636128",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:12:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.47884659340659",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:13:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.47140446841294",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:14:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.46432617246596",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:15:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.45725869242199",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:16:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.45075258394161",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:17:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.44429317073171",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:18:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.43804770098731",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:19:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.43227373092927",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:20:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.42495036734694",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:21:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.42024636241611",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:22:31"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.41451098784997",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:23:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.40916463858554",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:24:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.4042486687148",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:25:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.39940932929904",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:26:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.39456410829608",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:27:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.38996617719041",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:28:32"
      },
      {
        "bytes": {
          "reads": "5471",
          "writes": "568093"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "141",
        "cpu": "1.38541005788712",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:29:32"
      },
      {
        "bytes": {
          "reads": "5889",
          "writes": "623774"
        },
        "clients": {
          "current": "7",
          "max": "10"
        },
        "commands": "152",
        "cpu": "1.4019961959106",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7047136",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:30:32"
      },
      {
        "bytes": {
          "reads": "5889",
          "writes": "623774"
        },
        "clients": {
          "current": "7",
          "max": "10"
        },
        "commands": "152",
        "cpu": "1.39758171589311",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7047136",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:31:32"
      },
      {
        "bytes": {
          "reads": "6476",
          "writes": "679516"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "166",
        "cpu": "1.41154603651491",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7086080",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:32:32"
      },
      {
        "bytes": {
          "reads": "6711",
          "writes": "734856"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "173",
        "cpu": "1.42801358869129",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7114752",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:33:32"
      },
      {
        "bytes": {
          "reads": "6711",
          "writes": "734856"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "173",
        "cpu": "1.4237842780027",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7114752",
          "max": "7192952"
        },
        "sampletime": "2023-01-12 15:34:32"
      },
      {
        "bytes": {
          "reads": "7160",
          "writes": "763639"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "184",
        "cpu": "1.43892296427779",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7083976",
          "max": "7223688"
        },
        "sampletime": "2023-01-12 15:35:33"
      },
      {
        "bytes": {
          "reads": "7409",
          "writes": "818972"
        },
        "clients": {
          "current": "7",
          "max": "10"
        },
        "commands": "192",
        "cpu": "1.45381180205824",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7073704",
          "max": "7223688"
        },
        "sampletime": "2023-01-12 15:36:33"
      },
      {
        "bytes": {
          "reads": "8024",
          "writes": "875480"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "207",
        "cpu": "1.46887594553707",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7112744",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:37:33"
      },
      {
        "bytes": {
          "reads": "8210",
          "writes": "931246"
        },
        "clients": {
          "current": "7",
          "max": "10"
        },
        "commands": "214",
        "cpu": "1.48250486451888",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7073800",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:38:33"
      },
      {
        "bytes": {
          "reads": "8765",
          "writes": "960024"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "227",
        "cpu": "1.49658318938277",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7057504",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:39:33"
      },
      {
        "bytes": {
          "reads": "8765",
          "writes": "960024"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "227",
        "cpu": "1.49151015186187",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7057504",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:40:33"
      },
      {
        "bytes": {
          "reads": "8793",
          "writes": "960858"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "229",
        "cpu": "1.48657672077255",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7057504",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:41:33"
      },
      {
        "bytes": {
          "reads": "9256",
          "writes": "990079"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "241",
        "cpu": "1.5000836005683",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7055400",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:42:33"
      },
      {
        "bytes": {
          "reads": "9491",
          "writes": "1045417"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "248",
        "cpu": "1.51283501102867",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7084072",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:43:33"
      },
      {
        "bytes": {
          "reads": "9941",
          "writes": "1128213"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "260",
        "cpu": "1.52539300574599",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7096448",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:44:33"
      },
      {
        "bytes": {
          "reads": "9941",
          "writes": "1128213"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "260",
        "cpu": "1.52031619345996",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7096448",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:45:33"
      },
      {
        "bytes": {
          "reads": "9941",
          "writes": "1128213"
        },
        "clients": {
          "current": "9",
          "max": "10"
        },
        "commands": "260",
        "cpu": "1.5153710857364",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7096448",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:46:33"
      },
      {
        "bytes": {
          "reads": "10275",
          "writes": "1184123"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "270",
        "cpu": "1.54367320198929",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7057504",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:47:34"
      },
      {
        "bytes": {
          "reads": "10275",
          "writes": "1184123"
        },
        "clients": {
          "current": "8",
          "max": "10"
        },
        "commands": "270",
        "cpu": "1.53860083207262",
        "io": {
          "reads": "0",
          "writes": "0"
        },
        "memory": {
          "current": "7057504",
          "max": "7478824"
        },
        "sampletime": "2023-01-12 15:48:34"
      }
    ],
    "status": "Replicate",
    "type": "Leader"
  }
}
```