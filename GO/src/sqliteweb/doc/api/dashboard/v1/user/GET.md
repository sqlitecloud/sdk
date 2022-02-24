# API Documentation

Get all data and settings for logged in user

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/user" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/user

### Request object

```json
none
```

### Response object(s)

#### root Response:

```json
{
  status           = 0,                         ; status code: 0 = no error, error otherwise
  message          = "OK",                      ; "OK" or error message

  id               = tonumber( userid ),        ; UserID, 0 = static user defined in .ini file
  enabled          = false,                     ; Whether this user account is enabled or disabled
  name             = "",                        ; User name
  company          = "",                        ; User company
  email            = "",                        ; User email - also used as login
  password         = "*******",                 ; User password - this fiels is always 7 stars
  creationDate     = "1970-01-01 00:00:00",     ; Date and time when this user account was created
  lastRecoveryTime = "1970-01-01 00:00:00",     ; Last date and time when this user has tried to recover his password

  settings         = nil,                       ; Array with key/value pairs
}
```

#### Settings (key/value pair):

```json
{
  key   = "",                                   ; Key
  value = ""                                    ; Value
}
```

### Example Request:

```
POST /dashboard/v1/user HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 0
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 17:41:08 GMT
Content-Length: 282
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "company": "SQLiteCloud Inc.",
  "creationDate": "2021-11-22 19:01:18",
  "email": "my.address@domain.com",
  "enabled": true,
  "id": 1,
  "lastRecoveryTime": "2021-11-22 19:01:18",
  "message": "OK",
  "name": "Andreas Pfeil",
  "password": "*******",
  "settings": [
    {
      "key": "role",
      "value": "team member"
    }
  ],
  "status": 0
}
```

### Example Response (no auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 17:50:16 GMT
Content-Length: 231
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "company": "SQLiteCloud Inc.",
  "creationDate": "2022-02-22 17:47:39",
  "email": "admin@sqlitecloud.io",
  "enabled": true,
  "id": 0,
  "lastRecoveryTime": "2022-02-22 17:47:39",
  "message": "OK",
  "name": "Marco Bambini",
  "password": "*******",
  "status": 0
}
```