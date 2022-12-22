# API Documentation

Modify database REST API settings

## Requests

```sh
curl -X "PATCH" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/db1.sqlite/settings/rest" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'[
  {
    "name": "albums",
    "GET": true,
    "POST": false,
    "PATCH": true
  },
  {
    "name": "artists",
    "GET": true,
    "POST": false
  }
]'
```

### **PUT** - /dashboard/v1/{projectID}/node/{nodeID}/setting/{key}

### Request object

```json
[ table settings objects]
```

#### Value object:

```json
{
  tableName         = "Test",
  GET               = true,         // optional: false value is used if not specified
  POST              = true,         // optional: false value is used if not specified
  PATCH             = false,        // optional: false value is used if not specified
  DELETE            = false         // optional: false value is used if not specified
}
```

### Response object(s)

#### root Response:

```json
{
  message         = "OK",
  status          = 200
}
```

### Example Request:

```http
PATCH /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/chinook.sqlite/settings/rest HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3MTY3NDUyMywibmJmIjoxNjcxNjQ0NTIzLCJpYXQiOjE2NzE2NDQ1MjN9.E6YZBZCxcZPNJNuGdIThtv82XfVvZH342t4VyQXahIA
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.1.0) GCDHTTPRequest
Content-Length: 109

[
  {
    "tableName": "albums",
    "GET": true,
    "POST": false,
    "PATCH": true
  },
  {
    "tableName": "artists",
    "GET": true,
    "POST": false
  }
]
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Wed, 21 Dec 2022 23:32:25 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```