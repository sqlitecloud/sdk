#### Server side Pub/Sub implementation details

SQLiteCloud listens to a default port (default 8860) where clients can connect (using SSL and prior to authentication). The default port is used to read commands from clients, process that commands and sends a reply. Socket flows is:

1. READ FROM SOCKET
2. PROCESS REQUEST
3. SEND A REPLY

When a client send the first **LISTEN channel** command a special reply is sent back to it. This reply is a special command that must be re-executed as is on server side. Command is:

**PAUTH** **client_uuid** **client_secret** 

- When executed on **client-side**, the client opens a new connection to the server and a new thread is created listening for read events on this new socket.
- When executed on **server-side**, client is recognized (using uuid) and authenticated (using secret) and its pubsub socket is set to the current socket used by this new connection. This socket will be used exclusively to deliver notifications. Socket flow is different from the default one because it involves WRITE operations only.

The following commands are related to PUB/SUB:

1. **LISTEN channel**: registers the current client as a listener on the notification channel named `channel`. If the current session is already registered as a listener for this notification channel, nothing is done. If channel has the same name of the current database (if any) table then all WRITE operations will be notified. If channel is `*` the all WRITE operations of all the tables of the current database (if any) will be notified. LISTEN takes effect at transaction commit. If channel is not `*` and it is not the name of table in the current database (if any), then it represents a named channel that can be notified only by NOTIFY channel commands. 
2. **UNLISTEN channel**: remove an existing registration for `NOTIFY` events. `UNLISTEN` cancels any existing registration of the current SQLiteCloud session as a listener on the notification channel named `channel`. The special wildcard `*` cancels all listener registrations for the current session.
3. **NOTIFY channel [, payload]**: The `NOTIFY` command sends a notification event together with an optional "payload" string to each client application that has previously executed `LISTEN channel` for the specified channel name in the current database. The payload (if any) is broadcast as is to all other connections without any modification on server side.



**PUB/SUB FORMAT**

JSON is used to deliver payload to all listening clients. Format changes depending on the operation type. In case of database tables, notifications occur on COMMIT so the same JSON can collect more changes related to that table. Server guarantees **one JSON per channel**.



**1. NOTIFY payload**

Format:
```
{
    sender: "UUID",
    channel: "name",
    type: "MESSAGE",
    payload: "Message content here"	// payload is optional
}
```


**2. TABLE modification payload**
```
{
    sender: "UUID",
    channel: "tablename",
    type: "TABLE",
    pk: ["id", col1"]      // array of primary key name(s)
    payload: [             // array of operations that affect table name
        {
            type: "INSERT",
            id: 12,
            col1: "value1",
            col2: 3.14
        },
        {
            type: "DELETE",
            pv: [13]       // primary key value (s) in the same order as the pk array
        },
        {
            type: "UPDATE",
            id: 15,        // new value
            col1: "newvalue",
            col2: 0.0
            // if primary key is updated during this update then add it to:
            // UPDATE TABLE SET col1='newvalue', col2=0.0, id = 15 WHERE id=14
            pv: [14]       // primary key value (s) set prior to this UPDATE operation
           ]
        }
    ]
}
```

**Details:**

* `sender`: is the UUID of the client who sent the NOTIFY event or who initiated the WRITE operation that triggers the notification. It is common for a client that executes `NOTIFY` to be listening on the same notification channel itself. In that case it will get back a notification event, just like all the other listening sessions. Depending on the application logic, this could result in useless work, for example, reading a database table to find the same updates that that session just wrote out. It is possible to avoid such extra work by noticing whether the notifying `UUID` (supplied in the notification event message) is the same as one's `UUID` (available from SDK). When they are the same, the notification event is one's own work bouncing back, and can be ignored. If `UUID` is 0 it means that server sent that payload.
* `channel`: this field represents the channel/table affected.
* `type`: determine the type of operation, it can be: MESSAGE, TABLE, INSERT, UPDATE, or DELETE (more to come).
* `pk/pv`: these fields represent the primary key name(s) and value(s) affected by this table operation.
* `payload`: todo

**Examples:**
```
USE DATABASE test.sqlite
LISTEN foo
```
```
foo is a table created as:
CREATE TABLE "foo" ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "col1" TEXT, "col2" TEXT)
```

```
DELETE FROM foo WHERE id=14;

{
	"sender": "b7a92805-ef82-4ad1-8c2f-92da6df6b1d5",
	"channel": "foo",
	"type": "TABLE",
	"pk": ["id"],
	"payload": [{
		"type": "DELETE",
		"pv": [14]
	}]
}
```

```
INSERT INTO foo(col1, col2) VALUES ('test100', 'test101');

{
	"sender": "b7a92805-ef82-4ad1-8c2f-92da6df6b1d5",
	"channel": "foo",
	"type": "TABLE",
	"pk": ["id"],
	"payload": [{
		"type": "INSERT",
		"id": 15,
		"col1": "test100",
		"col2": "test101"
	}]
}
```

```
UPDATE foo SET id=14,col1='test200' WHERE id=15;

{
	"sender": "b7a92805-ef82-4ad1-8c2f-92da6df6b1d5",
	"channel": "foo",
	"type": "TABLE",
	"pk": ["id"],
	"payload": [{
		"type": "DELETE",
		"pv": [15]
	}, {
		"type": "INSERT",
		"id": 14,
		"col1": "test200",
		"col2": "test101"
	}]
}
```


#### Useful links:

https://www.postgresql.org/docs/current/sql-notify.html

https://www.postgresql.org/docs/current/sql-listen.html

https://tapoueh.org/blog/2018/07/postgresql-listen-notify/

https://www.postgresql.org/docs/9.1/libpq-example.html

https://www.postgresql.org/docs/9.5/libpq-example.html#LIBPQ-EXAMPLE-2

https://redis.io/topics/pubsub

https://thoughtbot.com/blog/redis-pub-sub-how-does-it-work

https://github.com/redis/hiredis/blob/master/test.c

https://www.toptal.com/go/going-real-time-with-redis-pubsub



#### Implementation:

Use pusher instead of re-implementing it? https://pusher.com

https://making.pusher.com/redis-pubsub-under-the-hood/

https://making.pusher.com/how-pusher-channels-has-delivered-10000000000000-messages/

https://github.com/nanopack/mist

https://github.com/lileio/pubsub

https://itnext.io/redis-as-a-pub-sub-engine-in-go-10eb5e6699cc
