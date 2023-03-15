# SQLite Cloud Javascript Client SDK 

Official SDK repository for SQLite Cloud databases and nodes.

This SDK client library supports web browsers. 

For a sample Web App using this SDK that demonstrates the power of the Pub/Sub capabilities built into SQLite Cloud, check out this [SQLite Cloud Chat](https://chat.sqlitecloud.io/) and the [relative code](https://github.com/sqlitecloud/sdk/tree/master/JS/example/sqlite-cloud-chat).


## Usage Overview

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