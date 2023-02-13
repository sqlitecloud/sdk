# API Documentation

Get a JSON with the list of queries recorded by the query analyzer for a specific database and sql (normalized_sql) specified by one of the queryID of that group

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/group/5" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/noode/6/analyzer/group/{queryID}

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

  value             = [QueryRecord, ...],        ; Array of QueryRecord object
}
```

#### QueryRecord object:

```json
{
  "database": "chinook.sqlite",
  "datetime": "2022-12-01 17:23:41",
  "id": 5,
  "query_time": 1.091406,
  "sql": "SELECT c.customerid,sum(i.total)FROM customers c JOIN invoices i ON c.customerid=i.customerid GROUP BY?ORDER BY?DESC;"
},
```

### Example Request:

```http
GET /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/group/5 HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY3MDAyNzgxOSwibmJmIjoxNjY5OTk3ODE5LCJpYXQiOjE2Njk5OTc4MTl9.MwdQlyGP8YAvoEJ2EayJR7vrD3D0KCxNqiZY7fyzQhw
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.0.0) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 01 Dec 2022 17:31:52 GMT
Content-Length: 693
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "database": "chinook.sqlite",
      "datetime": "2022-12-01 17:23:41",
      "id": 5,
      "query_time": 1.091406,
      "sql": "SELECT c.customerid,sum(i.total)FROM customers c JOIN invoices i ON c.customerid=i.customerid GROUP BY?ORDER BY?DESC;"
    },
    {
      "database": "chinook.sqlite",
      "datetime": "2022-12-01 11:00:28",
      "id": 4,
      "query_time": 1.001799,
      "sql": "SELECT c.customerid,sum(i.total)FROM customers c JOIN invoices i ON c.customerid=i.customerid GROUP BY?ORDER BY?DESC;"
    },
    {
      "database": "chinook.sqlite",
      "datetime": "2022-12-01 10:45:43",
      "id": 1,
      "query_time": 4.617585,
      "sql": "SELECT c.customerid,sum(i.total)FROM customers c JOIN invoices i ON c.customerid=i.customerid GROUP BY?ORDER BY?DESC;"
    }
  ]
}
```