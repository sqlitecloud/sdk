# API Documentation

DROP DATABASE % IF EXISTS

## Requests

```sh
curl -X "DELETE" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/x" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU' \
     -H 'Content-Type: application/json; charset=utf-8'
```

### **DELETE** - /dashboard/v1/{projectID}/database/{databaseName}

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  message         = "OK",
  status          = 200
}
```

### Example Request:

```
DELETE /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/x HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 2

{}
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 18 Mar 2022 11:38:52 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```