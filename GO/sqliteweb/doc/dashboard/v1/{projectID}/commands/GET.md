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
    privileges: "CLUSTERADMIN",
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
  "value": {
  "message": "OK",
  "status": 200,
  "value": [
    {
      "avgtime": 0,
      "command": "ADD ALLOWED IP \u003cip_address\u003e [ROLE \u003crole_name\u003e] [USER \u003cusername\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "ADD [LEARNER] NODE \u003cnodeid\u003e ADDRESS \u003cip_address:port\u003e CLUSTER \u003cip_address:port\u003e",
      "count": 0,
      "privileges": "CLUSTERADMIN"
    },
    {
      "avgtime": 0,
      "command": "APPLY BACKUP SETTINGS",
      "count": 0,
      "privileges": "BACKUP"
    },
    {
      "avgtime": 1.21,
      "command": "AUTH USER \u003cusername\u003e PASSWORD \u003cpassword\u003e",
      "count": 1,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "CLOSE CONNECTION \u003cconnectionid\u003e [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,USERADMIN,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "CREATE DATABASE \u003cdatabase_name\u003e [KEY \u003cencryption_key\u003e] [ENCODING \u003cencoding_value\u003e] [PAGESIZE \u003cpagesize_value\u003e] [IF NOT EXISTS]",
      "count": 0,
      "privileges": "CREATE_DATABASE"
    },
    {
      "avgtime": 0,
      "command": "CREATE ROLE \u003crole_name\u003e [PRIVILEGE \u003cprivilege_name\u003e] [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "CREATE USER \u003cusername\u003e PASSWORD \u003cpassword\u003e [ROLE \u003crole_name\u003e] [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "DECRYPT DATABASE \u003cdatabase_name\u003e",
      "count": 0,
      "privileges": "CREATE_DATABASE"
    },
    {
      "avgtime": 0,
      "command": "DISABLE PLUGIN \u003cplugin_name\u003e",
      "count": 0,
      "privileges": "PLUGIN"
    },
    {
      "avgtime": 0,
      "command": "DISABLE USER \u003cusername\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "DROP CLIENT KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "DROP DATABASE \u003cdatabase_name\u003e KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN"
    },
    {
      "avgtime": 0,
      "command": "DROP DATABASE \u003cdatabase_name\u003e [IF EXISTS]",
      "count": 0,
      "privileges": "DROP_DATABASE"
    },
    {
      "avgtime": 0,
      "command": "DROP KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "DROP ROLE \u003crole_name\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "DROP USER \u003cusername\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "ENABLE PLUGIN \u003cplugin_name\u003e",
      "count": 0,
      "privileges": "PLUGIN"
    },
    {
      "avgtime": 0,
      "command": "ENABLE USER \u003cusername\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "ENCRYPT DATABASE \u003cdatabase_name\u003e WITH KEY \u003cencryption_key\u003e",
      "count": 0,
      "privileges": "CREATE_DATABASE"
    },
    {
      "avgtime": 0,
      "command": "GET CLIENT KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "GET DATABASE \u003cdatabase_name\u003e KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN"
    },
    {
      "avgtime": 0,
      "command": "GET DATABASE [\u003cvalue\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "GET INFO \u003ckey\u003e [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "CLUSTERADMIN,CLUSTERMONITOR"
    },
    {
      "avgtime": 0,
      "command": "GET KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "GET LEADER",
      "count": 0,
      "privileges": "CLUSTERADMIN,CLUSTERMONITOR"
    },
    {
      "avgtime": 0,
      "command": "GET RUNTIME KEY \u003ckeyname\u003e",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "GET SQL \u003ctable_name\u003e",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE"
    },
    {
      "avgtime": 0,
      "command": "GET USER",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "GRANT PRIVILEGE \u003cprivilege_name\u003e ROLE \u003crole_name\u003e [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "GRANT ROLE \u003crole_name\u003e USER \u003cusername\u003e [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "HELP \u003c\u003ccommand\u003e\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "LIST ALLOWED IP [ROLE \u003crole_name\u003e] [USER \u003cusername\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUP SETTINGS",
      "count": 0,
      "privileges": "BACKUP"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUPS",
      "count": 0,
      "privileges": "BACKUP"
    },
    {
      "avgtime": 0,
      "command": "LIST BACKUPS DATABASE \u003cdatabase_name\u003e",
      "count": 0,
      "privileges": "BACKUP"
    },
    {
      "avgtime": 0,
      "command": "LIST CLIENT KEYS",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "LIST COMMANDS [DETAILED]",
      "count": 1,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "LIST CONNECTIONS [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,USERADMIN,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASE \u003cdatabase_name\u003e KEYS",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASE CONNECTIONS [ID] \u003cdatabase_name\u003e",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST DATABASES [DETAILED]",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "LIST INFO",
      "count": 0,
      "privileges": "CLUSTERADMIN,CLUSTERMONITOR"
    },
    {
      "avgtime": 0,
      "command": "LIST KEYS [DETAILED]",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "LIST KEYWORDS",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST LATENCY KEY \u003ckeyname\u003e [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,CLUSTERADMIN,CLUSTERMONITOR,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST LATENCY [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,CLUSTERADMIN,CLUSTERMONITOR,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST LOG [FROM \u003cstart_date\u003e] [TO \u003cend_date\u003e] [LEVEL \u003clog_level\u003e] [TYPE \u003clog_type\u003e] [ID] [ORDER DESC] [LIMIT \u003ccount\u003e] [CURSOR \u003ccursorid\u003e] [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "BACKUP,RESTORE,DOWNLOAD,CREATE_DATABASE,DROP_DATABASE,HOSTADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST NODES",
      "count": 0,
      "privileges": "CLUSTERADMIN,CLUSTERMONITOR"
    },
    {
      "avgtime": 0,
      "command": "LIST PLUGINS",
      "count": 0,
      "privileges": "PLUGIN"
    },
    {
      "avgtime": 0,
      "command": "LIST PRIVILEGES",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST ROLES",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST RUNTIME KEYS",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "LIST STATS [FROM \u003cstart_date\u003e TO \u003cend_date\u003e] [NODE \u003cnodeid\u003e]",
      "count": 0,
      "privileges": "CLUSTERADMIN"
    },
    {
      "avgtime": 0,
      "command": "LIST TABLES",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE"
    },
    {
      "avgtime": 0,
      "command": "LIST USERS [WITH ROLES] [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "LISTEN \u003cchannel_name\u003e",
      "count": 0,
      "privileges": "SUB"
    },
    {
      "avgtime": 0,
      "command": "LOAD PLUGIN \u003cplugin_name\u003e",
      "count": 0,
      "privileges": "PLUGIN"
    },
    {
      "avgtime": 0,
      "command": "NOTIFY \u003cchannel_name\u003e [\u003cpayload_value\u003e]",
      "count": 0,
      "privileges": "PUB"
    },
    {
      "avgtime": 0,
      "command": "PING",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "REMOVE ALLOWED IP \u003cip_address\u003e [ROLE \u003crole_name\u003e] [USER \u003cusername\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "REMOVE NODE \u003cnodeid\u003e",
      "count": 0,
      "privileges": "CLUSTERADMIN"
    },
    {
      "avgtime": 0,
      "command": "RENAME ROLE \u003crole_name\u003e TO \u003cnew_role_name\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "RENAME USER \u003cusername\u003e TO \u003cnew_username\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "RESTORE BACKUP DATABASE \u003cdatabase_name\u003e [GENERATION \u003cgeneration\u003e] [INDEX \u003cindex\u003e] [TIMESTAMP \u003ctimestamp\u003e]",
      "count": 0,
      "privileges": "RESTORE"
    },
    {
      "avgtime": 0,
      "command": "REVOKE PRIVILEGE \u003cprivilege_name\u003e ROLE \u003crole_name\u003e [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "REVOKE ROLE \u003crole_name\u003e USER \u003cusername\u003e [DATABASE \u003cdatabase_name\u003e] [TABLE \u003ctable_name\u003e]",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0.001,
      "command": "SET CLIENT KEY \u003ckeyname\u003e TO \u003ckeyvalue\u003e",
      "count": 1,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "SET DATABASE \u003cdatabase_name\u003e KEY \u003ckeyname\u003e TO \u003ckeyvalue\u003e",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN"
    },
    {
      "avgtime": 0,
      "command": "SET KEY \u003ckeyname\u003e TO \u003ckeyvalue\u003e",
      "count": 0,
      "privileges": "SETTINGS"
    },
    {
      "avgtime": 0,
      "command": "SET MY PASSWORD \u003cpassword\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "SET PASSWORD \u003cpassword\u003e USER \u003cusername\u003e",
      "count": 0,
      "privileges": "USERADMIN"
    },
    {
      "avgtime": 0,
      "command": "SLEEP \u003cms\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "TEST \u003ctest_name\u003e [COMPRESSED]",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "UNLISTEN \u003cchannel_name\u003e",
      "count": 0,
      "privileges": "NONE"
    },
    {
      "avgtime": 0,
      "command": "UNUSE DATABASE",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE"
    },
    {
      "avgtime": 0,
      "command": "USE [OR CREATE] DATABASE \u003cdatabase_name\u003e",
      "count": 0,
      "privileges": "READ,INSERT,UPDATE,DELETE,READWRITE,PRAGMA,CREATE_TABLE,CREATE_INDEX,CREATE_VIEW,CREATE_TRIGGER,DROP_TABLE,DROP_INDEX,DROP_VIEW,DROP_TRIGGER,ALTER_TABLE,ANALYZE,ATTACH,DETACH,DBADMIN,SUB,PUB,PUBSUB"
    }
  ]
}
```