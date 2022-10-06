# API Documentation

LIST DATABASE CONNECTIONS [ID] %

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'
```

### **GET** - /dashboard/v1/{projectID}/database/{databaseName}/connections

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

  value             = [ connection info objects] ; Array with Connection objects
}
```

#### Value object:

```json
{
  last_activity      = "2022-04-28 07:35:12",    ; Date and Time of last activity
  address            = "143.198.231.152",        ; Client IP address
  connection_date    = "2022-04-28 07:34:37",    ; Date and Time of connection creation
  id                 = 454,                      ; Conneciton ID
  username           = "admin"                   ; Username for this connection
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU3Mjk5NzAsImp0aSI6IjAiLCJpYXQiOjE2NDU2OTk5NzAsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1Njk5OTcwLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.8izk50ZCk4kQ7Mpf99tj3DuSOuJhPS2cFpAuhlvoGQw
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
Date: Thu, 24 Feb 2022 11:40:16 GMT
Content-Length: 195
Connection: close

{
  "value": [
    {
      "address": "5.100.32.221",
      "connectionDate": "2022-02-24 11:24:29",
      "database": "Dummy",
      "id": 21294,
      "lastActivity": "2022-02-24 11:40:13",
      "username": "admin"
    }
  ],
  "message": "OK",
  "status": 200
}
```