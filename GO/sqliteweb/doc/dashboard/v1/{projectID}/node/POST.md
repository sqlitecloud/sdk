# API Documentation

Create a new node

## Requests

```sh
curl -X "POST" "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "name": "test-dev",
  "region": "New York 3",
  "hardware": "1VCPU/1GB/25GB",
  "type": "worker",
  "counter": 1
}'
```

### **POST** - /dashboard/v1/{projectID}/node/

### Request object

```json
{
  "name": "test-dev",             // mandatory
  "region": "New York 3",         // mandatory
  "hardware": "1VCPU/1GB/25GB",   // mandatory
  "type": "worker",               // mandatory
  "counter": 1                    // mandatory
}
```

### Response object(s)

#### root Response:

```json
{
  "message": "OK",                                       ; "OK" or error message
  "status": 200,                                         ; status code: 200 = no error, error otherwise
  "value": [                                             ; array of job objects
    {
      "name": "test-dev-07",
      "nodeID": 7,
      "uuid": "5e1cf16b-fcb2-4ae6-abda-a6c2e917ca48"
    }
  ]
}
```

### Example Request:

```http
POST /dashboard/v1/9905669e-76a3-450f-ae92-0ff5e3537f96/node HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3NjYxNzY3NiwibmJmIjoxNjc2NTg3Njc2LCJpYXQiOjE2NzY1ODc2NzZ9.GnKx0CnIrMkgwG_x0GEpDQMDkBzysqzp8Z988Bv5bnY
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 97

{
  "name": "test-dev",
  "region": "New York 3",
  "hardware": "1VCPU/1GB/25GB",
  "type": "worker",
  "counter": 1
}
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 17 Feb 2023 03:11:17 GMT
Content-Length: 119
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "name": "test-dev-07",
      "nodeID": 7,
      "uuid": "5e1cf16b-fcb2-4ae6-abda-a6c2e917ca48"
    }
  ]
}
```