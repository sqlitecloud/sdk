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

import "os"
//import "io"
import "fmt"
//import "bytes"
//import "time"
//import "context"
//import "errors"
import "strings"
//import "strconv"
//import "sqlitecloud"

//import "github.com/kardianos/service"
import "net/http"


import "github.com/gorilla/mux"
// import "github.com/gorilla/websocket"


func init() {
	initializeSQLiteWeb()
}

func initStubs() {
	if PathExists( SQLiteWeb.APIPath ) {
		SQLiteWeb.router.HandleFunc( "/api/{endpoint:.*}", SQLiteWeb.stubHandler)
	}	
}

func (this *Server) stubHandler(writer http.ResponseWriter, request *http.Request) {
	v        := mux.Vars( request )
	endpoint := strings.ReplaceAll( v[ "endpoint" ] + "/", "//", "/" )
	file := fmt.Sprintf( "%s/%s%s.json", this.APIPath, endpoint, strings.TrimSpace( strings.ToUpper( request.Method ) ) )
	
	println( file )

	if PathExists( file ) {
		dat, err := os.ReadFile( file )
		if err == nil {
			writer.Header().Set("Content-Type", "application-json")
			writer.Header().Set("Content-Encoding", "utf-8")
			writer.Write( dat )
		} else {
			http.Error( writer, err.Error(), http.StatusInternalServerError )
		}
	} else {
		http.Error( writer, "Endpoint not found.", http.StatusNotFound )
	}
}