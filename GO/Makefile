# Pre-requirements
install-prerequirements:	liner term docopt-go dburl lz4 godoc gosec

# liner
liner:
ifeq ($(wildcard src/github.com/peterh),)
	-go get -v github.com/peterh/liner
endif

# term
term:
ifeq ($(wildcard src/golang.org/x/term),)
	-go get -v golang.org/x/term
endif

# Command line arguments
docopt-go:
ifeq ($(wildcard src/github.com/docopt),)
	go get -v github.com/docopt/docopt-go
endif

# Connection strings
dburl:
ifeq ($(wildcard src/github.com/xo),)
	go get -v github.com/xo/dburl
endif

# LZ4 compression
lz4:
ifeq ($(wildcard src/github.com/pierrec),)
	go get -v github.com/pierrec/lz4
endif

# Godoc
godoc:
ifeq ($(wildcard bin/godoc),)
	go get -u golang.org/x/tools/cmd/godoc
endif

# gosec
gosec:
ifeq ($(wildcard bin/gosec),)
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh
endif

# Test Programs
bin/test1:	sdk src/Test/test1.go
	go build -o bin/test1 src/Test/test1.go

bin/server:	sdk src/Test/server.go
	go build -o bin/server src/Test/server.go
	
bin/dump1:	sdk src/Test/dump1.go
	go build -o bin/dump1 src/Test/dump1.go
	
bin/dump2:	sdk src/Test/dump2.go
	go build -o bin/dump2 src/Test/dump2.go
	
bin/dump3:	sdk src/Test/dump3.go
	go build -o bin/dump3 src/Test/dump3.go

test:	sdk bin/test1 bin/server bin/dump1 bin/dump2 bin/dump3

# GO SDK
sdk:	install-prerequirements src/sqlitecloud/*.go


# CLI App
bin/sqlc:	sdk src/CLI/sqlc.go
	go build -o bin/sqlc src/CLI/sqlc.go
	
cli: bin/sqlc

github:
	open https://github.com/sqlitecloud/sdk 
	
diff:
	git difftool
	
checksec:	gosec src/sqlitecloud/*.go
	./bin/gosec ./src/sqlitecloud/
	# ./bin/gosec ./src/CLI/

# Documentation
doc:	godoc
	@echo "Hit CRTL-C to stop the documentation server..."
	@( sleep 1 && open http://localhost:6060/pkg/sqlitecloud/ ) &
	@./bin/godoc -http=:6060 -index -play -goroot ./src/sqlitecloud

clean:
	rm -rf rm -rf src/github.com src/golang.org/
	rm -rf src/Test/test1

all: sdk cli test