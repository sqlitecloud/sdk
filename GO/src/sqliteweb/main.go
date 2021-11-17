//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/11/17
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
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH
( should be something like: /Users/pfeil/GitHub/SqliteCloud/sdk/GO )
cd src/sqliteweb/
go run *.go

cd src
go get github.com/gorilla/websocket
go get -u github.com/gorilla/mux
go get gopkg.in/ini.v1
// go get https://github.com/judwhite/go-svc
go get github.com/kardianos/service

// Compile with:
GOOS=linux go build -o sqliteweb  *.go 

*/

import "os"
//import "io"
import "fmt"
//import "time"
//import "errors"
import "strings"
//import "strconv"
//import "sqlitecloud"
import "github.com/docopt/docopt-go"
import "github.com/kardianos/service"
import "gopkg.in/ini.v1"
//import "github.com/gorilla/mux"
//import "github.com/gorilla/websocket"

var app_name     = "sqliteweb"
var long_name    = "SQLite Cloud Web Server"
var version      = "version 0.0.1"
var copyright    = "(c) 2021 by SQLite Cloud Inc."

var cfg *ini.File

func main() {
	// Read the .ini file -> https://github.com/go-ini/ini
	var err error
  if cfg, err = ini.Load( "sqliteweb.ini" ); err == nil {

    // Overload the .ini file with the command line arguments
		if p, err := docopt.ParseArgs( strings.ReplaceAll( usage, "::", ": " ), nil, fmt.Sprintf( "%s %s, %s", app_name, version, copyright ) ); err == nil {
			for key, a := range p {
				if a != nil { 
					value := fmt.Sprintf( "%v", a )
					switch strings.TrimSpace( strings.ToLower( key ) ) {
					case "--stubs":		cfg.Section( "stubs" ).	Key( "path" ).	 SetValue( value )
					case "--www":	  	cfg.Section( "www" ).	  Key( "path" ).	 SetValue( value )
					case "--address": cfg.Section( "server" ).Key( "address" ).SetValue( value )
					case "--port": 		cfg.Section( "server" ).Key( "port"    ).SetValue( value )
					default: // fmt.Printf( "%10s = %v\r\n", key, a )
					}
				}
			}
		}
		// start the service -> https://github.com/kardianos/service
		svcConfig := &service.Config {
			Name:        app_name,  // "GoServiceExampleSimple",
			DisplayName: long_name, // "Go Service Example",
			Description: fmt.Sprintf( "%s %s, %s", long_name, version, copyright ), // "This is an example Go service.",
		}

		initializeSQLiteWeb()
		SQLiteWeb.Address 	= cfg.Section( "server" ).Key("address").String()
		SQLiteWeb.Port 			= cfg.Section( "server" ).Key("port").RangeInt( 8433, 0, 0xFFFF )
		SQLiteWeb.WWWPath 	= cfg.Section( "www" ).   Key("path").String()
		SQLiteWeb.StubsPath = cfg.Section( "stubs" ). Key("path").String()

		initStubs()
		initWWW()

		if s, err := service.New( SQLiteWeb, svcConfig ); err == nil {
			err = s.Run()
		} else {
				// log.Fatal(err)
			panic( err )
		}

  } else {
		fmt.Printf( "Fail to read file: %v", err )
    os.Exit(1)
	}
}