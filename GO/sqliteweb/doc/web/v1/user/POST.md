## Requests

```sh
curl -X "POST" "https://web1.sqlitecloud.io:8443/web/v1/user" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andinux@gmail.com",
  "Company": "sqlitecloud",
  "Referral": "Google",
  "Message": "Optional message"
}'
```

### **POST** - /web/v1/sendmail

### Request object

```json
{
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andinux@gmail.com",
  "Company": "SQLiteCloud",
  "Referral": "Google",
  "Message": "Optional message"
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
POST /web/v1/user HTTP/1.1
Origin: https://sqlitecloud.io
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: RapidAPI/4.0.0 (Macintosh; OS X/13.2.1) GCDHTTPRequest
Content-Length: 144

{
  "FirstName": "Andrea",
  "LastName": "Donetti",
  "Email": "andinux@gmail.com",
  "Company": "sqlitecloud",
  "Referral": "Google",
  "Message": "Optional Message"
}
```

### Example Response:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 16 Feb 2023 18:59:19 GMT
Content-Length: 29
Connection: close

{"status":200,"message":"OK"}
{
  "message": "OK",
  "status": 200
}
```