# API Documentation

Reset the password for the specified user, using a valid token

## Requests

```sh
## Request PATCH LUA Duplicate
curl "https://localhost:8443/admin/v1/user/reset" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{
  "token": "VODisbCozpy8UEsrEobstQBmsqpVR8F4YLI4L6fY5b7pJSjEJowa4KwEdWcUhT7z",
  "password": "$F4ng2ngl2"
}'
```

### **PATCH** - /admin/v1/user/reset

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
PATCH /admin/v1/user/reset HTTP/1.1
Authorization: Basic YWRtaW46cGxxxxxxxxQ=
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 100

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