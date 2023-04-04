# C SDK

SQCloud is the C application programmer's interface to SQLite Cloud. SQCloud is a set of library functions that allow client programs to pass queries and SQL commands to the SQLite Cloud backend server and to receive the results of these queries. In addition to the standard SQLite statements, several other [commands](https://docs.sqlitecloud.io/docs/commands) are supported.

The following files are required when compiling a C application:

- sqcloud.c/.h
- lz4.c/.h
- libtls.a (for static linking) or libtls.so/.dylib/.dll (for dynamic linking)

The header file `sqcloud.h` must be included in your C application.

All the communications between the client and the server are encrypted and so, you are required to link the LibreSSL (libtls) library with your client.

### Install LibreSSL

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



### Use the C SDK in your project

Prerequisites:

- [Install LibreSSL](#install-libressl)

Use SQCloud:

- Include the header file `sqcloud.h`

- Compile `sqcloud.c ` and `lz4.c`. Examples from the Makefile: 

  ```
  lz4.o:	lz4.c lz4.h
  	$(CC) $(OPTIONS) $(INCLUDES) lz4.c -c -o lz4.o
  	
  sqcloud.o:	sqcloud.c sqcloud.h
  	$(CC) $(OPTIONS) $(INCLUDES) sqcloud.c -c -o sqcloud.o
  
  libsqcloud.a: lz4.o sqcloud.o
  	$(AR) rcs libsqcloud.a *.o
  ```

- Link with the compiled sources 

- Link with the libtls shared library using the following gcc option flags: 

  - Add the libressl dir to the list of directories to be searched for: -L<path-to-libressl-lib>. Examples: 

    linux: `-L/usr/local/libressl/lib/` 

    macOS: `-L/opt/homebrew/opt/libressl/lib`

  - make the library available at runtime with the environment variable LD_LIBRARY_PATH or rpath. Example: `-Wl,-rpath=/usr/local/libressl/lib/` 

  - link the shared library: `-ltls`

  Example from the Makefile:

  ```
  CC        := gcc
  INCLUDES  := -I. -Icli 
  OPTIONS   := -Wno-macro-redefined -Wno-shift-negative-value -Os
  LDFLAGS = -ltls
  LIBFLAGS = -L/usr/local/libressl/lib/ -Wl,-rpath=/usr/local/libressl/lib/
  ...
  cli: libsqcloud.a cli/linenoise.c cli/linenoise.h cli/main.c
  	$(CC) $(OPTIONS) $(INCLUDES) cli/*.c libsqcloud.a -o sqlitecloud-cli ${LIBFLAGS} ${LDFLAGS}
  ```



# SQLite Cloud CLI

The command-line interface is a text-based user interface to interact with a SQLite Cloud server.

Prerequisites:

- [Install LibreSSL](#install-libressl)

Build:

```
make cli
```

Connect to a SQLite Cloud server:

```
./sqlitecloud-cli -h <host>.sqlite.cloud -p 8860 -n admin -m <password>
```
