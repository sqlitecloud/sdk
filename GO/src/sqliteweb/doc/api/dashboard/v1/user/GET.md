# API Documentation

Get all data and settings for logged in user

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/user" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/user

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status           = 200,                       ; status code: 200 = no error, error otherwise
  message          = "OK",                      ; "OK" or error message
	value						 = {
		id               = tonumber( userid ),        ; UserID, 0 = static user defined in .ini file
  	enabled          = false,                     ; Whether this user account is enabled or disabled
  	first_name       = "",                        ; First name
  	last_name				 = "",												; Last name
  	company          = "",                        ; User company
  	email            = "",                        ; User email - also used as login
  	creation_date     = "1970-01-01 00:00:00",    ; Date and time when this user account was created
  	settings         = nil,                       ; Array with key/value pairs
	}
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

```http
POST /dashboard/v1/user HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 0
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 17:41:08 GMT
Content-Length: 282
Content-Type: text/plain; charset=utf-8
Connection: close

{
	"status": 200,
	"message": "OK",
	"value": {
		"company": "SQLiteCloud Inc.",
		"creation_date": "2021-11-22 19:01:18",
		"email": "my.address@domain.com",
		"enabled": true,
		"id": 1,
		"first_name": "Andreas",
  	"last_name": "Pfeil",
  	"settings": [
    	{
      	"key": "role",
      	"value": "team member"
    	}
  	]
	}
}
```

### Example Response (no auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 17:50:16 GMT
Content-Length: 231
Content-Type: text/plain; charset=utf-8
Connection: close

{
	"status": 200,
	"message": "OK",
	"value": {
		"company": "SQLiteCloud Inc.",
		"creation_date": "2021-11-22 19:01:18",
		"email": "my.address@domain.com",
		"enabled": true,
		"id": 1,
		"first_name": "Andreas",
  	"last_name": "Pfeil",
  	"settings": [
    	{
      	"key": "role",
      	"value": "team member"
    	}
  	]
	}
}
```