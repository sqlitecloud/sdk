# API Documentation

LIST PRIVILEGES / List all privileges for this project

## Requests

```sh
curl "https://localhost:8443https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/privileges" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - https://localhost:8443/dashboard/v1/{projectID}/privileges

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  priveleges        = [ strings,... ],         ; String-Array of privileges 
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/privileges HTTP/1.1
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
Date: Thu, 24 Feb 2022 11:53:17 GMT
Content-Length: 442
Connection: close

{
  "message": "OK",
  "priveleges": [
    "NONE",
    "READ",
    "INSERT",
    "UPDATE",
    "DELETE",
    "READWRITE",
    "PRAGMA",
    "CREATE_TABLE",
    "CREATE_INDEX",
    "CREATE_VIEW",
    "CREATE_TRIGGER",
    "DROP_TABLE",
    "DROP_INDEX",
    "DROP_VIEW",
    "DROP_TRIGGER",
    "ALTER_TABLE",
    "ANALYZE",
    "ATTACH",
    "DETACH",
    "DBADMIN",
    "SUB",
    "PUB",
    "PUBSUB",
    "BACKUP",
    "RESTORE",
    "DOWNLOAD",
    "PLUGIN",
    "PREFERENCES",
    "USERADMIN",
    "CLUSTERADMIN",
    "CLUSTERMONITOR",
    "CREATE_DATABASE",
    "DROP_DATABASE",
    "HOSTADMIN",
    "ADMIN"
  ],
  "status": 0
}
```