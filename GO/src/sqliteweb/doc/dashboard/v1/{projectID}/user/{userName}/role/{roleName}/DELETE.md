# API Documentation

REVOKE ROLE % USER % [DATABASE %] [TABLE %] 

## Requests

```sh
curl -X "DELETE" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/newUser/role/rolename" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDkxNzU4MTMsImp0aSI6IjEiLCJpYXQiOjE2NDkxNDU4MTMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ5MTQ1ODEzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.A2P2wC9AyNcIFWm4AksF77RQWRVA2sRLTm9l7zy04uY' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
       "database": "databaseName"
     }'
```

### **DELETE** - /dashboard/v1/{projectID}/user/{userName}/role/{roleName}

### Request object

```json
{
  database           = "databaseName",    // optional
  table              = "tableName",       // optional
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
POST /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/user/newUser/newRole HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDkxNzU4MTMsImp0aSI6IjEiLCJpYXQiOjE2NDkxNDU4MTMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ5MTQ1ODEzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.A2P2wC9AyNcIFWm4AksF77RQWRVA2sRLTm9l7zy04uY
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 26

{
  "database": "*",
  "table": "*"
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
Date: Tue, 05 Apr 2022 08:05:15 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```