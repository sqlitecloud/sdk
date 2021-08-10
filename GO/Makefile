# Pre-requirements
install-prerequirements:	patch docopt-go dburl lz4 godoc

# Linenoise
go-cli:
ifeq ($(wildcard src/github.com/deadsy),)
	go get github.com/deadsy/go-cli
endif

# Command line arguments
docopt-go:
ifeq ($(wildcard src/github.com/docopt),)
	go get github.com/docopt/docopt-go
endif

# Connection strings
dburl:
ifeq ($(wildcard src/github.com/xo),)
	go get github.com/xo/dburl
endif

# LZ4 compression
lz4:
ifeq ($(wildcard src/github.com/pierrec),)
	go get github.com/pierrec/lz4
endif

# Godoc
godoc:
ifeq ($(wildcard bin/godoc),)
	go get -u golang.org/x/tools/cmd/godoc
endif

# Linenoise patching for OSX
src/github.com/deadsy/go-cli/line.go:	go-cli src/deadsy_line.patch
	cp src/deadsy_line.patch src/github.com/deadsy/go-cli/line.go
	
patch: src/github.com/deadsy/go-cli/line.go

# Test Programs
src/Test/test1:	csdk src/Test/test1.go
	go build -o src/Test/test1 src/Test/test1.go

test:	sdk src/Test/test1

# CSDK
src/sqlitecloud/libsqcloud.a:	
	make -C ../C libsqcloud.a
	cp ../C/libsqcloud.a  src/sqlitecloud/libsqcloud.a

csdk:	src/sqlitecloud/libsqcloud.a


# GO SDK
sdk:	proxy src/sqlitecloud/*.go


# CLI App
bin/sqlc:	sdk src/CLI/sqlc.go
	go build -o bin/sqlc src/CLI/sqlc.go
	
cli: bin/sqlc

github:
	open https://github.com/sqlitecloud/sdk 
	
diff:
	git difftool

# Documentation
doc:	godoc
	@echo "Hit CRTL-C to stop the documentation server..."
	@( sleep 1 && open http://localhost:6060/pkg/sqlitecloud/ ) &
	@./bin/godoc -http=:6060 -index -play -goroot ./src/sqlitecloud

src/sqlitecloud/synthesized-c-proxy.go: ../C/lz4.h ../C/sqcloud_private.h ../C/sqcloud.h ../C/lz4.c ../C/sqcloud.c
	echo "package sqlitecloud" > src/sqlitecloud/synthesized-c-proxy.go
	echo "" >> src/sqlitecloud/synthesized-c-proxy.go
	echo "// #cgo CFLAGS: -Wno-multichar -Wno-macro-redefined " >> src/sqlitecloud/synthesized-c-proxy.go
	echo "// #cgo LDFLAGS: -L. " >> src/sqlitecloud/synthesized-c-proxy.go
	echo "// #include <stdlib.h>" >> src/sqlitecloud/synthesized-c-proxy.go
  
	sed '/#include "lz4.h"/q' ../C/lz4.c | sed '/#include "lz4.h"/d' | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
	cat ../C/lz4.h | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
	sed -n -e '/#include "lz4.h"/,$$p' ../C/lz4.c | tail -n +2 | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
  
	sed '/#include "lz4.h"/q' ../C/sqcloud.c | sed '/#include "lz4.h"/d' | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
	cat ../C/sqcloud.h | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
	sed -n -e '/#include "sqcloud.h"/,$$p' ../C/sqcloud.c | tail -n +2 | sed 's/^/\/\/ /' >> src/sqlitecloud/synthesized-c-proxy.go
  
	echo "import \"C\"" >> src/sqlitecloud/synthesized-c-proxy.go
	echo "import \"unsafe\"" >> src/sqlitecloud/synthesized-c-proxy.go
  
	cat src/sqlitecloud/c_proxy.stubs >> src/sqlitecloud/synthesized-c-proxy.go

proxy: src/sqlitecloud/synthesized-c-proxy.go
	
clean:
	rm -rf src/Test/test1  src/sqlitecloud/libsqcloud.a
	# rm src/SQLiteCloud/c.go


all: sdk cli test
