# SQLiteWeb Server
## Getting started

### Requirements
1) Setup your GO environment:

```console
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH
```
This code snipped should output something like: `/Users/<your home directory>/GitHub/SqliteCloud/sdk/GO`

2) Create a id_rsa.pub on your machine:

```console
make id_rsa.pub
```

Don't worry - if you have done this already, the Makefile will detect this and leave your `~/.ssh/id_rsa.pub` file un-touched.

3) Install this file on the server:

```console
make web_install
```

You will have to enter the root password for the server `web1.sqlitecloud.io`. If you want to use your own login name, change the Makefile accordingly. (replace `root@web1.sqlitecloud.io` with `<your login name>@web1.sqlitecloud.io`).

After this, you are ready to work with the SQLiteWeb server!

### Compiling
```console
make bin/sqliteweb_linux
```
Will compile a fresh Linux binary. Add GOARCH=amd64 if your machine is not x86-64. You can also build binaries for other platform's and OS'es (on your machine) with:

```console
make bin/sqliteweb_osx
make bin/sqliteweb_win
```

accordingly.

### Setup the server file/folder structure

The **SQLiteWeb Server** requires a certain file/folder structure on the target machine:

```console
/opt/sqliteweb/
/opt/sqliteweb/www
/opt/sqliteweb/api/v1
/opt/sqliteweb/sbin
/opt/sqliteweb/etc
/opt/sqliteweb/etc/sqliteweb
/opt/sqliteweb/etc/sqliteweb/sqliteweb.ini
/opt/sqliteweb/etc/sqliteweb/certs
/opt/sqliteweb/etc/sqliteweb/certs/chain.pem
/opt/sqliteweb/etc/sqliteweb/certs/privkey.pem
/opt/sqliteweb/etc/init.d
/opt/sqliteweb/etc/init.d/sqliteweb
```

You then have to link the init.d script **ON THE SERVER** into the right place with:

```console
ON SERVER> ln -s /opt/sqliteweb/etc/init.d/sqliteweb /etc/.
```

## Installing/Updating the SQLiteWeb Server
Now, you can upload/update the previously compiled server executable to the server with a:

```console
make web_update
```

This command will at first stop the SQLiteWeb server on your remote host, then compile a fresh local version (if necessary) and after this install the Linux binary on the remote host. If everything went well, the new server is then started on the remote host again.

## Controlling the SQLiteWeb Server
You can **remote-control** the SQLiteWeb server with the following commands

```console
make web_stop
make web_start
make web_restart
```

## Testing the SQLiteWeb Server
To test the SQLiteWeb Server, enter a quick:

```console
make web_test
```
If everything went well, you should see an output like:

```console
Ping...success.
Auth...success.
```

# Using the SQLiteWeb Server
The following tasks have to be done **ON the remote host** where the SQLiteWeb Server is running (not your local machine).

## Command line arguments:
To see all possible command line arguments, enter the following command:

```console
/opt/sqliteweb/sbin/sqliteweb_linux --help
SQLite Cloud Web Server 

Usage:
  sqliteweb options
  sqliteweb -?|--help|--version

Examples:
  sqliteweb --config=../etc/sqliteweb.ini
  sqliteweb --version
  sqliteweb -?

General Options:
  --config=<PATH>          Use config file in <PATH> [default: /etc/sqliteweb/sqliteweb.ini]
  -?, --help               Show this screen
  --version                Display version information

Connection Options:
  -a, --address IP         Use IP address [default: 0.0.0.0]
  -p, --port PORT          Use Port [default: 8433]
  -c, --cert <FILE>        Use certificate chain in <FILE>        
  -k, --key <FILE>         Use private certificate key in <FILE>

Server Options:
  --www=<PATH>             Server static web sites from <PATH>
  --api=<PATH>             Server dummy REST stubs from <PATH>


```

The basic idea here is, that all parameters are configured in a config.ini file. Then, for quick test purposes, the most important parameters can be overwritten by command line arguments.


The default place where SQLiteWeb is looking for it's config.ini file is `/etc/sqliteweb/sqliteweb.ini`, but this config file can be located anywhere (like in /etc/sqliteweb.ini) or in `/opt/sqliteweb/etc/sqliteweb/sqliteweb.ini`. You can specify which configuration you want to use with the `--config` parameter.

For example, a typical start command for **debugging** would look similar like this:

```console
ON SERVER> /opt/sqliteweb/sbin/sqliteweb_linux --config=/opt/sqliteweb/etc/sqliteweb/sqliteweb.ini
```

## Configuration
The SQLiteWeb Server is configured with a config file, normally located under: `/opt/sqliteweb/etc/sqliteweb/sqliteweb.ini` or `/etc/sqliteweb/sqliteweb.ini`. The config file looks like this:

```console
[server]
  address = 0.0.0.0
  port    = 8443
  
  hostname    = web1.sqlitecloud.io
  cert_chain  = /opt/sqliteweb/etc/sqliteweb/certs/chain.pem
  cert_key    = /opt/sqliteweb/etc/sqliteweb/certs/privkey.pem
  
  logfile     = /var/log/sqliteweb.log
  clflogfile  = /var/log/sqliteweb-clf.log

[auth]
  jwt_key     = "my_secret_iwt_key"
  jwt_ttl     = 300

  host        = auth1.sqlitecloud.io
  port        = 8860
  login       = admin
  password    = secret

[www]
  path 	      = /opt/sqliteweb/www

[api]
  path        = /opt/sqliteweb/api
```

If you have made changes in the config file, you have to restart the server to make your changes take effect. You can restart the server with:

```console
ON SERVER> /etc/init.d/sqliteweb restart
```
or

```console
ON YOU LOCAL MACHINE> make web_restart
```

### The [server] section of the configuration file
- address: Sets the interface the server should use to serve the GUI and it's API: Common values are: 0.0.0.0 (serve on all interfaces), 127.0.0.1 (serve only on localhost), &quote network card ip address &quote (serve to the outside world).
- port: Sets the server port to use (default is 8443)
- hostname: This is important for the ssl encryption. Please set it to the name that the clients use to connect to this host.
- cert_chain: This is the path to the certificate PEM file.
- cert_key: This is the path to the key PEM file.
- logfile: This is the path to the file where SQLiteWeb should write it's log messages.'

### The [auth] section of the configuration file
- jwt_key: This is a static string of any length and complexity that is used as a secret to sign the JWT Access Token.
- jwt_ttl: This is the TimeToLive before a JWT Token will auto-expire.
- host: This is the hostname of the user authentication server. This server must be another SQLiteCloud instance with the user credentials table (users.sqlite).
- port: This is the port of the user authentication server (default is 8860)
- login: This is the login name for logging in to the user authentication server (default is admin).
- password: This is the password for logging into the user authentication server.

### The [www] section of the configuration file
- path: This is the path where (static) web-resources are served from. To access those resources, point your browser to the hostname and port specified in the [server] section. Example: [https://web1.sqlitecloud.io:8443/](https://web1.sqlitecloud.io:8443/)
or [https://web1.sqlitecloud.io:8443/firework/](https://web1.sqlitecloud.io:8443/firework/)


### The (dummy) [api] section of the configuration file
- path: This is the folder path where the (dummy) API requests are specified in the form of the directory structure and the responses are specified by <HTTP VERB>.json. files. To access those dummy request/response pairs, point your browser, or JSON clien to the hostname and port specified in the [server] section and add the path: `/api/vi/` to it. Example: [https://web1.sqlitecloud.io:8443/api/v1/ping](https://web1.sqlitecloud.io:8443/api/v1/ping)

**Please note: The path of the endpoint should start with: /api/v1/...**

## Serving the REACT GUI
Put all of your REACT files into the specified www.path folder (normally: /opt/sqliteweb/www). The effect of uploading new files is immediately, no server reload is necessary. A typical www folder contents could look like this for example:


```console
/opt/sqliteweb/www/manifest.json
/opt/sqliteweb/www/static
/opt/sqliteweb/www/static/css
/opt/sqliteweb/www/static/css/main.a2731a96.chunk.css
/opt/sqliteweb/www/static/css/...
/opt/sqliteweb/www/static/css/main.a2731a96.chunk.css.map
/opt/sqliteweb/www/static/media
/opt/sqliteweb/www/static/media/materialdesignicons-webfont.f60b16c8.ttf
/opt/sqliteweb/www/static/media/...
/opt/sqliteweb/www/static/js
/opt/sqliteweb/www/static/js/runtime-main.8264950d.js.map
/opt/sqliteweb/www/static/js/...
/opt/sqliteweb/www/static/js/main.bbc95380.chunk.js.map
/opt/sqliteweb/www/logo512.png
/opt/sqliteweb/www/index.html
/opt/sqliteweb/www/asset-manifest.json
/opt/sqliteweb/www/logo192.png
/opt/sqliteweb/www/favicon.ico
/opt/sqliteweb/www/robots.txt
```
You can then access those files with your browser at this address: [https://web1.sqlitecloud.io:8443/](https://web1.sqlitecloud.io:8443/)

## Using the JSON API
To access the JSON API, call the required endpoint with the corresponding HTTP VERB.

#### Example:

```console
curl --silent --insecure https://web1.sqlitecloud.io:8443/api/v1/ping 
```

The should output something like:

```console
{ 
  ResponseID: 0,
  Status:  0,
  Message: "pong",
}
```


### Auth
Authentication to the REST API is done with the help of JWT tokens. JWT Tokens consist of 3 parts separated with a '.'. Every part is base64 encoded. Let's have a look at the following token: 

```console
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc3NjEwNTUsImp0a
SI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIj
oiYXBpL3YxLyJ9.j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM`
```
The first part is the header of the token. As you can see, it is: >eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9<. This base64 string contains the follwoing information:

```console
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" | base64 -D

{"alg":"HS256","typ":"JWT"}
```

The second part is called the **"Claims"**. It contains the following information:

```console
echo "eyJleHAiOjE2Mzc3NjEwNTUsImp0aSI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIjoiYXBpL3YxLyJ9" | base64 -D

{"exp":1637761055,"jti":"1405","iat":1637760755,"nbf":1637760755,"sub":"api/v1/"}
```

The third part (`j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM`) is the cryptographic signature over the first and second part. The signature is salted by the secret jwt_key string (see config file).

### Authenticate / Loggin in
To authenticate to the server, call the authentication provider like this:

#### Example:

```console
curl --silent --insecure -X POST https://web1.sqlitecloud.io:8443/api/v1/auth -H 'Content-Type: application/json; charset=utf-8' -d '{"RequestID":1405,"Login":"admin","Password":"foo"}'
```

The result should look like this:

```console
{"ResponseID":1405,"Status":0,"Message":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc3NjEwNTUsImp0aSI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIjoiYXBpL3YxLyJ9.j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM"}
```

### Refreshing the Token
Calling this auth endpoint multiple times will create new JWT tokens and at the same time, will invalidate all the previous tokens for this user.

This behavior can be used to refresh a JWT token that is about to expire. However, for this case, another (easier) way exists. You can send the old JWT Token to the authentication provider endpoint (without username and password) and receive a fresh token in return (the old token is invalid after this call).

#### Example:

```console
curl --silent 
     --insecure 
     -X POST https://web1.sqlitecloud.io:8443/api/v1/auth 
     -H 'Content-Type: application/json; charset=utf-8' 
     -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc3NjEwNTUsImp0aSI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIjoiYXBpL3YxLyJ9.j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM'
```
Will give you a new JWT Token like this for Example:

```console
{"ResponseID":1405,"Status":0,"Message":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc3NjEwNTUsImp0aSI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIjoiYXBpL3YxLyJ9.j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM"}
```

### UnAuthenticate / Loging out
To render the actual JWT Token as invalid, just call the authentication provider with the DELETE HTTP Verb. This operation is equivalent with immediately logging out of the service.

#### Example:

```console
curl --silent 
     --insecure 
     -X DELETE https://web1.sqlitecloud.io:8443/api/v1/auth 
     -H 'Content-Type: application/json; charset=utf-8' 
     -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzc3NjEwNTUsImp0aSI6IjE0MDUiLCJpYXQiOjE2Mzc3NjA3NTUsIm5iZiI6MTYzNzc2MDc1NSwic3ViIjoiYXBpL3YxLyJ9.j4ECkdbLPzLnB76H5NK9X4cH4SGp-m7FYLfFApOwovM'
```

If everything went fine, the server will respond like this:

```console
{"ResponseID":0,"Status":0,"Message":"OK"}
```


### Dummy request
To speed up the development, the feature of creating dummy REST request/response pairs and performing calls against those dummy endpoints has been added to the server. This way, new endpoints can be developed and tested and can be used in the front-end code just right from the beginning.

To set up a new dummy endpoint, create a directory path under the api.path (see config file) like so:

```console
mkdir -p /opt/sqliteweb/api/v1/ping
```

This directory path maps the URL endpoint path 1:1.

Now, you can create a response for a specific HTTP verb like this:

```console
touch /opt/sqliteweb/api/v1/ping/GET.json
```
Please note, that the filename MUST follow the following scheme: < VERB >.json

Those HTTP Verbs are supported:

- GET
- HEAD
- POST
- PUT
- DELETE
- PATCH

Finally, you can specify the contents of your dummy response like this:

```console
echo "{                 " >> /opt/sqliteweb/api/v1/ping/GET.json
echo "  ResponseID: 0,  " >> /opt/sqliteweb/api/v1/ping/GET.json
echo "  Status:  0,     " >> /opt/sqliteweb/api/v1/ping/GET.json
echo "  Message: "pong"," >> /opt/sqliteweb/api/v1/ping/GET.json
echo "}                 " >> /opt/sqliteweb/api/v1/ping/GET.json
```

However, it is strongly recommended, that you use the editor of your choice or upload this file from your local machine.

**Please note: A dummy endpoint does not evaluate any dynamic input data from the request - whatsoever.**

## Using Lua to write endpoints
Lua version 5.2 is built right into the executable of sqliteweb. No external dll's or shared object files are required.

Using Lua endpoints is straight forward and pretty much the same as using dummy.json files. To use Lua endpoints, you first have to create the corresponding folder structure, beginning with: `/api/v1/`. Then add the desired endpoint name like `luatest/` followed by one of the supported Verbs (see above) and an `.lua`. A Complete Lua endpoint path could look like this for example:

```console
/api/v1/luatest/GET.lua
```

To execute this endpoint, point your browser or REST client to:

````console
https://web1.sqlitecloud.io:8443/api/v1/luatest/
````
In real life, you would most probably want to use some variable components in your endpoint path, like this for example:

````console
https://web1.sqlitecloud.io:8443/api/v1/node/{nodeID}
````
To implement this, you have to create the folder structure like this:

````console
/api/v1/node/{0}
````

You can add as many variable path components into your endpoint as you want, this could look like this for example:


````console
/api/v1/stat/{0}/{1}/
````

And then call it like this:

````console
https://web1.sqlitecloud.io:8443/api/v1/stat/cpu/1/
````

Those variable path components must be numbered like this {0}, {1}...{n}. No number must be used twice an a path. The variable path components can be of any data type (no type check is done). The variable endpoint path components are passed to the lua script in a global `args[0,1...n]` array. It is up to the programmer to do the type checks in lua.

This simple lua program demonstrate how to work with the variable path components:

````lua
print( args[ 0 ] ) --- please note, args[0] is always the script name
print( args[ 1 ] ) --- args[ 1 ] is what was {0} in the path 
print( args[ 2 ] ) --- this means, everything is shifted up by one position
````

It is very easy to create valid REST/JSON packages with Lua. To do this, the following commands where added to the Language:

````lua
SetHeader( "HeaderName", "HeaderValue" )
SetStatus( statusCode )
Write( data )
````

To write a very basic "Hello SQLiteCloud program, the following code is necessary:"

````lua
SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )
Write( '{ "Message": "Hello SQLite Cloud"}' )
SetStatus( 200 )
````

However, it is not very useful to create JSON packages like this. To create JSON packages, a JSON Lua package can be used. Those packages can be stored anywhere. The default path is: `/opt/sqliteweb/lib/`. To change this path, a new section in `/etc/sqliteweb/sqliteweb.ini` was added:

````console
[lua]
  package.path  = /opt/sqliteweb/lib/?.lua
````
In this folder, you can find a file called: `json.lua`
To use it in your Lua endpoint, write something like this:

````lua
json = require "json"
````

From now on, working with JSON is very easy. To parse the JSON data that was send - for example - in the body of a POST request, you can write:

````lua
request = json.decode( body )
````
Please note: `body` is a special predefined variable that was added to the Lua language to help and make it very easy for the programmer to access the body data of a request. It is so to say the opposite of the Write() function call.

A more complete Lua Example, that demonstrates JSON Packet parsing as well as the creation, makes use of variable path components and is doing a SQL SELECT in one of the nodes could look like this:

````lua
json = require "json"

SetHeader( "Content-Type", "application/json" )
SetHeader( "Content-Encoding", "utf-8" )

request = json.decode( body )

result = queryNode( "SELECT * FROM Dummy" )

Response = {
  Request = request,
  Parameter = {
    First  = args[ 1 ],
    Second = args[ 2 ]
  },

  ResponseID = request[ 'RequestID' ],
  Status = 0,
  Message = "OK",

  QueryResult = result
}

Write( json.encode( Response ) )
SetStatus( 200 )
````
