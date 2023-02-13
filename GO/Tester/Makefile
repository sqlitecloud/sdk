HOST ?= localhost
PORT ?= 8860
ADMINPASSWD ?= admin
TLS ?= SQLiteCloudCA

test-cluster:
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT)?tls=$(TLS) -path=scripts

test-nocluster:
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT)?tls=$(TLS) -path=scripts/general
	go test -v -connstring=sqlitecloud://admin:$(ADMINPASSWD)@$(HOST):$(PORT)?tls=$(TLS) -path=scripts/singlenode-or-leader