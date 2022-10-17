# API Documentation

CREATE APIKEY USER <username> NAME <key_name> [RESTRICTION <restriction_type>] [EXPIRATION <expiration_date>]

## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/apikey" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDkxNzU4MTMsImp0aSI6IjEiLCJpYXQiOjE2NDkxNDU4MTMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ5MTQ1ODEzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.A2P2wC9AyNcIFWm4AksF77RQWRVA2sRLTm9l7zy04uY' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "username": "apiuser",
  "name": "key1",
  "expiration": "2022-09-05 12:26:56"
}'
```

### **POST** - /dashboard/v1/{projectID}/apikey/

### Request object

```json
{
  username       = "apiuser"
  name           = "key2",
  expiration     = "2022-09-05 12:26:56",     // optional
  restriction    = 0,                         // optional
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
POST /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/apikey HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDkxNzU4MTMsImp0aSI6IjEiLCJpYXQiOjE2NDkxNDU4MTMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ5MTQ1ODEzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.A2P2wC9AyNcIFWm4AksF77RQWRVA2sRLTm9l7zy04uY
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 36

{
  "username": "apiuser",
  "name": "key2"
}
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Tue, 05 Apr 2022 08:05:15 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```