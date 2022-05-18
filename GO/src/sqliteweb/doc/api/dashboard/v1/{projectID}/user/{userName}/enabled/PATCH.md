# API Documentation

ENABLE USER %
DISABLE USER %

## Requests

```sh
## Request GET LUA
curl -X "PATCH" "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/admin/enabled" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMjQ0MzIsImp0aSI6IjEiLCJpYXQiOjE2NTA5OTQ0MzIsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTk0NDMyLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.1WLqJGTuPu-BEJJ3ExNFdIAaEv3iRc3bec4fVMZ9Jzk' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "enabled": false
}'
```

### **POST** - /dashboard/v1/{projectID}/user/{userName}/enabled

### Request object

```json
{
  enabled           = "<boolean value>",      -- optional: boolean value as boolean, int or string (default: true)
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
PATCH /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/admin/enabled HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMjQ0MzIsImp0aSI6IjEiLCJpYXQiOjE2NTA5OTQ0MzIsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTk0NDMyLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.1WLqJGTuPu-BEJJ3ExNFdIAaEv3iRc3bec4fVMZ9Jzk
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 18

{
  "enabled": "true"
}
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Tue, 26 Apr 2022 19:04:11 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```