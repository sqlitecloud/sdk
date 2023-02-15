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
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/swaggest/jsonschema-go v0.3.42 // indirect
	github.com/swaggest/refl v1.1.0 // indirect
	github.com/xo/dburl v0.12.4 // indirect
	golang.org/x/net v0.0.0-20220921155015-db77216a4ee9 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
	golang.org/x/time v0.0.0-20220922220347-f3bd1da661af // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/digitalocean/godo v1.95.0
	github.com/felixge/httpsnoop v1.0.3
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/uuid v1.3.0
	github.com/sendgrid/sendgrid-go v3.12.0+incompatible
	github.com/sqlitecloud/sdk v0.0.0
	github.com/swaggest/openapi-go v0.2.26
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	golang.org/x/exp v0.0.0-20221230185412-738e83a70c30
	golang.org/x/text v0.3.7
)

replace github.com/sqlitecloud/sdk v0.0.0 => ../sdk
