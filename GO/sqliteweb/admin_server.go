//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/04/12
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
	"fmt"
	"net"
	"net/http"
	"strings"

	//"github.com/Shopify/go-lua"
	"github.com/gorilla/mux"
)

func init() {
	initializeSQLiteWeb()
}

func initAdmin() {
	if PathExists(cfg.Section("admin").Key("path").String()) && cfg.Section("admin").Key("enabled").MustBool(false) {
		// SQLiteWeb.router.HandleFunc( "/dashboard/{endpoint:.*}", SQLiteWeb.executeLua )

		SQLiteWeb.router.HandleFunc("/admin/{endpoint:.*}", SQLiteWeb.executeLuaAdminServer)
	}
}

func (this *Server) executeLuaAdminServer(writer http.ResponseWriter, request *http.Request) {
	if clientIP, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		if _, maskNet, err := net.ParseCIDR(cfg.Section("admin").Key("allow").String()); err == nil {
			if maskNet.Contains(net.ParseIP(clientIP)) {

				if username, password, ok := request.BasicAuth(); ok && username == cfg.Section("auth").Key("login").String() {
					switch cfg.Section("admin").Key("password").String() {
					// case password:
					// 	fallthrough
					case Hash(password):
						v := mux.Vars(request)
						endpoint := strings.ReplaceAll(v["endpoint"]+"/", "//", "/") // "v1/fbf94289-64b0-4fc6-9c20-84083f82ee63/database/Foo/connections/"
						path := cfg.Section("admin").Key("path").String()            // "/Users/pfeil/GitHub/SqliteCloud/sdk/GO/src/sqliteweb/dashboard"
						this.executeLua(path, endpoint, -1, writer, request)
						return
					default:
						break
					}
				}
			}
		}
	}

	writer.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", cfg.Section("auth").Key("realm").String()))
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte(fmt.Sprintf("{\"status\":%d,\"message\":\"%s\"}", http.StatusUnauthorized, "Invalid Credentials")))
}
