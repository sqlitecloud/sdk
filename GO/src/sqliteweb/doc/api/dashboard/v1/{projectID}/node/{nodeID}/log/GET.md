# API Documentation

Filter log

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/log?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs'
```

### **GET** - https://localhost:8443/dashboard/v1/{projectID}/node/{nodeID}/log

### Query parameters for row based queries (LIST % ROWS FROM LOG [LEVEL %] [TYPE %])

```json
level = 4                                       -- optional, integer between 0..5 (default = 4)
type  = 4                                       -- optional, integer between 1..8 (default = 4)
```

If it is not a row based query (rows parameter is missing), the call to this endpoint is handled as a date based query.

### Query parameters for date based queries (LIST LOG FROM % TO % [LEVEL %] [TYPE %])

```json
from  = '2022-04-02 17:53:04'                   -- optional, date string in SQL format (default = now minus one day)
to    = '2022-04-26 18:53:04'                   -- optional, date string in SQL format (default = now)
level = 4                                       -- optional, integer between 0..5 (default = 4)
type  = 4                                       -- optional, integer between 1..8 (default = 4)
```

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

  value             = nil,                       -- Array with log file
}
```

#### Value object:

```json
{
 address     = "5.100.32.221",
 date        = "2022-04-26 16:58:59",
 description = "LIST LOG FROM '2022-04-02 17:53:04' TO '2022-04-26 18:53:04' LEVEL 4 TYPE 4 ORDER DESC;",
 level       = 4,
 type        = 4,
 username    = "admin",
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/log?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04 HTTP/1.1
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
Date: Tue, 26 Apr 2022 16:59:00 GMT
Connection: close
Transfer-Encoding: chunked

{
  "value": [
    {
      "address": "5.100.32.221",
      "date": "2022-04-26 16:58:59",
      "description": "LIST LOG FROM '2022-04-02 17:53:04' TO '2022-04-26 18:53:04' LEVEL 4 TYPE 4 ORDER DESC;",
      "level": 4,
      "type": 4,
      "username": "admin"
    },
    ...
    ],
  "message": "OK",
  "status": 200
}  
```