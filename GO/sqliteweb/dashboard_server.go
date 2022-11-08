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
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	// _ "net/http/pprof"
)

func init() {
	initializeSQLiteWeb()
}

func initDashboard() {
	if PathExists(cfg.Section("dashboard").Key("path").String()) && cfg.Section("dashboard").Key("enabled").MustBool(false) {
		// SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.executeLua )

		// special dashboard paths processed with WebSocket connections
		// SQLiteWeb.router.HandleFunc("/dwsTestClient", dwsTestClient) // only for test purpose
		SQLiteWeb.router.HandleFunc("/dashboard/{version:v[0-9]+}/{projectID}/database/{databaseName}/download", SQLiteWeb.Auth.JWTAuth(SQLiteWeb.Auth.getTokenFromCookie, SQLiteWeb.dashboardWebsocketDownload))
		SQLiteWeb.router.HandleFunc("/dashboard/{version:v[0-9]+}/{projectID}/database/{databaseName}/upload", SQLiteWeb.Auth.JWTAuth(SQLiteWeb.Auth.getTokenFromCookie, SQLiteWeb.dashboardWebsocketUpload))

		// catch all with executeLuaDashboardServer
		SQLiteWeb.router.HandleFunc("/dashboard/{endpoint:.*}", SQLiteWeb.Auth.JWTAuth(SQLiteWeb.Auth.getTokenFromAuthorization, SQLiteWeb.executeLuaDashboardServer))

		// SQLiteWeb.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	}
}

func (this *Server) executeLuaDashboardServer(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	this.Auth.cors(writer, request)

	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromAuthorization, request)

	v := mux.Vars(request)                                       // len(): 1, endpoint: v1/auth
	endpoint := strings.ReplaceAll(v["endpoint"]+"/", "//", "/") // "v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/database/Foo/connections/"

	path := cfg.Section("dashboard").Key("path").String() // "/Users/pfeil/GitHub/SqliteCloud/sdk/GO/src/sqliteweb/dashboard"

	this.executeLua(path, endpoint, id, writer, request)

	t := time.Since(start)
	SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s", request.Method, request.URL, request.RemoteAddr, id, t)
}
