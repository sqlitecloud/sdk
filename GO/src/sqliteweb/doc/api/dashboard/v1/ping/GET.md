# API Documentation

Ping endpoint (for testing)

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/ping" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDgyNDkxMzMsImp0aSI6IjEiLCJpYXQiOjE2NDgyMTkxMzMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ4MjE5MTMzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.9pc9rEdg3iFsDMzgkfR_rFVf0BtA38UVOZNGhdgHTCA' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{}'
```

### **GET** - /dashboard/v1/ping

### Request object

```code
none
```

### Response object(s)

#### root Response:

```code
{
  message         = "PONG",
  status          = 200
}
```

### Example Request:

```
GET /dashboard/v1/ping HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDgyNDkxMzMsImp0aSI6IjEiLCJpYXQiOjE2NDgyMTkxMzMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ4MjE5MTMzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.9pc9rEdg3iFsDMzgkfR_rFVf0BtA38UVOZNGhdgHTCA
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response (user is in the auth database):

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Fri, 25 Mar 2022 14:42:53 GMT
Content-Length: 31
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "message": "PONG",
  "status": 200
}
```