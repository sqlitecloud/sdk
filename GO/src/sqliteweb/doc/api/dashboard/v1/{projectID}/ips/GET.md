# API Documentation

LIST ALLOWED IP [ROLE %] [USER %]

## Requests

```sh
curl "/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ips?user=admin" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ips?user=admin

### Query parameters

The user MUST give only one of the following two query parameteres:

role = {roleName}
user = {userName}

If none is provided, an error is raised.
If both are provided, only the user parameter is taken into account.

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  ips               = ips.Rows,                  ; Array with allowded IP's for this role or user
}
```

#### Database object:

```code
{
  address = "127.0.0.1",                         ; IPv[4/6]
  name    = "name",                              ; Name
  type    = "type";                              ; Type
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ips?user=admin HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU3Mjk5NzAsImp0aSI6IjAiLCJpYXQiOjE2NDU2OTk5NzAsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1Njk5OTcwLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.8izk50ZCk4kQ7Mpf99tj3DuSOuJhPS2cFpAuhlvoGQw
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest

```

### Example Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 24 Feb 2022 12:18:11 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```