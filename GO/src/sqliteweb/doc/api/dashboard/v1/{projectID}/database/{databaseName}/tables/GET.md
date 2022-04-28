# API Documentation

LIST TABLES

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/db1.sqlite/tables" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY'
```

### **GET** - /dashboard/v1/{projectID}/database/{databaseName}/tables

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  tables            = [ table info objects],     -- Array of table info objects
}
```

#### Table info object:

```code
{
  columns           = 2,
  name              = "Test",
  schema            = "main",
  strict            = 0,
  type              = "table",
  wr                = 0
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/db1.sqlite/tables?from=2022-04-26%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY
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
Date: Wed, 27 Apr 2022 16:29:55 GMT
Content-Length: 117
Connection: close

{
  "message": "OK",
  "status": 200,
  "tables": [
    {
      "columns": 2,
      "name": "Test",
      "schema": "main",
      "strict": 0,
      "type": "table",
      "wr": 0
    }
  ]
}
```