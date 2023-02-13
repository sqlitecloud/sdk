//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.1
//     //             ///   ///  ///    Date        : 2022/02/18
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

/*
// To run, enter:
cd sdk/GO/sqliteweb/
go run *.go
go run *.go --config etc/sqliteweb/sqliteweb.ini

// Compile with:
GOOS=linux go build -o sqliteweb  *.go
*/

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/kardianos/service"
	"gopkg.in/ini.v1"
) //import "io"

//import "time"
//import "errors"

//import "strconv"

//import "github.com/gorilla/mux"
//import "github.com/gorilla/websocket"

var app_name = "sqliteweb"
var long_name = "SQLite Cloud Web Server"
var version = "version 0.1.0"
var copyright = "(c) 2022 by SQLite Cloud Inc."
var service_name = "web.sqlitecloud.io"
var jwt_issuer = "web.sqlitecloud.io"

var cfg *ini.File

func main() {
	// Read command line arguments
	if p, err := docopt.ParseArgs(strings.ReplaceAll(usage, "::", ": "), nil, fmt.Sprintf("%s %s, %s", app_name, version, copyright)); err != nil {
		panic(err)
	} else {

		// Look for a custom config file
		configfile := "/etc/sqliteweb/sqliteweb.ini"
		for key, a := range p {
			switch key = strings.TrimSpace(strings.ToLower(key)); {
			case a == nil:
				continue
			case key == "--config":
				configfile = fmt.Sprintf("%v", a)
			}
		}

		// Read the .ini file -> https://github.com/go-ini/ini
		var err error
		if cfg, err = ini.Load(configfile); err != nil {
			fmt.Printf("Fail to read file: %v", err)
			os.Exit(1)
		} else {

			cfg.Section("dashboard").Key("modified").SetValue("hallo")
			if file, err := os.Stat(configfile); err == nil {
				cfg.Section("dashboard").Key("modified").SetValue(file.ModTime().Format("2006-01-02 15:04:05"))
			}

			// Overload the .ini file with the command line arguments
			for key, a := range p {
				if a != nil {
					value := fmt.Sprintf("%v", a)
					fmt.Printf("%s = %v\r\n", key, a)

					switch strings.TrimSpace(strings.ToLower(key)) {
					case "--address":
						cfg.Section("server").Key("address").SetValue(value)
					case "--port":
						cfg.Section("server").Key("port").SetValue(value)
					case "--cert":
						cfg.Section("server").Key("cert_chain").SetValue(value)
					case "--key":
						cfg.Section("server").Key("cert_key").SetValue(value)

					case "--www":
						cfg.Section("www").Key("path").SetValue(value)
					case "--stubs":
						cfg.Section("stubs").Key("path").SetValue(value)
					default: // fmt.Printf( "%10s = %v\r\n", key, a )
					}
				}
			}

			// start the service -> https://github.com/kardianos/service
			svcConfig := &service.Config{
				Name:        app_name,                                                // "GoServiceExampleSimple",
				DisplayName: long_name,                                               // "Go Service Example",
				Description: fmt.Sprintf("%s %s, %s", long_name, version, copyright), // "This is an example Go service.",
			}

			initializeSQLiteWeb()
			SQLiteWeb.Address = cfg.Section("server").Key("address").String()
			SQLiteWeb.Port = cfg.Section("server").Key("port").RangeInt(8433, 0, 0xFFFF)

			SQLiteWeb.CertPath = cfg.Section("server").Key("cert_chain").String()
			SQLiteWeb.KeyPath = cfg.Section("server").Key("cert_key").String()

			SQLiteWeb.logLevel = cfg.Section("server").Key("loglevel").String()
			SQLiteWeb.logFile = cfg.Section("server").Key("logfile").String()
			SQLiteWeb.clfLogFile = cfg.Section("server").Key("clflogfile").String()

			SQLiteWeb.WWWPath = cfg.Section("www").Key("path").String()
			SQLiteWeb.WWW404URL = cfg.Section("www").Key("404").MustString("/")

			SQLiteWeb.StubsPath = cfg.Section("stubs").Key("path").String()

			SQLiteWeb.Auth.Realm = cfg.Section("auth").Key("realm").String()
			SQLiteWeb.Auth.JWTTTL = cfg.Section("auth").Key("jwt_ttl").RangeInt64(300, 0, 0xFFFF)
			SQLiteWeb.Auth.JWTSecret = []byte(cfg.Section("auth").Key("jwt_key").String())
			SQLiteWeb.Auth.host = cfg.Section("auth").Key("host").String()
			SQLiteWeb.Auth.port = cfg.Section("auth").Key("port").RangeInt(8860, 0, 0xFFFF)
			SQLiteWeb.Auth.login = cfg.Section("auth").Key("login").String()
			SQLiteWeb.Auth.password = cfg.Section("auth").Key("password").String()
			SQLiteWeb.Auth.cert = cfg.Section("auth").Key("cert").String()

			initLogger()
			initDashboard()
			initAdmin()
			initStubs()
			initApi()
			initWWW()
			initCors()

			// inittt()

			if s, err := service.New(SQLiteWeb, svcConfig); err == nil {
				err = s.Run()
			} else {
				// log.Fatal(err)
				panic(err)
			}

		}
	}
}
