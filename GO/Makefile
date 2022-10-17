CURL 			:= curl --silent --insecure -X
TEST_URL	:= https://web1.sqlitecloud.io:8443/dashboard/v1


# Pre-requirements
install-prerequirements:	liner term docopt-go dburl lz4 godoc gosec gorilla ini kardianos golang-jwt golua glob

golua:
ifeq ($(wildcard src/github.com/Shopify/go-lua),)
	-go get -v github.com/Shopify/go-lua
endif

glob:
ifeq ($(wildcard src/github.com/gobwas/glob),)
	-go get -v github.com/gobwas/glob
endif

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


# gorilla
gorilla:
ifeq ($(wildcard src/github.com/gorilla),)
	go get github.com/gorilla/websocket
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/handlers
endif

# ini
ini:
ifeq ($(wildcard src/gopkg.in/ini.v1),)
	git config --global http.sslverify "false"
	go get gopkg.in/ini.v1
	git config --global http.sslverify "true"
endif

# kardianos
kardianos:
ifeq ($(wildcard src/github.com/kardianos),)
	go get github.com/kardianos/service
endif


# golang-jwt
golang-jwt:
ifeq ($(wildcard src/github.com/golang-jwt),)
	go get github.com/golang-jwt/jwt
endif

#lua
go-lua:
ifeq ($(wildcard src/github.com/Shopify/go-lua),)
	go get github.com/Shopify/go-lua
endif

# Test Programs
bin/test-selectarray:	sdk src/test-selectarray
	go build -o bin/test-selectarray src/Test/test-selectarray.go

bin/tls-test1:	sdk src/Test/tls-test1.go
	go build -o bin/tls-test1 src/Test/tls-test1.go

bin/tls-server:	sdk src/Test/tls-server.go
	go build -o bin/tls-server src/Test/tls-server.go
	
bin/tls-dump1:	sdk src/Test/tls-dump1.go
	go build -o bin/tls-dump1 src/Test/tls-dump1.go
	
bin/tls-dump2:	sdk src/Test/tls-dump2.go
	go build -o bin/tls-dump2 src/Test/tls-dump2.go
	
bin/tls-dump3:	sdk src/Test/tls-dump3.go
	go build -o bin/tls-dump3 src/Test/tls-dump3.go

bin/tls-compress:	sdk src/Test/tls-compress.go
	go build -o bin/tls-compress src/Test/tls-compress.go
		
bin/tls-literals:	sdk src/Test/tls-literals.go
	go build -o bin/tls-literals src/Test/tls-literals.go

bin/tls-pubsub:	sdk src/Test/tls-pubsub.go
	go build -o bin/tls-pubsub src/Test/tls-pubsub.go

test:	bin/test-selectarray bin/tls-test1 bin/tls-server bin/tls-dump1 bin/tls-dump2 bin/tls-dump3 bin/tls-compress bin/tls-literals bin/tls-pubsub

# GO SDK
sdk:	install-prerequirements src/sqlitecloud/*.go


# CLI App
bin/sqlc:	sdk src/CLI/sqlc.go
	go build -o bin/sqlc src/CLI/sqlc.go
	
cli: bin/sqlc


# SqliteWeb
bin/sqliteweb_osx: src/sqliteweb/*.go
	go build -o bin/sqliteweb_osx src/sqliteweb/*.go
	
bin/sqliteweb_linux: src/sqliteweb/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/sqliteweb_linux src/sqliteweb/*.go
	
bin/sqliteweb_win: src/sqliteweb/*.go
	GOOS=windows go build -o bin/sqliteweb_win src/sqliteweb/*.go

id_rsa.pub:
ifeq ($(wildcard  ~/.ssh/id_rsa.pub),)
	ssh-keygen -t rsa
endif

web_install: id_rsa.pub
	ssh-copy-id ***REMOVED***
	
web_update: bin/sqliteweb_linux
	ssh ***REMOVED*** '/etc/init.d/sqliteweb stop'
	scp bin/sqliteweb_linux web1.sqlitecloud.io:/opt/sqliteweb/sbin/.
	ssh ***REMOVED*** '/etc/init.d/sqliteweb start &'

web_update_service: bin/sqliteweb_linux
	ssh ***REMOVED*** 'systemctl stop sqliteweb'
	scp bin/sqliteweb_linux web1.sqlitecloud.io:/opt/sqliteweb/sbin/.
	ssh ***REMOVED*** 'systemctl start sqliteweb'	

web_stop:
	ssh ***REMOVED*** '/etc/init.d/sqliteweb stop'
	
web_start:
	ssh ***REMOVED*** '/etc/init.d/sqliteweb start'
	
web_restart:
	ssh ***REMOVED*** '/etc/init.d/sqliteweb restart'

web_stop_service:
	ssh ***REMOVED*** 'systemctl stop sqliteweb'

web_start_service:
	ssh ***REMOVED*** 'systemctl start sqliteweb'

web_restart_service:
	ssh ***REMOVED*** 'systemctl restart sqliteweb'

web_test:
	$(eval JWT := $(shell $(CURL) "POST" $(TEST_URL)/auth -H 'Content-Type: application/json; charset=utf-8' -d '{ "Login": "my.address@domain.com ", "Password": "password" }' | grep -o 'message":"[^"]*' | cut -c 11-))
	@echo Auth...success: \(JWT = $(JWT)\)
	@$(CURL) "GET" $(TEST_URL)/ping -H 'Authorization: Bearer $(JWT)' | grep -q "PONG" && echo "Ping...success." || echo "Ping...ERROR."

sync:
	rsync -a src/sqliteweb/dashboard/v1/*  ***REMOVED***:/opt/sqliteweb/dashboard/v1/.
	rsync -a src/sqliteweb/lib/* ***REMOVED***:/opt/sqliteweb/lib/.

sync_admin:
	rsync -a src/sqliteweb/admin/v1/* ***REMOVED***:/opt/sqliteweb/admin/v1/.

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
	rm -rf rm -rf src/github.com src/golang.org/ src/honnef.co src/mvdan.cc src/gopkg.in
	rm -rf bin/tls-* bin/dump* bin/selc* 

all: sdk cli test

# Go Tester for sqlitecloud

install-tester-prerequirements: paesslerag-gval

paesslerag-gval:
ifeq ($(wildcard src/github.com/PaesslerAG/gval),)
	-go get -v github.com/PaesslerAG/gval
endif

test_gotester_leader_run: install-prerequirements install-tester-prerequirements
	go test src/Tester/tester_test.go -connstring=sqlitecloud://admin:admin@localhost:8860 -path=scripts -v
