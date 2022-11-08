# API Documentation

LIST APIKEYS 

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/apikeys" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'
```

### **GET** - /dashboard/v1/{projectID}/apikeys

### Query parameters

### Request object

```
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                          ; status code: 200 = no error, error otherwise
  message           = "OK",                         ; "OK" or error message

  value             = {                             ; Map with user objects, only users that have at least one apikey
                        "user1": [apikeysobj,...]   ; list of apikeys objects
                        "user2": [apikeysobj,...]
                      },                        
}
```

#### apikeysobj:

```json
{
  key           = "3456789",                         
  name          = "test1",   
  expiration    = "2022-09-05 12:26:56" 
  restriction   = 0                    
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/apikeys HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU3Mjk5NzAsImp0aSI6IjAiLCJpYXQiOjE2NDU2OTk5NzAsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1Njk5OTcwLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.8izk50ZCk4kQ7Mpf99tj3DuSOuJhPS2cFpAuhlvoGQw
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 24 Feb 2022 12:41:59 GMT
Content-Length: 109
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": {
    "apiuser": {
      "key": "rCd01HtTjqIefqWuTg780xVPWQvDpi6aSzwbXw5AAAA",
      "name": "key1",
      "expiration_date": "2022-09-05 12:26:56",
      restriction: 0
    },  
  }
}
```