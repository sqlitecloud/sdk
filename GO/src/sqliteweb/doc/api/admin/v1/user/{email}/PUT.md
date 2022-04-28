# API Documentation

Update User data

## Requests

```sh
curl -X "PUT" "https://localhost:8443/admin/v1/user/sqlitecloud1@synergiezentrum.com" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{
  "password": "passw2",
  "enabled": "true",
  "company": "Synergie",
  "name": "Andreas Pfeil"
}'
```

### **PUT** - /admin/v1/user/{email}

### Request object

```code
{
  name             = "customer name",             // mandatory, at least 2 chars
  company          = "Customer company name",     // can be empty
  password         = "new password",              // mandatory, at least 5 chars
  enabled          = 1                            // 0 = disabled, 1 = enabled
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
PUT /admin/v1/user/sqlitecloud1@synergiezentrum.com HTTP/1.1
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 82

{
  "password": "passw2",
  "name": "Andreas Pfeil",
  "enabled": "true",
  "company": "Synergie"
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
Date: Thu, 28 Apr 2022 14:29:10 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```