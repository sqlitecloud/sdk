# API Documentation

Return list of available backups

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backup/db1.sqlite" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY'
```

### **GET** - /dashboard/v1/{projectID}/backup/{databaseName}

### Request object

```
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                        ; status code: 200 = no error, error otherwise
  message           = "OK",                       ; "OK" or error message

  value             = [ {backup info objcts } ] ; Array with backup info objects
}
```

#### Value Object:

```json
{
  replica           = "file",
  size              = 462,
  time_stamp        = "2022-04-27T02:45:02Z",
  type              = "snapshot"
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backup/db1.sqlite HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY
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
Date: Wed, 27 Apr 2022 16:42:46 GMT
Content-Length: 360
Connection: close

{
  "message": "OK",
  "value": [
    {
      "replica": "file",
      "size": 119,
      "time_stamp": "2022-04-26T15:05:40Z",
      "type": "wal"
    },
    {
      "replica": "file",
      "size": 462,
      "time_stamp": "2022-04-27T02:45:02Z",
      "type": "snapshot"
    },
    {
      "replica": "file",
      "size": 504,
      "time_stamp": "2022-04-27T09:32:22Z",
      "type": "wal"
    },
    {
      "replica": "file",
      "size": 119,
      "time_stamp": "2022-04-27T09:32:22Z",
      "type": "wal"
    }
  ],
  "status": 200
}
```