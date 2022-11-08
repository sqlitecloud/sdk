# API Documentation

List available plugins

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/plugins" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'
```

### **GET** - /dashboard/v1/{projectID}/plugins

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {},                        -- List with Plugins objects
}
```

#### Plugin object :
  
```json
{
    enabled: 1,
    name: "crypto",
    type: "SQLite",
    version: "",
    copyright: "",
    description: ""
},
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/plugins HTTP/1.1
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
Date: Tue, 22 Feb 2022 22:22:42 GMT
Content-Length: 83
Content-Type: text/plain; charset=utf-8
Connection: close
{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "enabled": 1,
      "name": "crypto",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "fileio",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "fuzzy",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "ipaddr",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "json1",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "math",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "re",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "stats",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "text",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "unicode",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "uuid",
      "type": "SQLite"
    },
    {
      "enabled": 1,
      "name": "vsv",
      "type": "SQLite"
    }
  ]
}
```