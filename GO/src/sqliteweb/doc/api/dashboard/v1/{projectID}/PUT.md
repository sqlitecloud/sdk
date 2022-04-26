# API Documentation

Update an existing project

## Requests

```sh
## Request POST LUA Duplicate
curl -X "PUT" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "username": "NodeLoginName",
  "password": "NodeLoginPassword",
  "name": "ProjectName",
  "description": "ProjectDescription"
}'
```

### **PUT** - /dashboard/v1/{projectID}

### Request object

```code
{
  name          = "ProjectName",              ; Name of Project
  description   = "ProjectDescription",       ; Description for Project
  username      = "NodeLoginName",            ; Internal name for loggin into the nodes
  password      = "NodeLoginPassword",        ; Internal password for loggin into the nodes
}
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
PUT /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64?x=Hallo%20wie%20gehts&y=1&z=true&x=Second%20line HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 115

{
  "name": "ProjectName",
  "description": "ProjectDescription",
  "username": "NodeLoginName",
  "password": "NodeLoginPassword"
}
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 17 Mar 2022 17:39:27 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```