//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
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

import (
	// "encoding/json"
	// "fmt"
	// "text/template" // html/template
	// "io/ioutil"
	// "net"

	"net/http"
	"time"

	// "net/smtp"
	//"time"
	//"sqlitecloud"
	"strings"
	//"bytes"

	//"github.com/Shopify/go-lua"
	"github.com/gorilla/mux"
)

func init() {
  initializeSQLiteWeb()
}

func initDashboard() {
  if PathExists( cfg.Section( "dashboard" ).Key( "path" ).String() ) && cfg.Section( "dashboard" ).Key( "enabled" ).MustBool( false ) {
    // SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.executeLua )

    SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.executeLuaDashboardServer ) )
  }
}


func (this *Server) executeLuaDashboardServer( writer http.ResponseWriter, request *http.Request ) {
  start := time.Now()

  this.Auth.cors( writer, request )

	id, _    := SQLiteWeb.Auth.GetUserID( request )

  v        := mux.Vars( request ) // len(): 1, endpoint: v1/auth
  endpoint := strings.ReplaceAll( v[ "endpoint" ] + "/", "//", "/" ) // "v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/database/Foo/connections/"

  path     := cfg.Section( "dashboard" ).Key( "path" ).String()      // "/Users/pfeil/GitHub/SqliteCloud/sdk/GO/src/sqliteweb/dashboard"
  
	this.executeLua( path, endpoint, id, writer, request )

  t := time.Since( start )
  SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s", request.Method, request.URL, request.RemoteAddr, id, t)
}