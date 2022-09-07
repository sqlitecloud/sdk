# API Documentation

Create the pubsub websocket, used to send pubsub notifications to the Javascript SDK library. 

## WebSocket handshake

`wss://web1.sqlitecloud.io:8443/api/v1/wspsub?uuid={uuid}&secret={secret}`

### **WSS** - /api/v1/wspsub?uuid={uuid}&secret={secret}

Upgrade to websocket protocol or HTTP error in case of failed authentication or invalid projectID. 
The pubsub websocket is independet of the main websocket. The main websocket can be closed while the pubsub websocket remains active.

## Websocket Protocol

### Write 

The websocket endpoint is not supposed to receive messages, so any massage written by the client causes the endpoint to close the websocket.

### Read

#### Websocket Notification object 
```json
{
  "sender":  "<sender_uuid>",  // string. Example: "409b4d68-9b99-43ac-9c90-7370f5936793"
  "type":    "MESSAGE",
  "channel": "<channel_name>", // string. Example: "chan1"
  "payload": "<message>"       // string. Example: "Hello", or "{\"name\":\"Andrea\",\"msg\":\"Hello\"}"
}
```