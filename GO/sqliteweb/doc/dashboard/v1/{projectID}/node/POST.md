# API Documentation

Create a new node

## Requests

```sh
curl -X "POST" "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "name": "Dev1 Server"
  "hardware": "1VCPU/1GB/25GB",
  "region": "NYC3/US",
  "type": "worker",
  "counter": 1
}'
```

### **PUT** - /dashboard/v1/{projectID}/node/{nodeID}

### Request object

```json
{
  name      = "Dev1 Server",              // mandatory
  hardware  = "1VCPU/1GB/25GB",           // mandatory
  region    = "NYC3/US",                  // mandatory
  type      = "worker",                   // mandatory
  counter   = 1                           // mandatory
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
POST /dashboard/v1/9905669e-76a3-450f-ae92-0ff5e3537f96/node HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3NjQ2MzM3NywibmJmIjoxNjc2NDMzMzc3LCJpYXQiOjE2NzY0MzMzNzd9.l-b2MJ6Ubx8RqTZmEy9eUvHlGKO_xT30xUbb-Z2HEow
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.1.0) GCDHTTPRequest
Content-Length: 88

{
  "name": "test-dev",
  "region": "sfo3",
  "hardware": "s-1vcpu-1gb",
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
Date: Wed, 15 Feb 2023 04:13:34 GMT
Content-Length: 119
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "name": "test-dev-04",
      "nodeID": 4,
      "uuid": "b2f11b2c-0447-4b2e-b984-2b730eb6a63f"
    }
  ]
}
```