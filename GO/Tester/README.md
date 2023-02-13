# RUN TEST WITH MAKE
### Cluster 
`make test-cluster`

### Nocluster
`PORT=8850 make test-nocluster`

# RUN TESTS WITH COMMANDS
### Run all the tests in the default `scripts` folder on the dev1 server with user `admin` and password `admin`:
`go test -v -connstring=sqlitecloud://admin:admin@dev1.sqlitecloud.io?tls=SQLiteCloudCA`

### Run a specific test on localhost, with debug mode enabled
`go test -v -path=scripts/nwriters1.test -debug -connstring=sqlitecloud://admin:admin@localhost:8860?tls=SQLiteCloudCA`

### no tls
`go test -v -path=scripts/nwriters1.test -debug -connstring=sqlitecloud://localhost?tls=no`
