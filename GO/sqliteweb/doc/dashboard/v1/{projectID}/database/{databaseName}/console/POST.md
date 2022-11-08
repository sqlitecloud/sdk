# API Documentation

Execute a SQLiteCloud command, on the specified cluster/database

## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/dbname/console" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "command": "SELECT * FROM table1"
}'
```

### **POST** - /dashboard/v1/{projectID}/database/{databaseName}/console

### Request object

```json
{
  command           = "",             // SQLiteCloud command to execute
}
```

### Response object(s)

#### root Response:

```json
{
  message         = "OK",
  status          = 200,
  value           = [
    {
      b           = "b1",
      a           = 1
    },
    {
      b           = "b2",
      a           = 2
    }
  ]
}
```

### Example Request:

```http
POST /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/db1.sqlite HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 20

{
  "command": "SELECT * FROM t1"
}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 18 Mar 2022 11:30:08 GMT
Content-Length: 29
Connection: close

{
  "message":"OK",
  "status": 200,
  "value":[
    {
      "HEX(b)":"1BE2E8F457A4B96E25D586804D2B52BFCCD7A0E156D7E985A138954EE02527E1",
      "a":1
    },
    {
      "HEX(b)":"03B8B12C0980419BE8CDCD6B757685F39441D92B83C391C5989E9DC8BA031D30",
      "a":2
    },
    {
      "HEX(b)":"5317D9BE3288F6102429B29933C015750E7DBC38CD08563C9654DAE21CBBFD01",
      "a":3
    },
    {
      "HEX(b)":"A998DFFB12EAE00E777435F9B0E932ACF4F9B798C570E712BEDE5A40FE83D292",
      "a":4
    }
  ]      
}
```