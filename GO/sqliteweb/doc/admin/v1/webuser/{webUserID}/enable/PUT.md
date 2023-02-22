# API Documentation

Enable webUser

## Requests

```sh
curl -X "PUT" "https://localhost:8443/admin/v1/webuser/3/enable" \
     -u 'admin:password'
```

### **PUT** - /admin/v1/webuser/{webUserID}

### Request object

```code
none
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
PUT /admin/v1/webuser/3/enable HTTP/1.1
Authorization: Basic YWRtaW46cGF=
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 0

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