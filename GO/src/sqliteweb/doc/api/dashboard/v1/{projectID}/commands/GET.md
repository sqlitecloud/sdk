# API Documentation

List all builtin commands

## Requests

```sh
curl "https://localhost:8443/dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/commands" \
     -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxMjcuMC4wLjEiLCJleHAiOjE2NDU1NzY5NDUsImp0aSI6IjAiLCJpYXQiOjE2NDU1NDY5NDUsImlzcyI6IlNRTGl0ZSBDbG91ZCBXZWIgU2VydmVyIiwibmJmIjoxNjQ1NTQ2OTQ1LCJzdWIiOiJzcWxpdGVjbG91ZC5pbyJ9.Ru7lvh1tx72CWfsoL2-ZM2b1sB6bB59V6oSlN-gEs2E'
```

### **GET** - /dashboard/v1/{projectID}/commands

### Request object

```code
none
```

### Response object(s)

#### root Response:

```json
{
  status            = 200,                       -- status code: 200 = no error, error otherwise
  message           = "OK",                      -- "OK" or error message

  value             = {},                        -- List with Command objects
}
```

#### Command object :
  
```json
{
    command: "LIST STATS [FROM <start_date> TO <end_date>] [NODE <nodeid>]",
    count: 0,
    avgtime: 0.0,
    privileges: "33554432",
}
```

### Example Request:

```http
GET /dashboard/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/commands HTTP/1.1
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
Date: Tue, 22 Feb 2022 22:22:42 GMT
Content-Length: 83
Content-Type: text/plain; charset=utf-8
Connection: close
{
  "message": "OK",
  "status": 200,
  "value": [
    {
      "avgtime": 0,
      "command": "ADD ALLOWED IP <ip_address> [ROLE <role_name>] [USER <username>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "ADD [LEARNER] NODE <nodeid> ADDRESS <ip_address:port> CLUSTER <ip_address:port>",
      "count": 0,
      "privileges": "33554432"
    },
    {
      "avgtime": 0,
      "command": "APPLY BACKUP SETTINGS",
      "count": 0,
      "privileges": "524288"
    },
    {
      "avgtime": 0.416,
      "command": "AUTH USER <username> PASSWORD <password>",
      "count": 2,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "CLOSE CONNECTION <connectionid> [NODE <nodeid>]",
      "count": 0,
      "privileges": "423100416"
    },
    {
      "avgtime": 0,
      "command": "CREATE DATABASE <database_name> [KEY <encryption_key>] [ENCODING <encoding_value>] [PAGESIZE <pagesize_value>] [IF NOT EXISTS]",
      "count": 0,
      "privileges": "134217728"
    },
    {
      "avgtime": 0,
      "command": "CREATE ROLE <role_name> [PRIVILEGE <privilege_name>] [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "CREATE USER <username> PASSWORD <password> [ROLE <role_name>] [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "DECRYPT DATABASE <database_name>",
      "count": 0,
      "privileges": "134217728"
    },
    {
      "avgtime": 0,
      "command": "DISABLE PLUGIN <plugin_name>",
      "count": 0,
      "privileges": "4194304"
    },
    {
      "avgtime": 0,
      "command": "DISABLE USER <username>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "DROP CLIENT KEY <keyname>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "DROP DATABASE <database_name> KEY <keyname>",
      "count": 0,
      "privileges": "131071"
    },
    {
      "avgtime": 0,
      "command": "DROP DATABASE <database_name> [IF EXISTS]",
      "count": 0,
      "privileges": "268435456"
    },
    {
      "avgtime": 0,
      "command": "DROP KEY <keyname>",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "DROP ROLE <role_name>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "DROP USER <username>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "ENABLE PLUGIN <plugin_name>",
      "count": 0,
      "privileges": "4194304"
    },
    {
      "avgtime": 0,
      "command": "ENABLE USER <username>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "ENCRYPT DATABASE <database_name> WITH KEY <encryption_key>",
      "count": 0,
      "privileges": "134217728"
    },
    {
      "avgtime": 0,
      "command": "GET CLIENT KEY <keyname>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "GET DATABASE <database_name> KEY <keyname>",
      "count": 0,
      "privileges": "131071"
    },
    {
      "avgtime": 0,
      "command": "GET DATABASE [<value>]",
      "count": 0,
      "privileges": "406323200"
    },
    {
      "avgtime": 0,
      "command": "GET INFO <key> [NODE <nodeid>]",
      "count": 0,
      "privileges": "100663296"
    },
    {
      "avgtime": 0,
      "command": "GET KEY <keyname>",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "GET LEADER",
      "count": 0,
      "privileges": "100663296"
    },
    {
      "avgtime": 0,
      "command": "GET RUNTIME KEY <keyname>",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "GET SQL <table_name>",
      "count": 0,
      "privileges": "15"
    },
    {
      "avgtime": 0,
      "command": "GET USER",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "GRANT PRIVILEGE <privilege_name> ROLE <role_name> [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "GRANT ROLE <role_name> USER <username> [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "HELP <<command>>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "LIST ALLOWED IP [ROLE <role_name>] [USER <username>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUP SETTINGS",
      "count": 0,
      "privileges": "524288"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUPS",
      "count": 0,
      "privileges": "524288"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUPS DATABASE <database_name>",
      "count": 0,
      "privileges": "524288"
    },
    {
      "avgtime": 0,
      "command": "LIST CLIENT KEYS",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 4.67,
      "command": "LIST COMMANDS",
      "count": 2,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "LIST CONNECTIONS [NODE <nodeid>]",
      "count": 0,
      "privileges": "423100416"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASE <database_name> KEYS",
      "count": 0,
      "privileges": "131071"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASE CONNECTIONS [ID] <database_name>",
      "count": 0,
      "privileges": "406323200"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASES [DETAILED]",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "LIST INFO",
      "count": 0,
      "privileges": "100663296"
    },
    {
      "avgtime": 0,
      "command": "LIST KEYS [DETAILED]",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "LIST KEYWORDS",
      "count": 0,
      "privileges": "131071"
    },
    {
      "avgtime": 0,
      "command": "LIST LATENCY KEY <keyname> [NODE <nodeid>]",
      "count": 0,
      "privileges": "506986496"
    },
    {
      "avgtime": 0,
      "command": "LIST LATENCY [NODE <nodeid>]",
      "count": 0,
      "privileges": "506986496"
    },
    {
      "avgtime": 0,
      "command": "LIST LOG [FROM <start_date>] [TO <end_date>] [LEVEL <log_level>] [TYPE <log_type>] [ID] [ORDER DESC] [LIMIT <count>] [CURSOR <cursorid>] [NODE <nodeid>]",
      "count": 0,
      "privileges": "406323200"
    },
    {
      "avgtime": 0,
      "command": "LIST NODES",
      "count": 0,
      "privileges": "100663296"
    },
    {
      "avgtime": 0.108,
      "command": "LIST PLUGINS",
      "count": 5,
      "privileges": "4194304"
    },
    {
      "avgtime": 0,
      "command": "LIST PRIVILEGES",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "LIST ROLES",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "LIST RUNTIME KEYS",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "LIST STATS [FROM <start_date> TO <end_date>] [NODE <nodeid>]",
      "count": 0,
      "privileges": "33554432"
    },
    {
      "avgtime": 0,
      "command": "LIST TABLES",
      "count": 0,
      "privileges": "15"
    },
    {
      "avgtime": 0,
      "command": "LIST USERS [WITH ROLES] [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "LISTEN <channel_name>",
      "count": 0,
      "privileges": "131072"
    },
    {
      "avgtime": 0,
      "command": "NOTIFY <channel_name> [<payload_value>]",
      "count": 0,
      "privileges": "262144"
    },
    {
      "avgtime": 0,
      "command": "PING",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "REMOVE ALLOWED IP <ip_address> [ROLE <role_name>] [USER <username>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "REMOVE NODE <nodeid>",
      "count": 0,
      "privileges": "33554432"
    },
    {
      "avgtime": 0,
      "command": "RENAME ROLE <role_name> TO <new_role_name>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "RENAME USER <username> TO <new_username>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "RESTORE BACKUP DATABASE <database_name> [GENERATION <generation>] [INDEX <index>] [TIMESTAMP <timestamp>]",
      "count": 0,
      "privileges": "1048576"
    },
    {
      "avgtime": 0,
      "command": "REVOKE PRIVILEGE <privilege_name> ROLE <role_name> [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "REVOKE ROLE <role_name> USER <username> [DATABASE <database_name>] [TABLE <table_name>]",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "SET CLIENT KEY <keyname> TO <keyvalue>",
      "count": 1,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "SET DATABASE <database_name> KEY <keyname> TO <keyvalue>",
      "count": 0,
      "privileges": "131071"
    },
    {
      "avgtime": 0,
      "command": "SET KEY <keyname> TO <keyvalue>",
      "count": 0,
      "privileges": "8388608"
    },
    {
      "avgtime": 0,
      "command": "SET MY PASSWORD <password>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "SET PASSWORD <password> USER <username>",
      "count": 0,
      "privileges": "16777216"
    },
    {
      "avgtime": 0,
      "command": "SLEEP <ms>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "TEST <test_name> [COMPRESSED]",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "UNLISTEN <channel_name>",
      "count": 0,
      "privileges": "0"
    },
    {
      "avgtime": 0,
      "command": "UNUSE DATABASE",
      "count": 0,
      "privileges": "15"
    },
    {
      "avgtime": 0,
      "command": "USE [OR CREATE] DATABASE <database_name>",
      "count": 0,
      "privileges": "524287"
    }
  ]
}
```