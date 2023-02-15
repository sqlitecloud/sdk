## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/web/v1/sendmail" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "Template": "contactus.eml",
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andrea@sqlitecloud.io",
  "Company": "sqlitecloud",
  "Message": "First line of the message\nSecond line of the message.\nCiao"
}'
```

### **POST** - /web/v1/sendmail

### Request object

```json
{
  "Template": "contactus.eml",
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andrea@sqlitecloud",
  "Company": "sqlitecloud",
  "Message": "First line of the message\nSecond line of the message.\nCiao"
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
POST /web/v1/sendmail HTTP/1.1
Origin: https://sqlitecloud.io
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.1.0) GCDHTTPRequest
Content-Length: 196

{
  "Template": "contactus.eml",
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andrea@sqlitecloud",
  "Company": "sqlitecloud",
  "Message": "First line of the message\nSecond line of the message.\nCiao"
}
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Content-Encoding: utf-8
Content-Type: application/json
Date: Tue, 14 Feb 2023 16:24:55 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```