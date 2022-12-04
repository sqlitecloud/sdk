# API Documentation

Get a JSON with the list of queries recorded by the query analyzer, grouped by database and sql (normalized_sql)

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/{projectID}/node/{nodeID}/analyzer

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

  value             = [QueryGroupObj, ...],      ; Array of QueryGroup objects
}
```

#### QueryGroup object:

```json
{
  "AVG(query_time)": 551.5748515,
  "COUNT(query_time)": 2,
  "MAX(query_time)": 552.518287,
  "database": "chinook.sqlite",
  "group_id": 9,
  "sql": "WITH RECURSIVE cnt(x)AS(SELECT?UNION ALL SELECT x+?FROM cnt LIMIT?)SELECT x FROM cnt WHERE x=?;"
},
```

### Example Request:

```http
GET /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer HTTP/1.1
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
Date: Thu, 28 Apr 2022 08:13:15 GMT
Connection: close
Transfer-Encoding: chunked

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "AVG(query_time)": 551.5748515,
      "COUNT(query_time)": 2,
      "MAX(query_time)": 552.518287,
      "database": "chinook.sqlite",
      "group_id": 9,
      "sql": "WITH RECURSIVE cnt(x)AS(SELECT?UNION ALL SELECT x+?FROM cnt LIMIT?)SELECT x FROM cnt WHERE x=?;"
    },
    {
      "AVG(query_time)": 2.5735045,
      "COUNT(query_time)": 2,
      "MAX(query_time)": 4.850398,
      "database": "chinook.sqlite",
      "group_id": 21,
      "sql": "SELECT*FROM tracks WHERE albumid IN(SELECT albumid FROM albums WHERE artistid IN(SELECT artistid FROM artists WHERE name LIKE?)AND title LIKE?);"
    },
    {
      "AVG(query_time)": 0.51313925,
      "COUNT(query_time)": 4,
      "MAX(query_time)": 0.871526,
      "database": "chinook.sqlite",
      "group_id": 19,
      "sql": "SELECT c.customerid,sum(i.total)FROM customers c JOIN invoices i ON c.customerid=i.customerid GROUP BY?ORDER BY?DESC;"
    }
  ]
}
```