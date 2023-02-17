# API Documentation

Return the list of project's active and non archived node-related jobs

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/jobs/nodes" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NTEwMTE4NzUsImp0aSI6IjEiLCJpYXQiOjE2NTA5ODE4NzUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjUwOTgxODc1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.naId5iK5LSm9b52XvQVKytkQmFzTeDjSyamcGYVwWPs'
```

### **PUT** - /dashboard/v1/{projectID}/jobs/nodes

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  "message": "OK",                                       ; "OK" or error message
  "status": 200,                                         ; status code: 200 = no error, error otherwise
  "value": [                                             ; array of job objects
    job_object
  ]
}
```

#### job object:


```json
{
    "name": "test-dev-07",
    "nodeID": 7,
    "uuid": "5e1cf16b-fcb2-4ae6-abda-a6c2e917ca48"
}
```

### Example Request:

```http
GET /dashboard/v1/9905669e-76a3-450f-ae92-0ff5e3537f96/jobs/nodes HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3NjYxNzY3NiwibmJmIjoxNjc2NTg3Njc2LCJpYXQiOjE2NzY1ODc2NzZ9.GnKx0CnIrMkgwG_x0GEpDQMDkBzysqzp8Z988Bv5bnY
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 17 Feb 2023 03:12:02 GMT
Content-Length: 739
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "error": 0,
      "hostname": "flckgxltm.sqlite.cloud",
      "modified": "2023-02-15 03:59:09",
      "name": "Create Node test-dev-03",
      "node_id": 22,
      "node_name": "test-dev-03",
      "progress": 2,
      "status": "Completed",
      "steps": 2,
      "uuid": "87fe0142-1c9c-46a4-ad96-cb9eff034608"
    },
    {
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
    },
    {
      "error": 0,
      "modified": "2023-02-17 03:11:17",
      "name": "Create Node test-dev-07",
      "node_id": 25,
      "node_name": "test-dev-07",
      "progress": 0,
      "status": "Creating droplet",
      "steps": 2,
      "uuid": "5e1cf16b-fcb2-4ae6-abda-a6c2e917ca48"
    }
  ]
}
```