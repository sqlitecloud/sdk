//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2022/03/25
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
	"net/http"
)

func init() {
	initializeSQLiteWeb()

	SQLiteWeb.router.HandleFunc("/dashboard/v1/{path:.*}", func(writer http.ResponseWriter, request *http.Request) {
		SQLiteWeb.Auth.cors(writer, request)
	}).Methods("OPTIONS")
}

func (this *AuthServer) cors(writer http.ResponseWriter, request *http.Request) {
	// only for debuging Mauros front end
	// dont forget to remove it!!!

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with, X-SQLiteCloud-Api-Key")
}
