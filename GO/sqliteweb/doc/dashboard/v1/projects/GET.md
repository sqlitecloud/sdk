# API Documentation

List all userid projects

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/projects" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/projects

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

  value            = nil                                      -- Array with project objects
}
```

#### Value object:

```json
{
  id               = "00000000-0000-0000-0000-000000000000",  -- UUID of the project
 
  name             = "",                                      -- Project name
  description      = "",                                      -- Project description
  admin_username   = "admin"                                  -- Internal name for loggin into the nodes
}
```

### Example Request:

```http
GET /dashboard/v1/projects HTTP/1.1
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
  "value": [
    {
      "description": "Demo Project for internal testing",
      "id": "fbf94289-64b0-4fc6-9c20-84083f82ee63",
      "name": "Test"
    }
  ],
  "status": 200
}
```