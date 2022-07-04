# API Documentation

List all available regions

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/regions" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/regions

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status           = 200,                                     -- status code: 200 = no error, error otherwise
  message          = "OK",                                    -- "OK" or error message

  value            = nil                                      -- Array with strings (regions)
}
```

### Example Request:

```http
GET /dashboard/v1/regions HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 18:31:10 GMT
Content-Length: 150
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    "NYC1/US",
    "NYC2/US",
    "NYC3/US",
    "SFO1/US",
    "SFO2/US",
    "AMS2/NL",
    "AMS3/NL"
  ]
}
```