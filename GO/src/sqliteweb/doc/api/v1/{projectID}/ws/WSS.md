# API Documentation

Create the main websocket, used to process requests from the Javascript SDK library. 

## WebSocket handshake

`wss://web1.sqlitecloud.io:8443/api/v1/fbf94289-64b0-4fc6-9c20-84083f82ee64/ws?apikey=Rfk00KgQkqIzfqVuTmO87xVLWUwBos3zPzwbXw5UDVy`

### **WSS** - /api/v1/{projectID}/ws?apikey={apikey}

Upgrade to websocket protocol or HTTP error in case of failed authentication or invalid projectID. 

## Websocket Protocol

### Exec 

#### Websocket Request object 
```json
{
    "type":"exec",
    "command":"<sqlitecloud_command>", // example: "LIST DATABASES"
    "id":"Fa3Pe"
}
```

#### Websocket Response object 
Success:
```json
{
  "status": "success",
  "id": "Fa3Pe",
  "type": "exec",
  "data": {
    "columns": [
      "name"
    ],
    "rows": [
      {
        "name": "chinook.sqlite"
      },
      {
        "name": "db1.sqlite"
      },
      {
        "name": "dbempty.sqlite"
      },
      {
        "name": "encdb.sqlite"
      }
    ]
  }
}
```

Error:
```json
{
  "status": "error",
  "id": "Fa3Pe",
  "type": "exec",
  "data": {
    "code": 10002,
    "message": "Unable to find command LIST DATABASE"
  }
}
```


### Listen 

#### Websocket Request object 
```json
{
    "type":"listen",
    "channel":"<channel_name>", // example: "chan1"
    "id":"a4WpY"
}
```

#### Websocket Response object 
Success (pubsub websocket doens't exist). The client must create the pubsub websocket using the endopoint `/api/v1/wspsub?uuid={uuid}&secret={secret}` with the specified uuid and secret values.
```json
{
  "status": "success",
  "id": "a4WpY",
  "type": "listen",
  "data": {
    "uuid": "409b4d68-9b99-43ac-9c90-7370f5936793", 
    "secret": "203e4f8a-0e3d-432b-8fac-b1e8e42c8b90", 
    "channel" : "chan1"
  }
}
```

Success (pubsub websocket already exist). The client will receive notifications in the pubsub websocket.
```json
{
  "status": "success",
  "id": "a4WpY",
  "type": "listen",
  "data": {
    "channel" : "chan1"
  }
}
```

### Unlisten 

#### Websocket Request object 
```json
{
    "type":"unlisten",
    "channel":"<channel_name>", // example: "chan1"
    "id":"5Iyp9"
}
```

#### Websocket Response object 
Success:
```json
{
  "status": "success",
  "id": "5Iyp9",
  "type": "unlisten",
  "data": {
    "channel" : "chan1"
  }
}
```

### Notify 

#### Websocket Request object 
```json
{
    "type":"notify",
    "channel":"<channel_name>", // example: "chan1"
    "payload":"<payload>",      // example: "Hello"
    "id":"q2w15"
}
```

#### Websocket Response object 
Success:
```json
{
  "status": "success",
  "id": "q2w15",
  "type": "notify",
}
```
