HOST ?= localhost
PORT ?= 8860
ADMINPASSWD ?= admin

test-cluster:
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT) -path=scripts

test-nocluster:
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT) -path=scripts/general
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT) -path=scripts/singlenode-or-leader