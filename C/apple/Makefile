SDKROOT   ?= /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX13.3.sdk
TARGET    ?= macos
LIBTLS     = ./fat/libressl/$(TARGET)
LIBTLS_LIB = $(LIBTLS)/lib/libtls.a
LIBTLS_INC = $(LIBTLS)/include
AR        := ar
CC        ?= /usr/bin/clang -target arm64-apple-darwin -isysroot $(SDKROOT)
INCLUDES   = -I../ -I$(LIBTLS_INC)
OPTIONS   := -Wno-macro-redefined -Wno-shift-negative-value -O2 -pipe -no-cpp-precomp -Wall -fno-strict-aliasing  -fno-strict-overflow -fstack-protector-strong  -Qunused-arguments -Wno-pointer-sign
CFLAGS    ?= -arch arm64
LDFLAGS   ;= -lpthread
LIBFLAGS  ?=

all:	libsqcloud.a

lz4.o:	../lz4.c ../lz4.h
	$(CC) $(OPTIONS) $(CFLAGS) $(INCLUDES) ../lz4.c -c -o lz4.o
	
sqcloud.o:	../sqcloud.c ../sqcloud.h
	$(CC) $(OPTIONS) $(CFLAGS) $(INCLUDES) ../sqcloud.c -c -o sqcloud.o

libsqcloud.a: lz4.o sqcloud.o
	$(AR) rcs libsqcloud-tmp.a *.o
	libtool -static -o libsqcloud.a libsqcloud-tmp.a $(LIBTLS_LIB)

cli: libsqcloud.a ../cli/linenoise.c ../cli/linenoise.h ../cli/main.c
	$(CC) $(OPTIONS) $(CFLAGS) $(INCLUDES) ../cli/*.c libsqcloud.a -o sqlitecloud-cli ${LIBFLAGS} ${LDFLAGS}

clean:
	rm -rf *.o *.a *.so *.dylib sqlitecloud-cli
	
buildclean:
	rm -rf build Fat output
	
.PHONY: all cli clean buildclean
