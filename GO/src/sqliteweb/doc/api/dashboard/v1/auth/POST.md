# API Documentation

AUTH

## Requests

```sh
## Request - AUTH
curl -X "POST" "https://web1.sqlitecloud.io:8443/dashboard/v1/auth" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "Login": "email@domain.name",
  "Password": "secret"
}'


```

### **POST** - /dashboard/v1/auth

### Request object

```code
{
  Login           = "email@domain.name",        ; Email Adress of user
  Password        = "secret",                   ; Secret Password for user
}
```

### Response object(s)

#### root Response:

```code
{
  Status           = 0,                         ; status code: 0 = no error, error otherwise
  Message          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTIxNTEsImp0aSI6IjAiLCJpYXQiOjE2NDc1MjIxNTEsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyMTUxLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.hMPSCUP0hLYAA2UTddQgMqCAzOYepE6nXAU-iBspWZs"
}
```

### Example Request:

```
POST /dashboard/v1/auth HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 54

{
  "Login": "admin@sqlitecloud.io",
  "Password": "password"
}
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 17 Mar 2022 13:07:11 GMT
Content-Length: 290
Connection: close

{
  "Status": 0,
  "Message": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI0MzEsImp0aSI6IjAiLCJpYXQiOjE2NDc1MjI0MzEsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyNDMxLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Tp4cCDfapafLeSqZ5q8Cfok-LQaGi7szi686Vp9Zqeg"
}
```