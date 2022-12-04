# API Documentation

Get a JSON with the reports generated for query specified by the queryID parameter  

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/query/20/suggest" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/dashboard/v1/{projectID}/node/{nodeID}/analyzer/query/{queryID}/suggest

### Query parameters

```json
percentage  = 10                                 -- optional, integer between 0..100 (default = 100)  it represents the percentage of user
table rows that should be considered when generating sqlite_stat1 data
```

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

  value             = [StatementObj, ...],       ; Array of Statement objects
}
```

#### Statement object:

```json
{
  "candidates": "CREATE INDEX albums_idx_cafe21f7 ON albums(ArtistId, Title COLLATE NOCASE); -- stat1: 347 2 1\nCREATE INDEX albums_idx_7acba7bf ON albums(Title COLLATE NOCASE); -- stat1: 347 1\nCREATE INDEX artists_idx_6fd70cd6 ON artists(Name COLLATE NOCASE); -- stat1: 275 1\n",
  "indexes": "CREATE INDEX artists_idx_6fd70cd6 ON artists(Name COLLATE NOCASE);\nCREATE INDEX albums_idx_cafe21f7 ON albums(ArtistId, Title COLLATE NOCASE);\n",
  "plan": "SEARCH tracks USING INDEX IFK_TrackAlbumId (AlbumId=?)\nLIST SUBQUERY 2\nSEARCH albums USING COVERING INDEX albums_idx_cafe21f7 (ArtistId=? AND Title\u003e? AND Title\u003c?)\nLIST SUBQUERY 1\nSEARCH artists USING COVERING INDEX artists_idx_6fd70cd6 (Name\u003e? AND Name\u003c?)\n",
  "sql": "SELECT * FROM tracks WHERE AlbumId IN (SELECT AlbumId FROM albums WHERE ArtistId IN (SELECT ArtistId FROM artists WHERE Name like 'The %') AND Title LIKE \"The %\");"
}
```

### Example Request:

```http
GET /dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/node/6/analyzer/query/20/suggest HTTP/1.1
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
Date: Fri, 02 Dec 2022 16:34:46 GMT
Content-Length: 942
Connection: close

{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "candidates": "CREATE INDEX albums_idx_cafe21f7 ON albums(ArtistId, Title COLLATE NOCASE); -- stat1: 347 2 1\nCREATE INDEX albums_idx_7acba7bf ON albums(Title COLLATE NOCASE); -- stat1: 347 1\nCREATE INDEX artists_idx_6fd70cd6 ON artists(Name COLLATE NOCASE); -- stat1: 275 1\n",
      "indexes": "CREATE INDEX artists_idx_6fd70cd6 ON artists(Name COLLATE NOCASE);\nCREATE INDEX albums_idx_cafe21f7 ON albums(ArtistId, Title COLLATE NOCASE);\n",
      "plan": "SEARCH tracks USING INDEX IFK_TrackAlbumId (AlbumId=?)\nLIST SUBQUERY 2\nSEARCH albums USING COVERING INDEX albums_idx_cafe21f7 (ArtistId=? AND Title\u003e? AND Title\u003c?)\nLIST SUBQUERY 1\nSEARCH artists USING COVERING INDEX artists_idx_6fd70cd6 (Name\u003e? AND Name\u003c?)\n",
      "sql": "SELECT * FROM tracks WHERE AlbumId IN (SELECT AlbumId FROM albums WHERE ArtistId IN (SELECT ArtistId FROM artists WHERE Name like 'The %') AND Title LIKE \"The %\");"
    }
  ]
}
```
