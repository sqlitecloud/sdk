# API Documentation

List  connections to the specified node

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/connections" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c'
```

### **GET** - /dashboard/v1/{projectID}/node/{nodeID}/connections

### Request object

```
none
```

### Response object(s)

#### root Response:

```code
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  connections       = [{}],                      ; Array with connection objects
}
```

#### Connection object:

```json
{
  activityDateTime   = "2022-04-28 07:35:12",    ; Date and Time of last activity
  address            = "143.198.231.152",        ; Client IP address
  connectionDateTime = "2022-04-28 07:34:37",    ; Date and Time of connection creation
  id                 = 454,                      ; Conneciton ID
  username           = "admin"                   ; Username for this connection
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/connections HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTExNjA3MDksImp0aSI6IjEiLCJpYXQiOjE2NTExMzA3MDksImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMTMwNzA5LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.W7HGTl0uKcDLcdsM0wM6Jw-65Reu57WVRVIai9VAw1c
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 07:35:12 GMT
Content-Length: 187
Connection: close

{
  "connections": [
    {
      "activityDateTime": "2022-04-28 07:35:12",
      "address": "143.198.231.152",
      "connectionDateTime": "2022-04-28 07:34:37",
      "id": 454,
      "username": "admin"
    }
  ],
  "message": "OK",
  "status": 200
}
```