GOPATH		= $(shell go env GOPATH)

# Test SDK
.PHONY: test
test:	
	cd sdk/test; go test -v

# GO SDK
.PHONY: sdk
sdk:	sdk/*.go
	cd sdk; go build

# CLI App
$(GOPATH)/bin/sqlc:	sdk/*.go cli/sqlc.go
	cd cli; go build -o $(GOPATH)/bin/sqlc
	
cli: $(GOPATH)/bin/sqlc

github:
	open https://github.com/sqlitecloud/sdk 
	
diff:
	git difftool
	

# gosec
gosec:
ifeq ($(wildcard $(GOPATH)/bin/gosec),)
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
endif

checksec:	gosec sdk/*.go cli/*.go
	cd sdk; $(GOPATH)/bin/gosec -exclude=G304 .
	cd cli; $(GOPATH)/bin/gosec -exclude=G304,G302 .


# Documentation
godoc:
ifeq ($(wildcard $(GOPATH)/bin/godoc),)
	go install golang.org/x/tools/cmd/godoc
endif

doc:	godoc
ifeq ($(wildcard ./src),)
	ln -s . src			
endif
	@echo "Hit CRTL-C to stop the documentation server..."
	@( sleep 1 && open http://localhost:6060/pkg/sdk/ ) &
	@$(GOPATH)/bin/godoc -http=:6060 -index -play -goroot ./

clean:
	rm -rf $(GOPATH)/bin/sqlc* 

all: sdk cli test_dev1
