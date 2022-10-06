# API Documentation

Create new User setting with key

## Requests

```sh
curl -X "POST" "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com/setting/key" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{
  "value": "abc"
}'
```

### **POST** - /admin/v1/user/{email}/setting/{key}

### Request object

```code
{
  value             = "string value",           -- set the new key to "string value"
}
```

### Response object(s)

#### root Response:

```code
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
}
```

### Example Request:

```
POST /admin/v1/user/sqlitecloud@synergiezentrum.com/setting/key HTTP/1.1
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 15

{
  "value": "abc"
}
```

### Example Response :

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 14:18:26 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```