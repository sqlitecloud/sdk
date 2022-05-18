# API Documentation

Displays name and description for a project

## Requests

```sh
## Request re-AUTH Duplicate (2)
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs'
```

### **GET** - /dashboard/v1/{projectID}

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status           = 200,                                     -- status code: 200 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message

  value            = nil                                      -- Single project object
}
```

#### Value object:

```json
{
  id               = "00000000-0000-0000-0000-000000000000",  -- UUID of the project
 
  name             = "",                                      -- Project name
  description      = ""                                       -- Project description
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Tue, 26 Apr 2022 18:07:33 GMT
Content-Length: 149
Connection: close

{
  "message": "OK",
  "value": {
    "description": "Just for internal testing",
    "id": "fbf94289-64b0-4fc6-9c20-84083f82ee64",
    "name": "Test Project"
  },
  "status": 200
}
```