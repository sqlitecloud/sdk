//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.1
//     //             ///   ///  ///    Date        : 2021/12/20
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

import "sqlitecloud"

func init() {
  initializeSQLiteWeb()
  if db == nil {
    db, _ = sqlitecloud.Connect( "sqlitecloud://admin:admin@dev1.sqlitecloud.io/Test" );
  }
  if adb == nil {
    adb, _ = sqlitecloud.Connect( "sqlitecloud://admin:kAhqTYvgXX43@auth1.sqlitecloud.io/users.sqlite" );
  }
}

func initDashboard() {
  if PathExists( cfg.Section( "dashboard" ).Key( "path" ).String() ) {
    // SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.executeLua )

    SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.executeLua ) )
  }
}