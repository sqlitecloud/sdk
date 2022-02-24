# API Documentation

LIST DATABASES

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/databases" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'

```

### **GET** - /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/databases

### Request object

```json
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 0,                         ; status code: 0 = no error, error otherwise
  message           = "OK",                      ; "OK" or error message

  databases         = db,                        ; Array with Database objects
}
```

#### Database object:

```json
{
  name              = "Db1",
  size              = 18000000000,
  connections       = 5,
  encryption        = nil,
  backup            = "Daily",
  stats             = { 521, 12 },
  bytes             = { 8700000, 712 },
  fragmentation     = { Used = 2400000, total = 712000 }
}
```

### Example Request:

```
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/databases HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU3Mjk5NzAsImp0aSI6IjAiLCJpYXQiOjE2NDU2OTk5NzAsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1Njk5OTcwLCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.8izk50ZCk4kQ7Mpf99tj3DuSOuJhPS2cFpAuhlvoGQw
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response:

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 24 Feb 2022 11:10:46 GMT
Content-Length: 1725
Connection: close

{
  "databases": [
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "Compress",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "dummy",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "Dummy",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "mediastore.sqlite",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "Test",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "tester_linearizable1.sqlite",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "tester_nwriters1.sqlite",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "tester_nwriters2.sqlite",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "X",
      "size": 0,
      "stats": [
        521,
        12
      ]
    },
    {
      "backup": "Daily",
      "bytes": [
        8700000,
        712
      ],
      "connections": 0,
      "encryption": "",
      "fragmentation": {
        "Used": 2400000,
        "total": 712000
      },
      "name": "x",
      "size": 0,
      "stats": [
        521,
        12
      ]
    }
  ],
  "message": "OK",
  "status": 0
}
```

### Previous Response:

```
{
  "ResponseID": 0,
  "Message": "Databases List",
  "Databases": [
      {
          "Name": "Db1",
          "Size": 18000000000,
          "Connections": 5,
          "Encryption": null,
          "Backup": "Daily",
          "Stats": [
              521,
              12
          ],
          "Bytes": [
              8700000,
              712
          ],
          "Fragmentation": {
              "Used": 2400000,
              "Total": 712000
          }
      },
      {
          "Name": "Db1",
          "Size": 18000000000,
          "Connections": 5,
          "Encryption": null,
          "Backup": "Daily",
          "Stats": [
              521,
              12
          ],
          "Bytes": [
              8700000,
              712
          ],
          "Fragmentation": {
              "Used": 2400000,
              "Total": 712000
          }
      },
      {
          "Name": "Db1",
          "Size": 18000000000,
          "Connections": 5,
          "Encryption": null,
          "Backup": "Daily",
          "Stats": [
              521,
              12
          ],
          "Bytes": [
              8700000,
              712
          ],
          "Fragmentation": {
              "Used": 2400000,
              "Total": 712000
          }
      }
  ]
}
```