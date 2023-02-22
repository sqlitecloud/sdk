# API Documentation

Send Recover Email

The Email-Template recover.eml is used. 

Assume, the file path to this endpoint is ./admin/v1/user/{email}/recover/GET.lua, then the template file must be stored in: ./email/v1/recover.eml

The ./admin/v1/user/{email}/recover/GET.lua endpoint sending an email with a link that can be used in the next 10 minutes to reset the password



A note about the recover.eml file.

If the template file is stored in a subdirectory like: ./email/v1/de/recover.eml, it can be used by the LUA command: mail( "recover.eml", "en", template_data )

The "en" in this LUA command should be replaced by the language variable of the user. If the language subdirectory "en" is not found in the path, the system
looks for the recovery.eml file one directory hirachy higher (./email/v1/recover.eml), so in this parent directory, there should always be the default language template (like "en").

Each email can contains two types of content: the text/plain message and/or the text/html message
The mail function will look for the following files:
- `<templatename>.eml` or `<templatename>.txt` for the text/plain part
- `<templatename>.html` for the text/html
One or both parts can be used.

The eml/txt template follows the GO tex/template specifications.
The html template follows the GO html/template specifications.

## Requests

```sh
curl "https://localhost:8443/admin/v1/user/sqlitecloud@synergiezentrum.com/recover" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:password'
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
Authorization: Basic YWRtaW46cGFxxxxxxxx=
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