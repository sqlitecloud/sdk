# API Documentation

List all userid projects

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

```code
{
  status            = 0,                         -- status code: 0 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  node              = {
    id              = 0,                         -- NodeID - It is not good to have a simple int number!!!!!!
    name            = "",                        -- Name of this node
    type            = "",                        -- Type fo this node, for example: Leader, Worker
    provider        = "",                        -- Provider of this node
    image           = "",                        -- Image data for this node
    region          = "",                        -- Regin data for this node
    size            = "",                        -- Size info for this node
    address         = "",                        -- IPv[4,6] address or host name of this node
    port            = 0,                         -- Port this node is listening on
 
    status          = "unknown",                 -- Replicating
    details         = "?/?/?",                   -- "SFO1/1GB/25GB disk
    raft            = { 0, 0 },                  -- array 8960, 8960
    load            = { 0, 0 },                  -- some load info
    cpu             = { Used = 0, Total = 0 },   -- some cpu info
    ram             = { Used = 0, Total = 0 },   -- some ram info
    coordinates     = { Lat  = 0, Lng   = 0 },   -- coordinates of the machine
  },
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 22:07:00 GMT
Content-Length: 357
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "message": "OK",
  "node": {
    "address": "64.227.11.116",
    "coordinates": {
      "Lat": 0,
      "Lng": 0
    },
    "cpu": {
      "Total": 0,
      "Used": 0
    },
    "details": "?/?/?",
    "id": 1,
    "image": "i386/1/1MB/100MB",
    "load": [
      0,
      0
    ],
    "name": "Dev1 Server",
    "port": 8860,
    "provider": "DigitalOcean",
    "raft": [
      0,
      0
    ],
    "ram": {
      "Total": 0,
      "Used": 0
    },
    "region": "Rome/Italy",
    "size": "small",
    "status": "unknown",
    "type": "worker"
  },
  "status": 0
}
```