//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/08/16
//    ///             ///   ///  ///    Author      : Andrea Donetti
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

var apicm *ConnectionManager = nil

func init() {
	initializeSQLiteWeb()
	apicm, _ = NewConnectionManager()
}

func initApi() {
	if cfg.Section("api").Key("enabled").MustBool(false) {
		// SQLiteWeb.router.HandleFunc("/api/apiWebsocketTest", apiWebsocketTestClient) // only for test purpose
		SQLiteWeb.router.HandleFunc("/api/{version:v[0-9]+}/{projectID}/ws", SQLiteWeb.serveApiWebsocket)
		SQLiteWeb.router.HandleFunc("/api/{version:v[0-9]+}/wspsub", SQLiteWeb.serveApiWebsocketPubsub)

		SQLiteWeb.router.HandleFunc("/api/{version:v[0-9]+}/{projectID}/{databaseName}", SQLiteWeb.serveApiRest).Methods("GET")
		SQLiteWeb.router.HandleFunc("/api/{version:v[0-9]+}/{projectID}/{databaseName}/{tableName}", SQLiteWeb.serveApiRest)
		SQLiteWeb.router.HandleFunc("/api/{version:v[0-9]+}/{projectID}/{databaseName}/{tableName}/{id}", SQLiteWeb.serveApiRest)
	}
}
