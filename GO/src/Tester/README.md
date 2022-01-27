# ENVIRONMENT
go env -w GO111MODULE=off
export GOPATH=/Users/andrea/Documents/GitHub/SQLiteCloud/sdk/GO/

# RUN TESTS
### Run all the tests in the default `scripts` folder on the dev1 server with user `admin` and password `admin`:
`go test -v -connstring=sqlitecloud://admin:admin@dev1.sqlitecloud.io`

### Run a specific test on localhost, with debug mode enabled
`go test -v -path=scripts/nwriters1.test -debug -connstring=sqlitecloud://admin:admin@localhost:8860`

### no tls
`go test -v -path=scripts/nwriters1.test -debug -connstring=sqlitecloud://localhost?tls=no`
