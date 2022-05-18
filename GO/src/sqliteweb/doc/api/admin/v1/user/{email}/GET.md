# API Documentation

Get all customer data

## Requests

```sh
curl "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{}'
```

### **GET** - /admin/v1/user/{email}

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
  value              = {},                       -- User objects
}
```

#### User object:

```code
{
  company         = "SQLiteCloud Inc.",
  company_enabled = 1,
  created         = "2021-11-22 19:01:18",
  email           = "sqlitecloud@synergiezentrum.com",
  enabled         = 1,
  id              = 1,
  first_name      = "Andreas",
  last_name       = "Pfeil",
  password        = "password",
  recoveryRequest = "2021-11-22 19:01:18"
}
```

### Example Request:

```
GET /admin/v1/user/sqlitecloud@synergiezentrum.com HTTP/1.1
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 13:19:59 GMT
Content-Length: 247
Connection: close

{
  "message": "OK",
  "status": 200,
  "user": [
    {
      "company": "SQLiteCloud Inc.",
      "company_enabled": 1,
      "created": "2021-11-22 19:01:18",
      "email": "sqlitecloud@synergiezentrum.com",
      "enabled": 1,
      "id": 1,
      "first_name": "Andreas",
      "last_name": "Pfeil",
      "password": "password",
      "recoveryRequest": "2021-11-22 19:01:18"
    }
  ]
}
```