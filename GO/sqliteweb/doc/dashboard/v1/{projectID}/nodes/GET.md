# API Documentation

List all userid projects

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/nodes" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/{projectID}/nodes

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status           = 200,                       ; status code: 200 = no error, error otherwise
  message          = "OK",                      ; "OK" or error message

  value            = {}                         ; Array with node objects
}
```

#### Value object:

```json
{
  id            = 0,                            -- NodeID - 
  name          = "",                           -- Name of this node
  provider      = "",                           -- Provider of this node
  image         = "",                           -- Image data for this node
  region        = "",                           -- Regin data for this node
  size          = "",                           -- Size info for this node
  address       = "",                           -- IPv[4,6] address or host name of this node
  port          = "",                           -- Port this node is listening on
  latitude      = 44.931,       
  longitude     = 10.533,       
  node_id       = 0,                            -- id of the node inside de cluster
  type          = "",                           -- Type fo this node, for example: Leader, Follower, Worker
  status        = "",                           -- progress status of the node, for example: Probe, Replicate, Snapshot (cluster) or Running (nocluster).
  match         = 0,                            -- is the index of the highest known matched raft entry (LIST NODES)
  match_leader  = 0,                            -- is the index of the highest known matched raft entry of the Leader (LIST NODES)
  last_activity = "",                           -- date and time of the last contact with a follower. Leader has NULL. (LIST NODES)
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/nodes HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Date: Tue, 22 Feb 2022 21:15:18 GMT
Content-Length: 560
Content-Type: text/plain; charset=utf-8
Connection: close

{
  "message":"OK",
  "value":[
    {
      "address": "64.227.11.116",
      "details": "i386/1/1MB/100MB",
      "id": 1,
      "last_activity": "2022-04-28 07:58:36",
      "latitude": 40.793,
      "longitude": -74.0247,
      "match": 350,
      "match_leader": 350,
      "name": "Dev1 Server",
      "node_id": 1,
      "port": 9960,
      "status": "Replicate",
      "provider": "DigitalOcean",
      "region": "Rome/Italy",
      "size": "small",
      "type": "Follower",
    },
    ...
  } ],
  "status":200
}

```