# API Documentation

Modify the node info

## Requests

```sh
curl -X "PUT" "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "region": "Rome/Italy",
  "provider": "DigitalOcean",
  "address": "64.227.11.116",
  "port": "9960",
  "image": "i386/1/1MB/100MB",
  "size": "small",
  "type": "worker",
  "name": "Dev1 Server"
}'
```

### **PUT** - /dashboard/v1/{projectID}/node/{nodeID}

### Request object

```json
{
  name      = "Dev1 Server",              // mandatory
  type      = "worker",                   // mandatory
  provider  = "DigitalOcean",             // mandatory
  image     = "i386/1/1MB/100MB",         // mandatory
  region    = "Rome/Italy",               // mandatory
  size      = "small",                    // mandatory
  address   = "64.227.11.116",            // mandatory
  port      = 9960                        // mandatory
}
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message
}
```

### Example Request:

```http
PUT /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/ HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 168

{
  "name": "Dev1 Server",
  "type": "worker",
  "provider": "DigitalOcean",
  "image": "i386/1/1MB/100MB",
  "region": "Rome/Italy",
  "size": "small",
  "address": "64.227.11.116",
  "port": "9960"
}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 13:01:11 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```