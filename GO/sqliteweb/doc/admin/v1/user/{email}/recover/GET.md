# API Documentation

Send Recover Email

The Email-Template recover.eml is used. 

Assume, the file path to this endpoint is ./admin/v1/user/{email}/recover/GET.lua, then the template file must be stored in: ./email/v1/recover.eml

The ./admin/v1/user/{email}/recover/GET.lua endpoint is sending the password that is stored in the admin.USER database to the given email address.
This behaviour is not production ready. Instead, there should be a system that sends a first Email with: "A password recovery was requested, if not from you ignore it, if for you klick link". This requires another Endpoint that takes the klick and then generates a random password that is sent in a second email to the user. Any other password recover system is fine, too...

Since this endpoint is "experimental" and just to demonstrate a basic password recovery functionality (without propper GUI user interface), it does not take the
"last_recovery_request" column in the auth.USER database table into account. This field is there to limit the number of password requests in a certain time period.

A note about the recover.eml file.

If the template file is stored in a subdirectory like: ./email/v1/de/recover.eml, it can be used by the LUA command: mail( "recover.eml", "en", template_data )

The "en" in this LUA command should be replaced by the language variable of the user. If the language subdirectory "en" is not found in the path, the system
looks for the recovery.eml file one directory hirachy higher (./email/v1/recover.eml), so in this parent directory, there should always be the default language template (like "en").

The recover.eml template follows the GO template specifications. More info can be found here:

https://pkg.go.dev/text/template

To generate a rich media Email template, just use any Email program, design a rich recovery email (with embedded immages) and send this email to youself. Then open the email and look at the source code. Cut/Copy/Paste this sourcecode into the recover.eml file. Replace the variable data with tamplate variables.

## Requests

```sh
curl "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com/recover" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password' \
     -d $'{}'
```

### **GET** - /admin/v1/user/{email}/recover

### Request object

```code
none
```
### Response object(s)

#### root Response:

```code
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message
}
```

### Example Request:

```
GET /admin/v1/user/sqlitecloud@synergiezentrum.com/recover HTTP/1.1
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
Content-Type: application/json; charset=utf-8
Host: localhost:8443
Connection: close
User-Agent: Paw/3.3.6 (Macintosh; OS X/10.14.6) GCDHTTPRequest
```

### Example Response :

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
Access-Control-Allow-Origin: *
Content-Encoding: utf-8
Content-Type: application/json
Date: Thu, 28 Apr 2022 13:25:48 GMT
Content-Length: 29
Connection: close

{
  "message": "OK",
  "status": 200
}
```