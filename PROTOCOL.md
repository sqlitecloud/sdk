## SQLiteCloud Serialization Protocol

SQLite Cloud clients communicate with the SQLite Cloud server using a protocol called SCSP (**S**QLite **C**loud **S**erialization **P**rotocol). This protocol is largely inspired by some wise design decisions made by Redis.

The main design guidelines are:
* Simple to implement even with high-level programming languages
* Fast to parse
* Mostly human-readable
* Immune to big/little-endian processors differences
* Can contain binary data without any specific encoding
* Can support compression

Requests are sent from the client to the SQLite Cloud server as a string representing the command to execute. SQLite Cloud replies with a command-specific data type.

### Networking layer
A client connects to an SQLite Cloud server creating a TCP connection to port 8860 (port number can be changed from the user). While SCSP is technically non-TCP specific, in the context of SQLite Cloud the protocol is only used with TCP connections.

### Request-Response model
SQLite Cloud accepts strings representing the command to execute. These strings are parsed on the server-side and the command-specific reply is sent back to the client. A client is allowed to send to the server only **SCSP Strings**.

Multiple commands can be sent to the SQLite Cloud server as a one **SCSP String** with commands separated by the `;` character. In that case, the client will receive one reply representing the response of the last executed command. In case of error, the execution is interrupted, and the proper error code is returned.

In SCSP, the type of data depends on the first byte:
* For **Strings** the first byte is `+`
* For **Zero-Terminated Strings** the first byte is `!`
* For **Errors** the first byte is `-`
* For **Integer** the first byte is `:`
* For **Float** the first byte is `,`
* For **Blob** the first byte is `$`
* For **Rowset** the first byte is `*`
* For **Rowset chunk** the first byte is `/`
* For **Raw JSON** the first byte is `{` or `[` (NOT IMPLEMENTED)
* For **JSON** the first byte is `#`
* For **NULL** the first byte is `_`
* For **Compressed** rowset the first byte is `%`
* For **PSUB** the first byte is `|`
* For **Command** the first byte is `^`
* For **Reconnect** the first byte is `@`
* For **Array** the first byte is `=`

If the encoding does not include an explicit LEN value then the whole encoded value is terminated by a ` ` space character.

### SCSP Strings
The format is `+LEN STRING`. The whole command is built by four parts:
1. The single `+` character
2. LEN is a string representation of STRING length (theoretically the maximum supported value is UINT64_MAX but it is usually much lower. This value can be customized in SQLite Cloud but it can never exceed that maximum value). LEN does not include the length of the first `+LEN ` part.
3. A single space is used to separate LEN from STRING
4. STRING is the string to represent (string can also contain binary data)

For example to send the string "Hello World!" the command would be: `+12 Hello World!`

### SCSP Zero-Terminated Strings
The format is `!LEN STRING0`. See **SCSP Strings** for details. The only difference is that `STRING` is sent with a 0 at the end, represented in hexadecimal (`\x00`), for better performance in C-like environments that require zero-terminated strings. The LEN field includes the \x00 at the end.

For example, to send the string "Hello World!" the command would be: `!13 Hello World!\x00`.

### SCSP Errors
The format is `-LEN ERRCODE[:EXTCODE[:OFFCODE]] STRING`. Error replies are only sent by the server when something wrong happens. The ERRCODE field in the numeric error code. The optional part, EXTCODE and OFFCODE are numeric values specific to SQLite. EXTCODE represents the extended error code (as specified in the documentation at https://www.sqlite.org/rescode.html), while OFFCODE indicates the offset index within the SQL token where the syntax error occurs (or -1 if none). TThe remaining part of the string corresponds to the error message itself. The error code proves valuable to clients as it enables them to differentiate between various error conditions without resorting to pattern matching within the error message, which may undergo changes. LEN does not include the length of the first `-LEN ` part.

- ERRCODE < 10,000 are SQLite error codes as reported by https://www.sqlite.org/rescode.html
- ERRCODE >= 10,000 and <100,000 are SQLite Cloud error codes
- ERRCODE >= 100,000 are generated internally by the SDK

### SCSP Integer
The format is `:VALUE `. Where `VALUE` is a string representation of the integer value. `VALUE` can be negative and in C it can be parsed using the `strtol/strtoll` API. `VALUE` is guarantee to be an Integer 64 bit number.

### SCSP Float
The format is `,VALUE `. Where `VALUE` is a string representation of the float/double value. In C, `VALUE` can be parsed using the `strtod` API. In this first implementation `VALUE` is guarantee to be a Double number.

### SCSP **Rowset**
The format is `*LEN 0:VERS NROWS NCOLS DATA`. The whole command is built by ten parts:

1. The single `*` character
2. LEN is a string representation of Rowset length (theoretically the maximum supported value is UINT64_MAX but it is usually much lower. LEN does not include the length of the first `*LEN ` part.
3. A single space is used to separate LEN from 0:VERS
4. a single `0:` string followed by a VERS number (a string representation of the number) which specifies the version of the Rowset.
   * `1`: means that only column names are included in the header
   * `2`: means that column names, declared types, database names, table names, origin names, not null flags, primary key flags and autoincrement flags are included in the header (one value for each column)
5. A single space is used to separate 0:VERS from NROWS
6. NROWS  is a string representation of the number of rows contained in the Rowset (can be zero)
7. A single space is used to separate NROWS from NCOLS 
8. NCOLS  is a string representation of the number of columns contained in the Rowset (cannot be zero)
9. A single space is used to separate NCOLS from DATA
10. DATA is a continuos stream of SCSP encoded values:
   1. The first NCOLS are SCSP Strings representing column names
   2. The next NROWS * NCOLS fields are SCSP encoded values

### SCSP **Rowset** Chunk
A Rowset can be sent as a series of multiple chunks (based on user-specific settings) when its size exceeds a pre-defined value. This can be useful to reduce the memory usage (especially on the server-side) with the disadvantage of increasing the time required to process the whole Rowset on the client-side (due to the increased network latency).

The format is `/LEN IDX:VERS NROWS NCOLS DATA`. The command is equal to the SCSP Rowset specification, except for the following differences:

1. The first character is `/`
2. IDX represents the index of the chunk. The first chunk has IDX equals to 1 (0 is reserved for the latest chunk). This value is followed by a fixed `:` character and then by a string representation of the Rowset version (see point 4 of the SCPR Rowset for more information about the version value).
3. NROWS represents the number of rows contained in the chunk. The total number of rows in the final Rowset will be the sum of each NROWS contained in each chunk
4. NCOLS will be the same for all chunks, which means that it does not need to be computed (as a sum) in the final Rowset, and it means that a logical line is never break
5. To mark the end of the Rowset, the special string `/LEN 0 0 0  ` is sent (LEN is always 6 in this case)

When the Rowset is sent in chuck, it is guaranteed that the first chuck contains a complete header and that all the chunks contain complete rows (the individual fields are not truncated in any way).

### SCSP RAW JSON
When the first character is `{` that means that the whole packet is guarantee to be a valid JSON value that can be parsed with a JSON parser.
**This reply is not used in the current implementation.**

### SCSP JSON
The format is `#LEN JSON`. The whole command is built by four parts:
1. The single `#` character
2. LEN is a string representation of JSON length (theoretically the maximum supported value is UINT64_MAX but it is usually much lower. This value can be customized in SQLite Cloud but it can never exceed that maximum value). LEN does not include the length of the first `+LEN ` part.
3. A single space is used to separate LEN from JSON
4. JSON is the string payload

### SCSP NULL
The null type is encoded just as `_ `, which is just the underscore character followed by the ` ` space character.

### SCSP Compression
Any value can be compressed and the format is `%LEN COMPRESSED UNCOMPRESSED BUFFER`. Compression is used ONLY if client notify the server that it supports compression with a dedicated command. Compression is performed using the high performance LZ4 algorithm (https://github.com/lz4/lz4). The whole command is built by six parts:
1. The single `%` character
2. LEN is a string representation of total command length
3. A single space is used to separate LEN from COMPRESSED
4. COMPRESSED is a string representation of the compressed BUFFER size
5. UNCOMPRESSED is a string representation of the uncompressed BUFFER size
6. BUFFER can be any SCSP value containing a UNCOMPRESSED SCSP header followed by LZ4 compressed data. At the time of this writing, the only supported compressed SCSP value is the Rowset. In the chunked version, the `/LEN` header is replaced by `/0` (for performance optimizations on the server side). The data `LEN` value is equivalent to the `UNCOMPRESSED` value of the initial compression header.

### SCSP Command
The format is `^LEN COMMAND`. The whole command is built by four parts:
1. The single `^` character
2. LEN is a string representation of COMMAND length (theoretically the maximum supported value is UINT64_MAX but it is usually much lower. This value can be customized in SQLite Cloud but it can never exceed that maximum value). LEN does not include the length of the first `+LEN ` part.
3. A single space is used to separate LEN from COMMAND
4. COMMAND is the raw string to be executed as is on client side

### SCSP Reconnect
The format is `@LEN COMMAND`. The whole command is built by four parts:
1. The single `^` character
2. LEN is a string representation of COMMAND length (theoretically the maximum supported value is UINT64_MAX but it is usually much lower. This value can be customized in SQLite Cloud but it can never exceed that maximum value). LEN does not include the length of the first `+LEN ` part.
3. A single space is used to separate LEN from COMMAND
4. COMMAND is the raw string to be parsed and to be used to close current connection and reconnect to a new host.
**This reply is not used in the current implementation.**

### SCSP Array
The format is `=LEN N VALUE1 VALUE2 ... VALUEN`. The whole command is built by N+3 parts:
1. The single `=` character
2. LEN is a string representation of the whole command. LEN does not include the length of the first `=LEN ` part.
3. N is the number of items in the array
4. N values separated by a space ` ` character

### SQLite Statements
SQLite statements are usually sent from client to server as `SCSP Strings`.  
In case of bindings the whole statement can be sent as `SCSP Array`.

The server replies to READ statements (like SELECT) with a `SCSP Rowset` or SCSP Rowset Chunk. In case of WRITE statements (like INSERT, UPDATE, DELETE) the SQLite Cloud server replies with a `SCSP Array` in the following format: `=LEN 6 TYPE INDEX ROWID CHANGES TOTAL_CHANGES FINALIZED`:

1. TYPE is always 10 in this case
2. INDEX is always 0 in this case
3. ROWID is the result of the [sqlite3_last_insert_rowid](https://www.sqlite.org/c3ref/last_insert_rowid.html) function
4. CHANGES is the result of the [sqlite3_changes64](https://www.sqlite.org/c3ref/changes.html) function
5. TOTAL_CHANGES is the result of the [sqlite3_total_changes64](https://www.sqlite.org/c3ref/total_changes.html) function
6. FINALIZED is always 1 in this case

---
```Last revision: January 24th, 2024```
