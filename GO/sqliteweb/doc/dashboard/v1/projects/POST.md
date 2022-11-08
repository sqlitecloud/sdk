# API Documentation

Create new setting for key for logged in user

## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/dashboard/v1/projects" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "admin_username": "NodeLoginName",
  "admin_password": "NodeLoginPassword",
  "name": "ProjectName",
  "description": "ProjectDescription"
}'

```

### **POST** - /dashboard/v1/projects

### Request object

```json
{
  name              = "ProjectName",              ; Name of Project
  description       = "ProjectDescription",       ; Description for Project
  admin_username    = "NodeLoginName",            ; Internal name for loggin into the nodes
  admin_password    = "NodeLoginPassword",        ; Internal password for loggin into the nodes
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
POST /dashboard/v1/projects HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 115

{
  "name": "ProjectName",
  "description": "ProjectDescription",
  "admin_username": "NodeLoginName",
  "admin_password": "NodeLoginPassword"
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
Date: Thu, 17 Mar 2022 13:51:39 GMT
Content-Length: 85
Connection: close

{
  "message": "OK",
  "project": [
    {
      "id": "75bc0c5f-fc53-458e-9fbc-e617ab5843aa"
    }
  ],
  "status": 0
}
```