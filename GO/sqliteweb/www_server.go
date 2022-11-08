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

import "strings"
import "net/http"

type fileHandler404 struct {
  root        http.FileSystem
  notFoundURL string
}
func FileServer404( root http.FileSystem, NotFoundURL string ) http.Handler {
  return &fileHandler404{ root, NotFoundURL }
}
func ( this *fileHandler404 ) ServeHTTP( writer http.ResponseWriter, request *http.Request ) {
  // URL rewriter
  // ... if necessary

  // 404 checker
  file := SQLiteWeb.WWWPath + strings.ReplaceAll( "/" + request.URL.Path, "//", "/" )
  if PathExists( file ) {
    fs := http.FileServer( this.root )
    fs.ServeHTTP( writer, request )
  } else {
    http.Redirect( writer, request, this.notFoundURL, http.StatusMovedPermanently )
  } 
}


func init() {
  initializeSQLiteWeb()
}

func initWWW() {
  if PathExists( SQLiteWeb.WWWPath ) {
    SQLiteWeb.router.PathPrefix( "/" ).Handler( http.StripPrefix( "/", FileServer404( http.Dir( SQLiteWeb.WWWPath ), SQLiteWeb.WWW404URL ) ) )  
  } 
}