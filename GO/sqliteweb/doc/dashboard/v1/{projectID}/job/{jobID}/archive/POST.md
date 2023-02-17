# API Documentation

Archive the job

## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/job/87fe0142-1c9c-46a4-ad96-cb9eff034608/archive" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDkxNzU4MTMsImp0aSI6IjEiLCJpYXQiOjE2NDkxNDU4MTMsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ5MTQ1ODEzLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.A2P2wC9AyNcIFWm4AksF77RQWRVA2sRLTm9l7zy04uY'
```

### **POST** - /dashboard/v1/{projectID}/job/{jobID}/archive

### Request object

```code
none
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
POST /dashboard/v1/9905669e-76a3-450f-ae92-0ff5e3537f96/job/87fe0142-1c9c-46a4-ad96-cb9eff034608/archive HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3NjYxNzY3NiwibmJmIjoxNjc2NTg3Njc2LCJpYXQiOjE2NzY1ODc2NzZ9.GnKx0CnIrMkgwG_x0GEpDQMDkBzysqzp8Z988Bv5bnY
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 0
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 17 Feb 2023 03:31:26 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": {
    "error": 0,
    "hostname": "upokxe2ap.sqlite.cloud",
    "modified": "2023-02-15 04:16:11",
    "name": "Create Node test-dev-04",
    "node_id": 23,
    "node_name": "test-dev-04",
    "progress": 2,
    "status": "Completed",
    "steps": 2,
    "uuid": "b2f11b2c-0447-4b2e-b984-2b730eb6a63f"
  }
}
```