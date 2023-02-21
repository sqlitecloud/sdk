# API Documentation

Reset the password for the specified user, using a valid token

## Requests

```sh
## Request PATCH LUA Duplicate
curl "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com/recover" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{
  "token": "VODisbCozpy8UEsrEobstQBmsqpVR8F4YLI4L6fY5b7pJSjEJowa4KwEdWcUhT7z",
  "password": "$F4ng2ngl2"
}'
```

### **PATCH** - /admin/v1/user/{email}/recover

### Request object

```json
{
  "token": "VODisbCozpy8UEsrEobstQBmsqpVR8F4YLI4L6fY5b7pJSjEJowa4KwEdWcUhT7z", ; mandatory
  "password": "$F4ng2ngl2"                                                     ; mandatory
}
```

### Response object(s)

#### root Response:

```json
{
  message         = "OK",                     -- "OK" or error message
  status          = 200                       -- status code: 200 = no error, error otherwise
}
```

### Example Request:

```http
PATCH /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64?x=Hallo%20wie%20gehts&y=1&z=true&x=Second%20line HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 115

{
  "token": "VODisbCozpy8UEsrEobstQBmsqpVR8F4YLI4L6fY5b7pJSjEJowa4KwEdWcUhT7z",
  "password": "$F4ng2ngl2"
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