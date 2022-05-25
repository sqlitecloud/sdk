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
    id              = 1,                         ; Unique node id 
    name            = "Dev1 Server",             ; Name of this node
    node_id         = 1,                         ; node_id of the machine in the cluster
    provider        = "DigitalOcean",            ; Provider of this node
    region          = "Rome/Italy",              ; Regin data for this node
    size            = "small",                   ; Size info for this node
  },
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
    "type": "worker"
  },
  "status": 200
}
```