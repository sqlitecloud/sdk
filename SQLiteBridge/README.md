# SQLite Bridge

The SQLite Bridge is ANSI-C source code. It implements all the official SQLite's public API. The SQLite Bridge can be used as a drop-in replacement for the original SQLite amalgamation to use the SQLite Cloud service instead of the local SQLite file without any change to you application code, just replace the SQLite filename with a SQLite Cloud connection string. 

All the communications between the client and the server are encrypted and so, you are required to link the LibreSSL (libtls) library with your client.

## Install LibreSSL

##### Linux:

- download the latest portable source from [www.libressl.org](https://www.libressl.org/)

- extract the tarball and change the directory to the libressl dir

- compile and install LibreSSL. By default, the install script will install LibreSSL to the `/usr/local/` folder. In order to avoid issue with other SSL libraries installed on the system, you can specify a different install directory, for example `/usr/local/libressl`, with the following command:

  ```
  ./configure --prefix=/usr/local/libressl --with-openssldir=/usr/local/libressl && make && make install
  ```

##### macOS:

```
brew install libressl
```



## Use the C SDK in your project

Just replace the SQLite amalgamation files with the SQLite Bridge files (sqlite.h and sqlite.c) and link your project with the libtls shared library using the following gcc option flags: 

- Add the libressl dir to the list of directories to be searched for: -L<path-to-libressl-lib>. Examples: 

  linux: `-L/usr/local/libressl/lib/` 

  macOS: `-L/opt/homebrew/opt/libressl/lib`

- make the library available at runtime with the environment variable LD_LIBRARY_PATH or rpath. Example: `-Wl,-rpath=/usr/local/libressl/lib/` 

- link the shared library: `-ltls`

#### Example: Compiling The SQLite Command-Line Interface with the SQLite Bridge

Prerequisites:

- [Install LibreSSL](#install-libressl)

 A build of the [SQLite command-line interface](https://www.sqlite.org/cli.html) requires three source files:

- sqlite3.c: The SQLite Bridge amalgamation source file
- sqlite3.h: The SQLite Bridge amalgamation header file
- shell.c: The command-line interface program from the [SQLite amalgamation tarball](https://www.sqlite.org/download.html#amalgtarball).

To build the CLI, simply put these three files in the same directory and compile them together. For example, if LibreSSL was installed in the `/usr/local/libressl` folder you can use the following command:

```
gcc shell.c sqlite3.c -I. -L/usr/local/libressl/lib/ -Wl,-rpath=/usr/local/libressl/lib/ -ltls -o sqlitecloud
```

Then, connect to a database on your SQLite Cloud server with something like this:

```
./sqlitecloud "sqlitecloud://<username>:<password>@<hostname>.sqlite.cloud:8860/<database-name>?create=1&sqlite=1"
```

