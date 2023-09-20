AR        := ar
CC        := gcc
INCLUDES  := -I. -Icli 
OPTIONS   := -Wno-shift-negative-value -Os -fPIC


ifeq ($(OS),Windows_NT)
	# Windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		# MacOS
		LDFLAGS = -lpthread -ltls
		ifeq ($(TLS_STATIC), 1)
			LIBFLAGS = -L./SSL/macos_fat
		else
			LIBRESSLDIR ?= /opt/homebrew/opt/libressl/lib/
			LIBFLAGS = -L$(LIBRESSLDIR)
		endif
	else
		# Linux
		LDFLAGS = -lpthread -ltls
		ifeq ($(TLS_STATIC), 1)
			LIBFLAGS = -L./SSL/linux_64bit
		else
			LIBRESSLDIR ?= /usr/local/lib/
			# LIBRESSLDIR ?= /usr/local/libressl/lib/
			LIBFLAGS = -L$(LIBRESSLDIR) -Wl,-rpath=$(LIBRESSLDIR)
		endif
	endif
endif

all: libsqcloud.a

lz4.o:	lz4.c lz4.h
	$(CC) $(OPTIONS) $(INCLUDES) lz4.c -c -o lz4.o
	
sqcloud.o:	sqcloud.c sqcloud.h
	$(CC) $(OPTIONS) $(INCLUDES) sqcloud.c -c -o sqcloud.o

libsqcloud.a: lz4.o sqcloud.o
	$(AR) rcs libsqcloud.a *.o

libsqcloud.so: lz4.o sqcloud.o
	$(CC) -fPIC -shared *.o -o libsqcloud.so ${LIBFLAGS} ${LDFLAGS}
	
libsqcloud.dylib: lz4.o sqcloud.o
	$(CC) -dynamiclib -install_name -flat_namespace *.o -o libsqcloud.dylib

cli: sqlitecloud-cli

sqlitecloud-cli: libsqcloud.a cli/linenoise.c cli/linenoise.h cli/main.c
	$(CC) $(OPTIONS) $(INCLUDES) cli/*.c libsqcloud.a -o sqlitecloud-cli ${LIBFLAGS} ${LDFLAGS}

clean:
	rm -rf *.o *.a *.so *.dylib sqlitecloud-cli
	
test: cli
	./sqlitecloud-cli -h localhost
