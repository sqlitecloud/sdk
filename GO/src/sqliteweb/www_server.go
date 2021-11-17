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

import "net/http"

func init() {
	initializeSQLiteWeb()
}

func initWWW() {
	if PathExists( SQLiteWeb.WWWPath ) {
		SQLiteWeb.router.PathPrefix( "/" ).Handler(http.StripPrefix( "/", http.FileServer( http.Dir( SQLiteWeb.WWWPath ) ) ) )
	}	
}