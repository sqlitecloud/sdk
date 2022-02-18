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

func init() {
  initializeSQLiteWeb()
}

func initDashboard() {
  if PathExists( cfg.Section( "dashboard" ).Key( "path" ).String() ) {
    // SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.executeLua )

    SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.Auth.JWTAuth( SQLiteWeb.executeLua ) )
  }
}