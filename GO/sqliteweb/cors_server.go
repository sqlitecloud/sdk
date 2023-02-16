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

	"github.com/gorilla/mux"
)

const allowAllOrigins = "*"

func initCors() {
	SQLiteWeb.router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	apiRoute := SQLiteWeb.router.PathPrefix("/api/")
	webRoute := SQLiteWeb.router.PathPrefix("/web/")
	routeMatch := mux.RouteMatch{}

	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with"

			// allow all origins only for /api/* endpoints
			// TODO: set the allowed origins for the project in the API settings page of the dashboard
			switch {
			case apiRoute.Match(req, &routeMatch):
				allowHeaders += ", X-SQLiteCloud-Api-Key"
				w.Header().Set("Access-Control-Allow-Origin", allowAllOrigins)
			case webRoute.Match(req, &routeMatch):
				w.Header().Set("Access-Control-Allow-Origin", allowAllOrigins)
			}

			if DEBUG_SQLITEWEB {
				w.Header().Set("Access-Control-Allow-Origin", allowAllOrigins)
			}

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

			// and call next handler!
			next.ServeHTTP(w, req)
		})
}
