# API Documentation

Filter log

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/log?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs'
```

### **GET** - https://localhost:8443/dashboard/v1/{projectID}/node/{nodeID}/log



### Query parameters

```json
from   = '2022-04-02 17:53:04'                   -- optional, date string in SQL format (default = unix epoc)
to     = '2022-04-26 18:53:04'                   -- optional, date string in SQL format (default = now)
level  = 4                                       -- optional, integer between 0..5 (default = null -> not filtered)
type   = 4                                       -- optional, integer between 1..8 (default = null -> not filtered)
limit  = 100                                     -- optional, integer (default 100)
cursor = 1234                                    -- optional, integer (default nil), use the next_cursor value from previous responses) 
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

  value             = {
    count           = nil,          -- Number of logs for the current filters, only returned if the CURSOR arg is empty
    next_cursor     = nil,          -- Value to be used in the next request to get the next page
    logs            = {},           -- Array of logs
  }
}
```

#### log object:

```json
{
 address        = "5.100.32.221",
 date           = "2022-04-26 16:58:59",
 description    = "LIST LOG FROM '2022-04-02 17:53:04' TO '2022-04-26 18:53:04' LEVEL 4 TYPE 4;",
 level          = 4,
 type           = 4,
 username       = "admin",
 database       = "db1.sqlite"
 connection_id  = 5
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/1/log?level=4&from=2022-04-02%2017%3A53%3A04&to=2022-04-26%2018%3A53%3A04&limit=100 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with, X-SQLiteCloud-Api-Key
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 13 Jan 2023 19:44:52 GMT
Content-Length: 949
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": {
    "count": 175530,
    "logs": [
      {
        "date": "2022-10-14 21:45:39",
        "description": "database_execute_v2 error executing PRAGMA application_id;  (26 26 file is not a database) changes: 0",
        "level": 2,
        "type": 1
      },
      {
        "date": "2022-10-14 21:45:39",
        "description": "/home/andrea/data/9960/databases/wrongdb9.sqlite is not a valid sqlite3 database (file is not a database).",
        "level": 2,
        "type": 1
      },
      {
        "date": "2022-10-14 21:45:39",
        "description": "database_execute_v2 error executing PRAGMA application_id;  (26 26 file is not a database) changes: 0",
        "level": 2,
        "type": 1
      },
      {
        "date": "2022-10-14 21:45:39",
        "description": "database_execute_v2 error executing PRAGMA application_id;  (11 11 database disk image is malformed) changes: 0",
        "level": 2,
        "type": 1
      },
      {
        "date": "2022-10-14 21:45:39",
        "description": "/home/andrea/data/9960/databases/wrongdb5.sqlite is not a valid sqlite3 database (database disk image is malformed).",
        "level": 2,
        "type": 1
      }
    ],
    "next_cursor": 4
  }
}
```