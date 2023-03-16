# SQLite Cloud Javascript Client SDK 

Official SDK repository for SQLite Cloud databases and nodes.

This Javascript SDK allows a WebApp to communicate with an SQLite Cloud cluster using 2 channels.
* A **main channel** used to:
  * execute commands
  * get channels list
  * create a new channel
  * remove an existing channel
  * send notifications to a specific channel
  * listen a channel
  * listen a database table
  * unlisten a channel
  * unlisten a table

* **Pub/Sub channel**  used to:
  * receive notifications sent by others on listen channels


## How to use
First of all you have to get your `APP_KEY` and `PROJECT_ID` from the [SQLiteCloud dashboard](https://dashboard.sqlitecloud.io/) associated to an SQLite Cloud cluster.

Once you have your credatials, you can create a new main **main channel**. Every instance of the Javascript Client allows to create a **main channel**.

As describe above **main channel** is used for all outgoing communications from your WebApp.

It is very important to emphasize that once you have created and registered channels or tables for PUB/SUB communications and you are only interested in receiving messages, you can close the **main channel**, leaving open only the **Pub/Sub channel**.


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
`async connect()`|Invoked after initialization opens a new **main channel** connection. Returns how connection process completed.
`close(closePubSub = true)`|Invoked closes both the **main channel** and the **PUB/SUB channel**. If invoked with `closePubSub = false`, closes only the **main channel**. Returns how closing process completed.
`connectionState()`|Returns the actual state of the **main channel** connection.
`pubSubState()`|Returns the actual state of the **PUB/SUB channel** connection.
`requestsStackState()`|Return the lits of pending requests.
`subscriptionsStackState()`|Returns the lits of active subscriptions.
`async exec(command)`|Invoked after connection sends command through the **main channel**. On command exectution success returns the `response`, if not return error.
`async notify(channel, payload)`|Invoked after connection sends notification through the **main channel**. On command exectution success returns the `response`, if not return error.
`async listenChannel(channel, callback)`|Invoked after connection send through the **main channel** the request to start listening for incoming message on the selected channel. It is on the first channel listen request that the SDK open the **PUB/SUB channel**. On the following request the SDK simply add the subscription to the supscriptionStack. For each registered channel is registered the callback to be invoked when a new message arrives. The callback can be different for each channel.  On command exectution success returns the `response`, if not return error.
`async listenTable(channel, callback)`|Invoked after connection send through the **main channel** the request to start listening for incoming message on the selected table. It is on the first table listen request that the SDK open the **PUB/SUB channel**. On the following request the SDK simply add the subscription to the supscriptionStack. For each registered table is registered the callback to be invoked when a new message arrives. The callback can be different for each channel.  On command exectution success returns the `response`, if not return error.
`async unlistenChannel(channel)`|Invoked after connection send through the **main channel** the request to unlistening for incoming message on the selected channel. On command exectution success returns the `response`, if not return error.
`async unlistenTable(table)`|Invoked after connection send through the **main channel** the request to unlistening for incoming message on the selected table. On command exectution success returns the `response`, if not return error.
`async listChannels()`|Invoked after connection send through the **main channel** the request to receive the list of all active channels for the current SQLite Cloud cluster. On command exectution success returns the channels list, if not return error.
`async createChannel(channelName, ifNotExist = true)`|Invoked after connection send through the **main channel** the request to create a new channel with the specified name. On command exectution success returns the `response`, if not return error.
`async removeChannel(channelName)`|Invoked after connection send through the **main channel** the request to remove the channel with the specified name. On command exectution success returns the `response`, if not return error.


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

This method returns the following object

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


