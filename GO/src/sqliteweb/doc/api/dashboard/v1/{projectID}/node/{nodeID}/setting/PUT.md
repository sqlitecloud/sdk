# API Documentation

Modify setting with key to value

## Requests

```sh
curl -X "PUT" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/10/setting/keyName?x=Hallo%20wie%20gehts&y=1&z=true&x=Second%20line" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "value": "SomeValue2"
}'
```

### **PUT** - /dashboard/v1/{projectID}/node/{nodeID}/setting/{key}

### Request object

```code
{
  value           = "SomeValue",
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
PUT /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/node/10/setting/keyName HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 22

{
  "value": "SomeValue2"
}
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 Success
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 17 Mar 2022 17:52:21 GMT
Content-Length: 43
Connection: close

{
  "message": "OK",
  "status": 200
}
```