# SQLite Cloud Javascript Client SDK 

Official SDK repository for SQLite Cloud databases and nodes.

This Javascript SDK allows a WebApp to communicate with an SQLite Cloud cluster using ****2 WebSocket****.
* A **main WebSocket** used to:
  * execute commands
  * get channels list
  * create a new channel
  * remove an existing channel
  * send notifications to a specific channel
  * listen a channel
  * listen a database table
  * unlisten a channel
  * unlisten a table

* **Pub/Sub WebSocket** used to (only after at least one channel/table listening sent command from **main WebSocket**):
  * receive notifications sent by others users


## How to use
First of all you have to get your `APP_KEY` and `PROJECT_ID` associated to an SQLite Cloud cluster from the [SQLiteCloud dashboard](https://dashboard.sqlitecloud.io/).

Once you have your credentials, you can create a new **main WebSocket**. Every instance of the Javascript Client SDK allows to create a **main WebSocket**.

As describe above **main WebSocket** is used for all outgoing communications from your WebApp.

It is very important to emphasize that once you have created and registered channels or tables for PUB/SUB communications, if you are only interested in receiving messages, you can close the **main WebSocket**, leaving open only the **Pub/Sub WebSocket**. 


## Example
For a simple, but comprehensive example of the functionality of this SDK check out this [project](https://jstest.sqlitecloud.io/) and the [relative code](https://github.com/sqlitecloud/sdk/tree/master/JS/example/test-js-sdk).

For a sample WebApp using this SDK that demonstrates the power of the Pub/Sub capabilities built into SQLite Cloud, check out this [SQLite Cloud Chat](https://chat.sqlitecloud.io/) and the [relative code](https://github.com/sqlitecloud/sdk/tree/master/JS/example/sqlite-cloud-chat).


## Topics

The following topics are covered:

* [Installation](https://github.com/sqlitecloud/sdk/tree/master/JS#installation)
  * [Web](https://github.com/sqlitecloud/sdk/tree/master/JS#web)
* [Initialization](https://github.com/sqlitecloud/sdk/tree/master/JS#initialization)
* [Configuration
](https://github.com/sqlitecloud/sdk/tree/master/JS#configuration)
* [SDK Methods
](https://github.com/sqlitecloud/sdk/tree/master/JS#sdk-methods)
  * [Connection](https://github.com/sqlitecloud/sdk/tree/master/JS#connection)
  * [Close](https://github.com/sqlitecloud/sdk/tree/master/JS#close)
  * [Main WebSocket connection state](https://github.com/sqlitecloud/sdk/tree/master/JS#main-websocket-connection-state)
  * [PUB/SUB WebSocket connection state](https://github.com/sqlitecloud/sdk/tree/master/JS#pubsub-websocket-connection-state)
  * [List Channels](https://github.com/sqlitecloud/sdk/tree/master/JS#list-channels)
  * [Create Channel](https://github.com/sqlitecloud/sdk/tree/master/JS#create-channel)
  * [Remove Channel](https://github.com/sqlitecloud/sdk/tree/master/JS#remove-channel)
  * [Exec Command](https://github.com/sqlitecloud/sdk/tree/master/JS#exec-command)
  * [Notify](https://github.com/sqlitecloud/sdk/tree/master/JS#notify)
  * [Listen channel](https://github.com/sqlitecloud/sdk/tree/master/JS#listen-channel)
  * [Listen table](https://github.com/sqlitecloud/sdk/tree/master/JS#listen-table)
  * [Unlisten table](https://github.com/sqlitecloud/sdk/tree/master/JS#unlisten-table)
  * [Main WebSocket pending requests](https://github.com/sqlitecloud/sdk/tree/master/JS#main-websocket-pending-requests)
  * [PUB/SUB WebSocket subscriptions state](https://github.com/sqlitecloud/sdk/tree/master/JS#pubsub-websocket-subscriptions-state)
* [Developing](https://github.com/sqlitecloud/sdk/tree/master/JS#developing)
  * [Building](https://github.com/sqlitecloud/sdk/tree/master/JS#building)



## Supported platforms

* Web
  * We test against Chrome, Firefox and Safari
  * Works in web pages


## Installation

### Web

You can install the library via:

#### NPM (or Yarn)

You can use any NPM-compatible package manager, including NPM itself and Yarn.

```bash
npm install sqlitecloud-sdk
```

Then:

```javascript
import SQLiteCloud from 'sqlitecloud-sdk';
```

Or, if you're not using ES6 modules:

```javascript
var SQLiteCloud = require('sqlitecloud-sdk');
```

#### CDN

```html
<script src="https://js.sqlitecloud.io/1.0/sqlitecloud-sdk.min.js"></script>
```

Then:

```javascript
var SQLiteCloud = window.SQLiteCloud;
```


## Initialization

```js
var client = new SQLiteCloud(PROJECT_ID, API_KEY);
```

Optionally, during initialization you can pass two callbacks functions:
* `onErrorCallback` called on WebSocket error event
* `onCloseCallback` called on WebSocket close event

```js
var onErrorCallback = function (event, msg) {
  console.log("WebSocket onError callback:" + msg);
  console.log(event);
}
var onCloseCallback = function (msg) {
  console.log("WebSocket OnClose callback:" + msg);
}
var client = new SQLiteCloud(PROJECT_ID, API_KEY, onErrorCallback, onCloseCallback);

```

You can get your APP_KEY and PROJECT_ID from the [SQLiteCloud dashboard](https://dashboard.sqlitecloud.io/).


## Configuration

After initialization it is possibile to configure your client.

#### `SQLiteCloud.setRequestTimeout` (Int value in milliseconds)
Default value is `3000 ms`

#### `SQLiteCloud.setFilterSentMessages` (Boolean)
Default value is `false`

If `true` during PUB/SUB communications library does not return sent messages, but only incoming messages. 


## SDK Methods

**Method**|**Description**
--- | ---
`async connect()`|Creates a new **main WebSocket**. Returns how creation process completed.
`close(closePubSub = true)`|By default, closes both the **main WebSocket** and the **Pub/Sub WebSocket**. If invoked with `closePubSub = false`, closes only the **main WebSocket**. Returns how closing process completed.
`connectionState`|Returns the actual state of the **main WebSocket**.
`pubSubState`|Returns the actual state of the **Pub/Sub WebSocket**.
`async listChannels()`|Uses **main WebSocket** to request the list of all active channels for the current SQLite Cloud cluster. If method executes successfully returns the channels list, else returns error.
`async createChannel(channelName, ifNotExist = true)`|Uses **main WebSocket** to create a new channel with the specified name. If method executes successfully returns the `response`, else returns error.
`async removeChannel(channelName)`|Uses **main WebSocket** to remove the channel with the specified name. If method executes successfully returns the `response`, else returns error.
`async exec(command)`|Uses **main WebSocket** to send commands. If method executes successfully returns the `response`, else returns error.
`async notify(channel, payload)`|Uses **main WebSocket** to send notification to an available channel for the current SQLite Cloud cluster. If method executes successfully returns the `response`, else returns error.
`async listenChannel(channel, callback)`|Uses **main WebSocket** to start listening for incoming messages on the selected channel. On the first `listenChannel()` request the SDK creates the **Pub/Sub WebSocket**. On the following `listenChannel()` request the SDK simply adds the new subscription to the `supscriptionStack`. For each listened channel a callback is registered to be invoked when a new message arrives. The callback can be different for each channel. If method executes successfully returns the `response`, else returns error.
`async unlistenChannel(channel)`|Uses **main WebSocket** to stop listening for incoming messages on the selected channel. If method executes successfully returns the `response`, else returns error.
`async listenTable(table, callback)`|Uses **main WebSocket** to start listening for incoming messages on the selected table. On the first `listenTable()` request the SDK creates the **Pub/Sub WebSocket**. On the following `listenTable()` request the SDK simply adds the new subscription to the `supscriptionStack`. For each listened  table a callback is registered to be invoked when a new message arrives. The callback can be different for each table. If method executes successfully returns the `response`, else returns error.
`async unlistenTable(table)`|Uses **main WebSocket** to stop listening for incoming messages on the selected table. If method executes successfully returns the `response`, else returns error.
`requestsStackState`|Returns the list of pending requests.
`subscriptionsStackState`|Returns the list of active subscriptions.


### Connection

#### `async SQLiteCloud.connect`() 

After initialization and configuration you can connect invoking the `async` method `SQLiteCloud.connect()`.

```js
var connect = async function () {
  var response = await client.connect();
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
connect();
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      message: String,
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Close

#### `SQLiteCloud.close`(Boolean) 

You can close **main WebSocket** and **PUB/SUB WebSocket** invoking the method `SQLiteCloud.close(Boolean)`.
- Passing `true` closes both "main WebSocket" and "PUB/SUB WebSocket"
- Passing `false` closes only "main WebSocket" leaving open "PUB/SUB WebSocket" to receive incoming messages on subscripted channels and tables 

```js
var close = function (closeAll) {
  var response = client.close(closeAll);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
//closes both "main WebSocket" and "PUB/SUB WebSocket"
close(true);
//closes only "main WebSocket" leaving open "PUB/SUB WebSocket" to receive incoming messages on subscripted channels and tables 
close(false);
```

This method returns the following object:

```js
/* 
  //on success
  response = {
    status: "success" 
    data: {
      message: String,
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/

```

### Main WebSocket connection state

#### `SQLiteCloud.connectionState` 

You can monitor the state of **main WebSocket** invoking the method `SQLiteCloud.connectionState`.

```js
setInterval(function () {
  var mainWebSocketState = client.connectionState;
  console.log(mainWebSocketState);
}, 500)
```

This method returns the following object:

```js
/*
mainWebSocketState = {
  state: -1 | 0 | 1 | 2 | 3 | 
  description: [string]
}
*/
```

**Method**|**Description**
--- | ---
-1|main WebSocket connection not exist
0|CONNECTING
1|OPEN
2|CLOSING
3|CLOSED

### PUB/SUB WebSocket connection state

#### `SQLiteCloud.pubSubState` 

You can monitor the state of **PUB/SUB WebSocket** invoking the method `SQLiteCloud.connectionState`.

```js
setInterval(function () {
  var pubSubWebSocketState = client.pubSubState;
  console.log(pubSubWebSocketState);
}, 500)
```
This method returns the following object:

```js
/*
pubSubWebSocketState = {
  state: -1 | 0 | 1 | 2 | 3 | 
  description: [string]
}
*/
```

**Method**|**Description**
--- | ---
-1|PubSub WebSocket connection not exist
0|CONNECTING
1|OPEN
2|CLOSING
3|CLOSED

### List Channels

#### `async SQLiteCloud.listChannels`()

You can request the list of all active channels for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.listChannels()`.

```js
var listChannels =  async function () {
  var response = await client.listChannels();
  if (response.status == 'success') {
    var channels = response.data.rows;
    for (var i = 0; i < channels.length; i++) {
      //do your stuff  
    }    
  }
  if (response.status == 'error') {
    //error handling
  }  
}
listChannels();
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success",
    data: {
      columns: ['chname'],  
      rows: [
        {chname: ch0},
        {chname: ch1},
        {chname: ch2}
      ],  
    }
  }

  //on success
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/

```

### Create Channel

#### `async SQLiteCloud.createChannel`(String) 

You can request the creation of a new channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.createChannel(String)`.

```js
var createChannel = async function (channelName) {
  var response = await client.createChannel(channelName);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var newChannel = "test-ch";
createChannel(newChannel);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      message: String,
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/

```

### Remove Channel

#### `SQLiteCloud.removeChannel`(String) 

You can request the removal of a channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.removeChannel(String)`.

```js
var removeChannel = async function (channelName) {
  var response = await client.removeChannel(channelName);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }  
}
var removeChannel = "test-ch";
removeChannel(removeChannel);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      message: String,
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Exec Command

#### `SQLiteCloud.exec`(String) 

You can execute a command for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.exec(String)`.

```js
var exec = async function (command) {
  var response = await client.exec(command);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var command = "USE DATABASE db1.sqlite; LIST TABLES PUBSUB";
exec(command);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: Object //depends on submited command
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Notify

#### `SQLiteCloud.notify`(String, Object) 

You can notify a message on an available channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.notify(String, Object)`.

```js
var notify = async function (channel, payload) {
  var response = await client.notify(channel, payload);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var channel = "test channel";
var payload = { message: "hello world" };
notify(channel, payload);

```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      message: String,
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Listen channel

#### `SQLiteCloud.listenChannel`(String, Callback) 

You can start listening for incoming messages on an available channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.listenChannel(String, Callback)`. You have to provide a callback function to consume the incoming messages.

```js
var listenChannel = async function (channel, callback) {
  var response = await client.listenChannel(channel, callback);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var channel = "test channel";
var newMessageCallback = function(incomingMessage) {
  //do your stuff 
}
listenChannel(channel, newMessageCallback);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      channel: String, //the name of the channel started listening correctly
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

The callback returns the following object:

```js
/*
  incomingMessage = {
    channel: String, //the name of the channel that received the message
    ownMessage: Boolean,//true if the user that sent that message is the same that is receiving the message
    payload: {
      message: String, //text of the incoming message
    },
    sender: String //ID of the sender
  }
*/
```


### Unlisten channel

#### `SQLiteCloud.unlistenChannel`(String) 

You can stop listening for incoming messages on an available channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.unlistenChannel(String)`.

```js
var unlistenChannel = async function (channel) {
  var response = await client.unlistenChannel(channel);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var channel = "test channel";
unlistenChannel(channel);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      //TODO
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Listen table

#### `SQLiteCloud.listenTable`(String, Callback) 

You can start listening for incoming messages on an available database table for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.listenTable(String, Callback)`. You have to provide a callback function to consume the incoming messages.

```js
var listenTable = async function (table, callback) {
  var response = await client.listenTable(table, callback);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var table = "test table";
var newMessageCallback = function(incomingMessage) {
  //do your stuff 
}
listenTable(channel, newMessageCallback);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      channel: String, //the name of the table started listening correctly
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

The callback returns the following object:

```js
/*
  incomingMessage = {
    channel: String, //the name of the channel that received the message
    ownMessage: Boolean,//true if the user that sent that message is the same that is receiving the message 
    payload: Object, //the object representing the change made to the table //TODO
    sender: String //ID of the sender
  }
*/
```

### Unlisten table

#### `SQLiteCloud.unlistenTable`(String) 

You can stop listening for incoming messages on an available database table for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.unlistenTable(String)`.

```js
var unlistenTable = async function (table) {
  var response = await client.unlistenTable(table);
  if(response.status == "success"){
    //do your stuff 
  }
  if(response.status == "error"){
    //error handling
  }
}
var table = "test table";
unlistenTable(table);
```

This method returns the following object:

```js
/*
  //on success
  response = {
    status: "success" 
    data: {
      //TODO
    }
  }

  //on error
  response = {
    status: "error" 
    data: {
      message: String,
      error: Error //optional
    }
  }
*/
```

### Main WebSocket pending requests 

#### `SQLiteCloud.requestsStackState` 

You can monitor the pending requests sent on **main WebSocket** invoking the method `SQLiteCloud.requestsStackState`.

```js
setInterval(function () {
  var pendingRequests = client.requestsStackState;
  //do yout stuff
}, 500)
```

This method returns a `Map` containing all the IDs of the pending requests.

### PUB/SUB WebSocket subscriptions state 

#### `SQLiteCloud.subscriptionsStackState` 

You can monitor which channels and tables you are actually listening to on **PUB/SUB WebSocket** invoking the method `SQLiteCloud.subscriptionsStackState`.

```js
setInterval(function () {
  var subscriptionsStack = client.subscriptionsStackState;
  //do yout stuff
}, 500)
```

This method returns a `Map` containing all the channels and tables currently listening.

## Developing
Install all dependencies via npm:

```bash
npm install
```
Run a development server which serves bundled javascript from <http://localhost:8080/sqlitecloud-sdk.js> so that you can edit files in /src freely.

```bash
npm run start
```

### Building
In order to build run:

```bash
npm run build //not minified version sqlitecloud-sdk.js
npm run buildMinify //minified version sqlitecloud-sdk.min.js
```
