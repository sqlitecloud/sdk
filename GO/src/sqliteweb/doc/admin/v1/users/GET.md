# API Documentation

Get a List of Customers

## Requests

```sh
curl "https://localhost:8443/admin/v1/users/" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{}'
```

### **GET** - /admin/v1/users

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

  users             = [{}],                      -- Array with User objects
}
```

#### User object:

```code
{
  company     = "SQLiteCloud Inc.",                 
  email       = "sqlitecloud@synergiezentrum.com",
  id          = 1,
  first_name  = "Andreas",
  last_name   = "Pfeil"
}
```

### Example Request:

```
GET /admin/v1/users/ HTTP/1.1
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
Date: Thu, 28 Apr 2022 13:17:04 GMT
Content-Length: 414
Connection: close

{
  "message": "OK",
  "status": 200,
  "users": [
    {
      "company": "SQLiteCloud Inc.",
      "email": "andrea@sqlitecloud.io",
      "id": 1,
      "first_name": "Andrea",
      "last_name": "Donetti"
    },
    {
      "company": "SQLiteCloud Inc.",
      "email": "marco@sqlitecloud.io",
      "id": 3,
      "first_name": "Marco",
      "last_name": "Bambini"
    },
    {
      "company": "Synergie",
      "email": "my2.address@domain.com",
      "id": 23,
      "name": "Andreas Pfeil"
    }
  ]
}
```