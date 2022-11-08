module github.com/sqlitecloud/sqliteweb

go 1.18

require (
	github.com/Shopify/go-lua v0.0.0-20221004153744-91867de107cf
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/gobwas/glob v0.2.3
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/kardianos/service v1.2.2
	gopkg.in/ini.v1 v1.67.0
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/xo/dburl v0.12.4 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
)

require (
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/sqlitecloud/sdk v0.0.0
)

replace github.com/sqlitecloud/sdk v0.0.0 => ../sdk
