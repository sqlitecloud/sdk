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
	"bufio"
	"fmt"
	"github.com/Shopify/go-lua" //import "io"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//import "bytes"
//import "time"
//import "context"
//import "errors"
//import "strconv"
//import "sqlitecloud"
//import "github.com/kardianos/service"
//import "github.com/gorilla/websocket"

var out = bufio.NewWriter( os.Stdout )

func init() {
  initializeSQLiteWeb()
}

func initStubs() {
  if PathExists( SQLiteWeb.StubsPath ) {
    SQLiteWeb.router.HandleFunc( "/stubs/{endpoint:.*}", SQLiteWeb.stubHandler)
  }
}

func (this *Server) stubHandler(writer http.ResponseWriter, request *http.Request) {
  this.Auth.cors( writer, request )

  v        := mux.Vars( request )
  endpoint := strings.ReplaceAll( v[ "endpoint" ] + "/", "//", "/" )

  args     := []string{}
  path     := fmt.Sprintf( "%s", this.StubsPath )
  for _, part := range strings.Split( endpoint, "/" ) {
    //fmt.Printf( "i=%d, path=%s, part=%s\r\n", i, path, part )
    //fmt.Printf( "%s/{%d}\r\n", part, len( args ) )

    if PathExists( fmt.Sprintf( "%s/%s", path, part ) ) {
      path = fmt.Sprintf( "%s/%s", path, part )
    } else if PathExists( fmt.Sprintf( "%s/{%d}", path, len( args ) ) ) {
      path = fmt.Sprintf( "%s/{%d}", path, len( args ) )
      args = append( args, part )
    } else {
      panic( "not found" )
    }
  }
  path = fmt.Sprintf( "%s%s.", path, strings.TrimSpace( strings.ToUpper( request.Method ) ) )
  //fmt.Printf( "PATH=%s\r\n", path )

  switch {
  case PathExists( fmt.Sprintf( "%slua", path ) ):

    l := lua.NewState()
    lua.OpenLibraries( l )

    l.Register( "SetStatus", func(L *lua.State) int {
      if status, ok := L.ToInteger( 1 ); ok { writer.WriteHeader( status ) }
      return 0
    } )

    l.Register( "SetHeader", func(L *lua.State) int {
      if key, ok := L.ToString( 1 ); ok {
        if value, ok := L.ToString( 2 ); ok {
          //fmt.Printf( "SET HEADER: %s:%s\r\n", key, value )
          writer.Header().Set( key, value )
        }
      }
      return 0
    } )

    l.Register( "Write", func(L *lua.State) int {
      if data, ok := L.ToString( 1 ); ok {
        //fmt.Printf( "WRITE: %s\r\n", data )
        writer.Write( []byte( data ) )
      }
      return 0
    } )

    ////// l.Register( "queryNode", QueryNode )

    l.NewTable()
    l.PushInteger( 0 )
    l.PushString( fmt.Sprintf( "%slua", path ) )
    l.SetTable( -3 )

    for i, arg := range args {
      l.PushInteger( i + 1 )
      l.PushString( arg )
      l.SetTable( -3 )
    }
    l.SetGlobal( "args" )

    body, err := ioutil.ReadAll( request.Body )
    if err == nil {
      l.PushString( string( body ) )
      l.SetGlobal( "body" )
    }

    fmt.Printf( "will execute lua script: '%s'\r\n", path )
    fmt.Printf( "%v\r\n", args )

    //lua.DoString( l, fmt.Sprintf( `package.path = "%s"`, this.LUAPath ) )

    err = lua.DoFile( l, fmt.Sprintf( "%slua", path ) )
    if err != nil {
      panic( err )
    }
    return

  case PathExists( fmt.Sprintf( "%sjson", path ) ):
    fmt.Printf( "will return json\r\n", path )
  default:
    panic( "not found" )
  }



  file := fmt.Sprintf( "%s/%s%s.json", this.StubsPath, endpoint, strings.TrimSpace( strings.ToUpper( request.Method ) ) )
  //println( file )

  if PathExists( file ) {
    dat, err := os.ReadFile( file )
    if err == nil {

      //this.cors( writer, request )

      writer.Header().Set("Content-Type", "application/json")
      writer.Header().Set("Content-Encoding", "utf-8")
      writer.Write( dat )
    } else {
      http.Error( writer, err.Error(), http.StatusInternalServerError )
    }
  } else {
    http.Error( writer, "Endpoint not found.", http.StatusNotFound )
  }
}