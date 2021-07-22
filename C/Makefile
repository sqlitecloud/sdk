AR        := ar
CC        := gcc
INCLUDES  := -I. -Icli
OPTIONS   := -Wno-macro-redefined -Wno-shift-negative-value -Os

lz4.o:	lz4.c lz4.h
	$(CC) $(OPTIONS) $(INCLUDES) lz4.c -c -o lz4.o
	
sqcloud.o:	sqcloud.c sqcloud.h sqcloud_private.h
	$(CC) $(OPTIONS) $(INCLUDES) sqcloud.c -c -o sqcloud.o

libsqcloud.a: lz4.o sqcloud.o
	$(AR) rcs libsqcloud.a *.o

libsqcloud.so: lz4.o sqcloud.o
	$(CC) -fPIC -shared *.o -o libsqcloud.so
	
libsqcloud.dylib: lz4.o sqcloud.o
	$(CC) -dynamiclib -install_name -flat_namespace *.o -o libsqcloud.dylib

all:	lz4.o

cli:	libsqcloud.a cli/linenoise.c cli/linenoise.h cli/main.c
	$(CC) $(OPTIONS) $(INCLUDES) cli/*.c libsqcloud.a -o sqlitecloud-cli

clean:
	rm -rf *.o *.a *.so *.dylib
	
test: cli
	./sqlitecloud-cli -h ***REMOVED***