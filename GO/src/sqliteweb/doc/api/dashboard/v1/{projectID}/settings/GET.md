# API Documentation

List all project settings

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/settings" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'
```

### **GET** - /dashboard/v1/{projectID}/settings

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

  value             = {},                        -- Array with Setting object (key value pairs)
}
```

#### Setting object (key/value pair):

```json
{
    key:"autocheckpoint",
    value:"1000",
    readonly:0,
    default_value:"1000",
    description:"This is a default description about a setting key." 
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64settings HTTP/1.1
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

         "description":"This is a default description about a setting key.",
         "key":"tcpkeepalive",
         "readonly":0,
         "value":"300"
      },
      {
         "default_value":"10",
         "description":"This is a default description about a setting key.",
         "key":"tcpkeepalive_count",
         "readonly":0,
         "value":"10"
      },
      {
         "default_value":"0",
         "description":"This is a default description about a setting key.",
         "key":"tls_verify_client",
         "readonly":0,
         "value":"0"
      },
      {
         "default_value":"0",
         "description":"This is a default description about a setting key.",
         "key":"use_concurrent_transactions",
         "readonly":0,
         "value":"0"
      },
      {
         "default_value":"60",
         "description":"This is a default description about a setting key.",
         "key":"zombie_timeout",
         "readonly":0,
         "value":"60"
      }
   ]
}
```