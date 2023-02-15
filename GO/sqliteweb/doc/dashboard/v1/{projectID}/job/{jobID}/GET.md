### Example Request:

```http
GET /dashboard/v1/9905669e-76a3-450f-ae92-0ff5e3537f96/job/b2f11b2c-0447-4b2e-b984-2b730eb6a63f HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3NjQ2MzM3NywibmJmIjoxNjc2NDMzMzc3LCJpYXQiOjE2NzY0MzMzNzd9.l-b2MJ6Ubx8RqTZmEy9eUvHlGKO_xT30xUbb-Z2HEow
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.1.0) GCDHTTPRequest
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Wed, 15 Feb 2023 04:18:09 GMT
Content-Length: 280
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