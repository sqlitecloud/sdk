# API Documentation

Return list of the backup settings

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwOTM4MzUsImp0aSI6IjEiLCJpYXQiOjE2NTEwNjM4MzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUxMDYzODM1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.6oTRZEBprnPjHoPpxd89RDfHifXn38MQmvureXl2XbY'
```

### **GET** - /dashboard/v1/{projectID}/backups

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

  value             = [ list of settings objects ]; Array with backup settings for each database
}
```

#### settings object:

```json
{ 
  name                      = "",                                   -- database name
  enabled                   = 1,                                    -- backup enabled or disabled
  backup_retention          = "24h",                                -- retention (null if default value)
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/backups/settings HTTP/1.1
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
Date: Wed, 27 Apr 2022 16:37:02 GMT
Content-Length: 69
Connection: close

{
   "message":"OK",
   "status":200,
   "value":[
      {
         "backup_retention":"24h",
         "enabled":1,
         "name":"chinook.sqlite"
      },
      {
         "backup_retention":"168h",
         "enabled":1,
         "name":"db1.sqlite"
      },
      {
         "enabled":0,
         "name":"db3.sqlite"
      },
      {
         "enabled":0,
         "name":"db4enc.sqlite"
      },
      {
         "enabled":0,
         "name":"db5enc.sqlite"
      },
      {
         "enabled":0,
         "name":"dbempty.sqlite"
      }
   ]
}
```