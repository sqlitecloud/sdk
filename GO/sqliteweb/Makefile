CURL 			:= curl --silent --insecure -X
TEST_URL		:= https://web1.sqlitecloud.io:8443/dashboard/v1
TEST_USER		?= my.address@domain.com
TEST_PASSWORD 	?= password
GOPATH		= $(shell go env GOPATH)

# SqliteWeb
$(GOPATH)/bin/sqliteweb_osx: ../sdk/*.go *.go
	GOOS=darwin GOARCH=arm64 go build -o $(GOPATH)/bin/sqliteweb_osx
	
$(GOPATH)/bin/sqliteweb_linux: ../sdk/*.go *.go
	GOOS=linux GOARCH=amd64 go build -o $(GOPATH)/bin/sqliteweb_linux *.go
	
$(GOPATH)/bin/sqliteweb_win: ../sdk/*.go sqliteweb/*.go
	GOOS=windows go build -o $(GOPATH)/bin/sqliteweb_win *.go

sqliteweb_osx: $(GOPATH)/bin/sqliteweb_osx
sqliteweb_linux: $(GOPATH)/bin/sqliteweb_linux
sqliteweb_win: $(GOPATH)/bin/sqliteweb_win

id_rsa.pub:
ifeq ($(wildcard  ~/.ssh/id_rsa.pub),)
	ssh-keygen -t rsa
endif

web_install: id_rsa.pub
	ssh-copy-id root@web1.sqlitecloud.io
	
web_update: sqliteweb_linux
	ssh root@web1.sqlitecloud.io '/etc/init.d/sqliteweb stop'
	scp $(GOPATH)/bin/sqliteweb_linux web1.sqlitecloud.io:/opt/sqliteweb/sbin/.
	ssh root@web1.sqlitecloud.io '/etc/init.d/sqliteweb start &'

web_update_service: sqliteweb_linux
	ssh root@web1.sqlitecloud.io 'systemctl stop sqliteweb'
	scp $(GOPATH)/bin/sqliteweb_linux web1.sqlitecloud.io:/opt/sqliteweb/sbin/.
	ssh root@web1.sqlitecloud.io 'systemctl start sqliteweb'	

web_stop:
	ssh root@web1.sqlitecloud.io '/etc/init.d/sqliteweb stop'
	
web_start:
	ssh root@web1.sqlitecloud.io '/etc/init.d/sqliteweb start'
	
web_restart:
	ssh root@web1.sqlitecloud.io '/etc/init.d/sqliteweb restart'

web_stop_service:
	ssh root@web1.sqlitecloud.io 'systemctl stop sqliteweb'

web_start_service:
	ssh root@web1.sqlitecloud.io 'systemctl start sqliteweb'

web_restart_service:
	ssh root@web1.sqlitecloud.io 'systemctl restart sqliteweb'

web_test:
	$(eval JWT := $(shell $(CURL) "POST" $(TEST_URL)/auth -H 'Content-Type: application/json; charset=utf-8' -d '{ "Login": "$(TEST_USER)", "Password": "$(TEST_PASSWORD)" }' | grep -o 'message":"[^"]*' | cut -c 11-))
	@echo Auth...success: \(JWT = $(JWT)\)
	@$(CURL) "GET" $(TEST_URL)/ping -H 'Authorization: Bearer $(JWT)' | grep -q "PONG" && echo "Ping...success." || echo "Ping...ERROR."

sync:
	rsync -a dashboard/v1/*  root@web1.sqlitecloud.io:/opt/sqliteweb/dashboard/v1/.
	rsync -a lib/* root@web1.sqlitecloud.io:/opt/sqliteweb/lib/.

sync_admin:
	rsync -a admin/v1/* root@web1.sqlitecloud.io:/opt/sqliteweb/admin/v1/.
