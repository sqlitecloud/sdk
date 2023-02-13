# API Documentation

ENCRYPT DATABASE % KEY % 
DECRYPT DATABASE %    

## Requests

```sh
curl -X "PUT" "https://web1.sqlitecloud.io:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/xx/encryption" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "key": ""
}'
```

### **PUT** - /dashboard/v1/{projectID}/database/{databaseName}/encryption

### Request object

```json
{
  key               = "",             // optional: if encryption key is empty, the database will be decrypted, 
                                      // if there is a key present, the database will be encrypted with this key
}
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
PUT /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/database/xx/encryption HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1LjEwMC4zMi4yMjEiLCJleHAiOjE2NDc2MjA5NTcsImp0aSI6IjEiLCJpYXQiOjE2NDc1OTA5NTcsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ3NTkwOTU3LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.erjwvn7RsILHA5cmcrCWdlaOvoyzvysutkab1CGyZGU
Content-Type: application/json; charset=utf-8
Host: web1.sqlitecloud.io:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
Content-Length: 10

{
  "key": ""
}
```

### Example Response (user is in the auth database):

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Fri, 18 Mar 2022 11:58:27 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```