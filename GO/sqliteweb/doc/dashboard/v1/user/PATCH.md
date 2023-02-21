# API Documentation

Update my user

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/user" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "first_name": "Andrea",
  "last_name": "Donetti",
  "email": "mynewmail@sqlitecloud.io",
  "password": "xxxxxxxxxxx"
}'

```

### **PATCH** - /dashboard/v1/user

### Request object

```json
{
  first_name    = "Andrea",              		; optional
  last_name     = "Donetti",       				; optional
  email  		= "mynewmail@sqlitecloud.io",   ; optional
  password  	= "myNewPassword",         		; optional
}
```

### Response object(s)

#### root Response:

```json
{
  message         = "OK",
  status          = 200
}
```

### Example Request:

```http
PATCH /dashboard/v1/user HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3Njk5MjcxNywibmJmIjoxNjc2OTYyNzE3LCJpYXQiOjE2NzY5NjI3MTd9.mBGnYlDDeUWKo8IOpv-XthoFmRBHmY-wVQI6-q4zXyQ
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 101

{
  "first_name": "Andrea",
  "last_name": "Donetti",
  "email": "mynewmail@sqlitecloud.io",
  "password": "xxxxxxxxxxx"
}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Tue, 21 Feb 2023 14:01:14 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```