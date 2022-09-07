# API Documentation

List all key/value paires for the given user

## Requests

```sh
curl "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com/settings" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{}'
```

### **GET** - /admin/v1/user/{email}/settings

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

  settings          = [{}],                      -- Array with Setting objects (key/value pairs) 
}
```

#### Setting objects (key/value pairs):

```code
{
  key               = "key",
  value             = "value",
}
```

### Example Request:

```
GET /admin/v1/user/sqlitecloud@synergiezentrum.com/settings HTTP/1.1
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
Date: Thu, 28 Apr 2022 14:12:04 GMT
Content-Length: 218
Connection: close

{
  "message": "OK",
  "settings": [
    {
      "key": "key",
      "value": "yyyy"
    },
    {
      "key": "role",
      "value": "Team Member"
    },
    {
      "key": "testKey",
      "value": "bla"
    },
    {
      "key": "testKey2",
      "value": "TestValue"
    },
    {
      "key": "testkey3",
      "value": "SomeValue"
    }
  ],
  "status": 200
}
```