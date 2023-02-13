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

	"github.com/gorilla/mux"
	// _ "net/http/pprof"
)

var dashboardcm *ConnectionManager = nil

func init() {
	initializeSQLiteWeb()
	dashboardcm, _ = NewConnectionManager()
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
	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromAuthorization, request)

	v := mux.Vars(request)                                       // len(): 1, endpoint: v1/auth
	endpoint := strings.ReplaceAll(v["endpoint"]+"/", "//", "/") // "v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/database/Foo/connections/"

	path := cfg.Section("dashboard").Key("path").String() // "/Users/pfeil/GitHub/SqliteCloud/sdk/GO/src/sqliteweb/dashboard"

	this.executeLua(path, endpoint, id, writer, request)
}

// // SQLiteWeb.router.HandleFunc("/dashboard/{version:v[0-9]+}/{projectID}/node", SQLiteWeb.Auth.JWTAuth(SQLiteWeb.Auth.getTokenFromAuthorization, SQLiteWeb.createNodeHandler)).Methods(http.MethodPost)
//
// func (this *Server) createNodeHandler(writer http.ResponseWriter, request *http.Request) {
// 	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromAuthorization, request)
// 	v := mux.Vars(request)

// 	projectUUID, _, err := verifyProjectID(id, v["projectID"], dashboardcm)
// 	if err != nil {
// 		writeError(writer, http.StatusUnauthorized, err.Error(), "")
// 		return
// 	}

// 	insertjob := fmt.Sprintf("INSERT INTO Jobs (uuid, name, steps, progress, user_id, project_uuid) VALUES ('%s', 'Create Node', 4, 0, %d, %s); SELECT uuid FROM Jobs WHERE rowid = last_insert_rowid();", uuid.New(), id, projectUUID)
// 	res, err, _, _, _ := dashboardcm.ExecuteSQL("auth", insertjob)
// 	if err != nil {
// 		writeError(writer, http.StatusInternalServerError, err.Error(), "")
// 		return
// 	}

// 	jobuuid, err := res.GetStringValue(0, 0)
// 	if err != nil {
// 		writeError(writer, http.StatusInternalServerError, err.Error(), "")
// 		return
// 	}

// 	go createNode(jobuuid, v)

// 	statusCode := http.StatusOK
// 	writer.Header().Set("Content-Type", "application/json")
// 	writer.Header().Set("Content-Encoding", "utf-8")
// 	writer.WriteHeader(statusCode)
// 	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\", \"value\": \"{\"uuid\": \"%s\" }\"}", statusCode, "OK", jobuuid)))
// }
