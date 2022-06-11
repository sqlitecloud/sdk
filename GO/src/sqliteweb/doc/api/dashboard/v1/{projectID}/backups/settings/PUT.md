# API Documentation

Update an existing project

## Requests

```sh
curl -X "PUT" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
        "values": [
            {
                "name": "db1.sqlite",
                "enabled": "1",
                "retention": "12h",
                "snapshot_interval": "6h"
            },
            {
                "name": "chinook.sqlite",
                "enabled": "1",
                "retention": "24h"
            },
            {
               "name": "db3.sqlite",
                "enabled": "0"
            }
        ]
    }'
```

### **PUT** - /dashboard/v1/{projectID}

### Request object

```json
{
  value           = [list of settings objects]            
}
```

#### settings object:

```json
{ 
  name                      = "",                                   -- database name
  enabled                   = 1,                                    -- backup enabled or disabled
  backup_retention          = "12h",                                -- retention (null if default value)
  backup_snapshot_interval  = "6h"                                  -- snapshot interval, must be equal or smaller than  backup_retention (null if default value, by default it is equal to backup_retention )
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
PUT /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 115

{
  "values": [
    {
      "name": "db1.sqlite",
      "enabled": "1",
      "retention": "12h",
      "snapshot_interval": "6h"
    },
    {
      "name": "chinook.sqlite",
      "enabled": "1",
      "retention": "24h"
    },
    {
      "name": "db3.sqlite",
      "enabled": "0"
    }
  ]
}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 17 Mar 2022 17:39:27 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```