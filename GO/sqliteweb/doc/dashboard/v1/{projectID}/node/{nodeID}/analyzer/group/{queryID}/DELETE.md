# API Documentation

Delete the query analizer's records for the group (database, normalized_sql) specified by the queryID parameter 

## Requests

```sh
curl -X "DELETE" "https://web1.sqlitecloud.io:8443/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/group/10" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc1NTI5OTUsImp0aSI6IjEiLCJpYXQiOjE2NDc1MjI5OTUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTIyOTk1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.SaOn2-XJbf6_irYDvhTGEkDHNHJobiNeEO7CPQVHUi8' \
     -H 'Content-Type: application/json; charset=utf-8'
```

### **DELETE** - /dashboard/v1/dashboard/v1/{projectID}/node/{nodeID}/analyzer/group/{queryID}

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       ; status code: 200 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message
}
```

### Example Request:

```http
DELETE /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/group/10 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3MDAyNzgxOSwibmJmIjoxNjY5OTk3ODE5LCJpYXQiOjE2Njk5OTc4MTl9.MwdQlyGP8YAvoEJ2EayJR7vrD3D0KCxNqiZY7fyzQhw
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.0.0) GCDHTTPRequest
Content-Length: 2

{}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 02 Dec 2022 07:59:45 GMT
Content-Length: 29
Connection: close

{
  "message":"OK",
  "status":200
}
```
