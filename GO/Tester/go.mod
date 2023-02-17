module github.com/sqlitecloud/sdk/go/tester

go 1.18

require github.com/PaesslerAG/gval v1.2.1

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/xo/dburl v0.12.4 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.1.0 // indirect
)

require github.com/sqlitecloud/sdk v0.0.0

replace github.com/sqlitecloud/sdk v0.0.0 => ../sdk