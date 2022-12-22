LIST TABLES with REST API settings

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/db1.sqlite/settings/rest" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY'
```

### **GET** - /dashboard/v1/{projectID}/database/{databaseName}/settings/rest

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = [ table settings objects],     -- Array of table settings objects
}
```

#### Value object:

```json
{
  tableName         = "Test",
  GET               = true,
  POST              = true,
  PATCH             = false,
  DELETE            = false
}
```

### Example Request:

```http
GET /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/chinook.sqlite/settings/rest HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3MTY3NDUyMywibmJmIjoxNjcxNjQ0NTIzLCJpYXQiOjE2NzE2NDQ1MjN9.E6YZBZCxcZPNJNuGdIThtv82XfVvZH342t4VyQXahIA
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.1.0) GCDHTTPRequest
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Wed, 21 Dec 2022 23:27:26 GMT
Content-Length: 915
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "DELETE": false,
      "GET": true,
      "PATCH": true,
      "POST": false,
      "tableName": "albums"
    },
    {
      "DELETE": false,
      "GET": true,
      "PATCH": false,
      "POST": false,
      "tableName": "artists"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "playlists"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "customers"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "employees"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "genres"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "tracks"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "media_types"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "invoices"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "playlist_track"
    },
    {
      "DELETE": false,
      "GET": false,
      "PATCH": false,
      "POST": false,
      "tableName": "invoice_items"
    }
  ]
}
```