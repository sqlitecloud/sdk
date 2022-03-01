# API Documentation

LIST DATABASE CONNECTIONS [ID] %

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  status            = 0,                         ; status code: 0 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  connections       = c,                         ; Array with Connection objects
}
```

#### Connection object:

```code
{
  id              = 0,                          ; Internal connection id
  address         = "127.0.0.1",                ; Clients IPv[4/6]address
  username        = "admin",                    ; Login username
  database        = "Dummy",                    ; Database name in use
  connectionDate  = "1970-01-01 00:00:00",      ; Date of connection in SQL format
  lastActivity    = "1970-01-01 00:00:00"       ; Date of last Activity in SQL format
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/Dummy/connections HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU3Mjk5NzAsImp0aSI6IjAiLCJpYXQiOjE2NDU2OTk5NzAsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1Njk5OTcwLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.8izk50ZCk4kQ7Mpf99tj3DuSOuJhPS2cFpAuhlvoGQw
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
Date: Thu, 24 Feb 2022 11:40:16 GMT
Content-Length: 195
Connection: close

{
  "connections": [
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
  "status": 0
}
```

### Previous Response:

```
{
  "ResponseID": 0,
  "Message": "Connections List",
  "Connections": [
      {
          "Id": 2,
          "Address": "192.168.1.23",
          "Username": "admin",
          "Database": "db1",
          "ConnectionDate": "January 1, 1970 00:00:00 UTC",
          "LastActivity": "January 1, 1970 00:00:00 UTC"
      },
      {
          "Id": 4,
          "Address": "192.168.1.23",
          "Username": "admin",
          "Database": "db1",
          "ConnectionDate": "January 1, 1970 00:00:00 UTC",
          "LastActivity": "January 1, 1970 00:00:00 UTC"
      },
      {
          "Id": 7,
          "Address": "192.168.1.23",
          "Username": "admin",
          "Database": "db1",
          "ConnectionDate": "January 1, 1970 00:00:00 UTC",
          "LastActivity": "January 1, 1970 00:00:00 UTC"
      }
  ]
}
```