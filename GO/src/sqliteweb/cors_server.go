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

import (
	// "fmt"
	// "time"
	// "encoding/json"
	"net/http"

	// "github.com/golang-jwt/jwt"
)

func init() {
	initializeSQLiteWeb()

	SQLiteWeb.router.HandleFunc("/dashboard/v1/{path:.*}", func(writer http.ResponseWriter, request *http.Request) {
		SQLiteWeb.Auth.cors(writer, request)
	}).Methods("OPTIONS")
}


func (this *AuthServer) cors( writer http.ResponseWriter, request *http.Request ) {
	// only for debuging Mauros front end
	// dont forget to remove it!!!
	
	writer.Header().Set( "Access-Control-Allow-Origin", "*" )
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}