# SQLite Cloud Javascript Client SDK 

Official SDK repository for SQLite Cloud databases and nodes.

This Javascript SDK allows a WebApp to communicate with an SQLite Cloud cluster using 2 WebSocket.
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
For a simple, but comprehensive example of the functionality of this SDK check out this [project]() and the [relative code](https://github.com/sqlitecloud/sdk/tree/master/JS/example/simple).

For a sample WebApp using this SDK that demonstrates the power of the Pub/Sub capabilities built into SQLite Cloud, check out this [SQLite Cloud Chat](https://chat.sqlitecloud.io/) and the [relative code](https://github.com/sqlitecloud/sdk/tree/master/JS/example/sqlite-cloud-chat).


## Topics

The following topics are covered:

* [Installation](https://github.com/sqlitecloud/sdk/tree/master/JS#installation)
  * [Web](https://github.com/sqlitecloud/sdk/tree/master/JS#web)



## Supported platforms

* Web
  * We test against Chrome, Firefox and Safari.
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
const SQLiteCloud = require('sqlitecloud-sdk');
```

#### CDN

```html
<script src="https://js.sqlitecloud.io/1.0/sqlitecloud-sdk.min.js"></script>
```

Then:

```javascript
const SQLiteCloud = window.SQLiteCloud;
```


## Initialization

```js
const client = new SQLiteCloud(PROJECT_ID, API_KEY);
```

Optionally during initialization you can pass two callbacks functions:
* `onErrorCallback` called on WebSocket error event
* `onCloseCallback` called on WebSocket close event

```js
const onErrorCallback = function (event, msg) {
  console.log("WebSocket onError callback:" + msg);
  console.log(event);
}
const onCloseCallback = function (msg) {
  console.log("WebSocket OnClose callback:" + msg);
}
const client = new SQLiteCloud(config.PROJECT_ID, config.API_KEY, onErrorCallback, onCloseCallback);

```

You can get your APP_KEY and PROJECT_ID from the [SQLiteCloud dashboard](https://dashboard.sqlitecloud.io/).


## Configuration

After initializazion it is possibile to configure your client.

#### `SQLiteCloud.setRequestTimeout` (Int value in milliseconds)
Default value is `3000 ms`

#### `SQLiteCloud.setFilterSentMessages` (Boolean)
Default value is `false`

If `true` during PUB/SUB communications library not return messages sent by the user. 


## SDK Methods

**Method**|**Description**
--- | ---
`async connect()`|Creates a new **main WebSocket*. Returns how creation process completed..
`close(closePubSub = true)`|By default, closes both the **main WebSocket** and the **Pub/Sub WebSocket**. If invoked with `closePubSub = false`, closes only the **main WebSocket**. Returns how closing process completed.
`connectionState`|Returns the actual state of the **main WebSocket**.
`pubSubState`|Returns the actual state of the **Pub/Sub WebSocket**.
`async listChannels()`|Uses **main WebSocket** to request the list of all active channels for the current SQLite Cloud cluster. On command execution success returns the channels list, if not return error.
`async createChannel(channelName, ifNotExist = true)`|Uses **main WebSocket** to create a new channel with the specified name. On command exectution success returns the `response`, if not return error.
`async removeChannel(channelName)`|Uses **main WebSocket** to remove the channel with the specified name. On command exectution success returns the `response`, if not return error.
`async exec(command)`|Uses **main WebSocket** to send commands. On command execution success returns the `response`, if not return error.
`async notify(channel, payload)`|Invoked after connection sends notification through the **main WebSocket**. On command exectution success returns the `response`, if not return error.
`async listenChannel(channel, callback)`|Invoked after connection send through the **main WebSocket** the request to start listening for incoming message on the selected channel. It is on the first channel listen request that the SDK open the **Pub/Sub WebSocket**. On the following request the SDK simply add the subscription to the supscriptionStack. For each registered channel is registered the callback to be invoked when a new message arrives. The callback can be different for each channel.  On command exectution success returns the `response`, if not return error.
`async listenTable(channel, callback)`|Invoked after connection send through the **main WebSocket** the request to start listening for incoming message on the selected table. It is on the first table listen request that the SDK open the **Pub/Sub WebSocket**. On the following request the SDK simply add the subscription to the supscriptionStack. For each registered table is registered the callback to be invoked when a new message arrives. The callback can be different for each channel.  On command exectution success returns the `response`, if not return error.
`async unlistenChannel(channel)`|Invoked after connection send through the **main WebSocket** the request to unlistening for incoming message on the selected channel. On command exectution success returns the `response`, if not return error.
`async unlistenTable(table)`|Invoked after connection send through the **main WebSocket** the request to unlistening for incoming message on the selected table. On command exectution success returns the `response`, if not return error.
`requestsStackState()`|Returns the list of pending requests.
`subscriptionsStackState()`|Returns the list of active subscriptions.



### Connection

#### `SQLiteCloud.connect()` 

After initializazion and configuration you can connect invoking the `async` method `SQLiteCloud.connect()`.

```js
async function () {
  const connectionResponse = await client.connect();
  if (connectionResponse.status == 'success') {
    console.log(connectionResponse.data.message);
  } else {
    console.log(connectionResponse.data.message);
  }
}
```

This method returns the following object:

```js
//success or warning response
/*
connectionResponse = {
  status: "success" | "warning"
  data: {
    message: "..."
  }
}
*/

//error response
/*
connectionResponse = {
  status: "error"
  data: error
}
*/

```

### Close

#### `SQLiteCloud.close()` 

To close **main WebSocket** and **PUB/SUB WebSocket** you can invoking the method `SQLiteCloud.close()`.

```js
const close = function (closeAll) {
  //try to close websocket connection
  var closeResponse = client.close(closeAll);
  //check how websocket close completed  
  console.log(closeResponse);
  closeResult.innerHTML = closeResponse.data.message;
  if (closeResponse.status == 'success') {
    //successful websocket close
    logThis(closeResponse.data.message);
  } else {
    //error on websocket close
    logThis(closeResponse.data.message);
  }
}
//close both "main WebSocket" and "PUB/SUB WebSocket"
close(true);
//close only "main WebSocket" leaving open "PUB/SUB WebSocket" to receive incoming messages on subscripted channels and tables 
close(true);
```

This method returns the following object:

```js
//success or error response
/*
connectionResponse = {
  status: "success" | "error"
  data: {
    message: "..."
  }
}
*/
```

### Main WebSocket connection state

#### `SQLiteCloud.connectionState` 

You can monitor the state of **main WebSocket** invoking the method `SQLiteCloud.connectionState`.

```js
setInterval(function () {
  console.log(client.connectionState);
}, 500)
```

### PUB/SUB WebSocket connection state

#### `SQLiteCloud.pubSubState` 

You can monitor the state of **PUB/SUB WebSocket** invoking the method `SQLiteCloud.connectionState`.

```js
setInterval(function () {
  console.log(client.pubSubState);
}, 500)
```

### List Channels

#### `SQLiteCloud.listChannels()` 

You can request the list of all active channels for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.listChannels()`.

```js
async function () {
  const listChannelsResponse = await client.listChannels();
  if (listChannelsResponse.status == 'success') {
    console.log(listChannelsResponse.data);
    var channels = listChannelsResponse.data.rows;
    for (var i = 0; i < channels.length; i++) {
      console.log(channels[i]);
    }    
  } else {
    console.log(listChannelsResponse.data.message);
  }
}
```

This method returns the following object:

```js
//success or warning response
/*
connectionResponse = {
  status: "success"
  data: {
    columns: ['chname'],  
    rows: [
      {chname: ch0},
      {chname: ch1},
      {chname: ch2}
    ],  
  }
}
*/

//error response
/*
connectionResponse = {
  status: "error"
  data: {
    message: "..."
  }
}
*/

```

### Create Channel

#### `SQLiteCloud.createChannel()` 

You can request the creation of a new channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.createChannel()`.

```js
const createChannel = async function (channelName) {
  const createChannelResponse = await client.createChannel(channelName);
  if (createChannelResponse.status == 'success') {
    console.log(createChannelResponse.data);   
  } else {
    console.log(createChannelResponse.data.message);
  }
}
const newChannel = "test-ch";
createChannel(newChannel);
```

This method returns the following object:

```js
//success or warning response
/*
createChannelResponse = {
  status: "success"
  data: "OK"
}
*/

//error response
/*
createChannelResponse = {
  status: "error"
  data: {
    message: "..."
  }
}
*/

```

### Remove Channel

#### `SQLiteCloud.removeChannel()` 

You can request the removal of a channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.removeChannel()`.

```js
const removeChannel = async function (channelName) {
  const removeChannelResponse = await client.removeChannel(channelName);
  if (removeChannelResponse.status == 'success') {
    console.log(removeChannelResponse.data);   
  } else {
    console.log(removeChannelResponse.data.message);
  }
}
const removeChannel = "test-ch";
removeChannel(removeChannel);
```

This method returns the following object:

```js
//success or warning response
/*
removeChannelResponse = {
  status: "success"
  data: "OK"
}
*/

//error response
/*
removeChannelResponse = {
  status: "error"
  data: {
    message: "..."
  }
}
*/

```

### Exec Command

#### `SQLiteCloud.exec()` 

You can execute a command for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.exec()`.

```js
const execCommand = async function (command) {
  const execCommandResponse = await client.exec(command);
  if (execCommandResponse.status == 'success') {
    console.log(execCommandResponse.data);   
  } else {
    console.log(execCommandResponse.data.message);
  }
}
const command = "USE DATABASE db1.sqlite; LIST TABLES PUBSUB";
execCommand(command);
```

This method returns the following object:

```js
//success response
/*
execCommandResponse = {
  status: "success"
  data: [depend on submitted command]
}
*/

//error response
/*
execCommandResponse = {
  status: "error"
  data: {
    code: [int value]
    message: "..."
  }
}
*/

```

### Notify

#### `SQLiteCloud.notify()` 

You can notify a message on an avaible channel for the the current SQLite Cloud cluster invoking the `async` method `SQLiteCloud.notify()`.

```js

```

This method returns the following object:

```js

```