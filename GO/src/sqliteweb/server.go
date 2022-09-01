//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/02/15
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

//import "os"
//import "io"
import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

//import "strings"
//import "strconv"
//import "sqlitecloud"

// import "github.com/gorilla/websocket"

type Server struct{
  Address       string
  Port          int

  Hostname      string
  CertPath      string
  KeyPath       string

  Logger        *Logger
  CLFWriter     io.Writer

  Auth          AuthServer

  WWWPath       string
  WWW404URL     string
  StubsPath     string

  server        *http.Server
  router        *mux.Router
  ticker        *time.Ticker

  logLevel      string
  logFile       string
  clfLogFile    string
}

var SQLiteWeb *Server = nil

func sendError( writer http.ResponseWriter, message string, statusCode int ) {
  writer.Header().Set( "Content-Type", "application/json" )
  writer.Header().Set( "Content-Encoding", "utf-8" )
	writer.WriteHeader( statusCode )
  writer.Write( []byte( fmt.Sprintf( "{\"status\":%d,\"message\":\"%s\"}", statusCode, message ) ) )
}

func initializeSQLiteWeb() {
  if SQLiteWeb == nil {
    SQLiteWeb = &Server{
      Address:    "127.0.0.1",
      Port:       8433,

      Hostname:   "",
      CertPath:   "",
      KeyPath:    "",

      Logger:     StdErrLogger,
      CLFWriter:  ioutil.Discard,

      Auth:       AuthServer{
        JWTSecret:  []byte( "" ),
        JWTTTL:     0,
      },

      WWWPath:    "",
      WWW404URL:  "/",
      StubsPath:  "",
      server:     nil,
      router:     mux.NewRouter(),
      ticker:     nil,

      logFile:    "",
      clfLogFile: "",
    }
  }
}

func init() {
  initializeSQLiteWeb()
}

func ( this *Server ) Start( s service.Service ) error {
  if this.router == nil {
    this.router = mux.NewRouter()
  }
  if this.router == nil { return errors.New( "XXX" ) }

  this.router.Use(CommonLogFormatMiddleware)

  this.server = &http.Server{
    Addr:         fmt.Sprintf( "%s:%d", this.Address, this.Port ),
    Handler:      this.router,

    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }
  if this.server == nil { return errors.New( "YYY" ) }

  if this.ticker == nil {
    this.ticker = time.NewTicker( time.Second )
    go func() {
      for {
        <-this.ticker.C
        this.tick()
    } }()
  }
  if this.ticker == nil { return errors.New( "ZZZ" ) }

  SQLiteWeb.Logger.Infof( "%s, starting at: %s:%d\r\n", long_name, this.Address, this.Port )
  go this.run()
  return nil
}
func ( this *Server ) run() {
   // This line now hangs until shutdown is called
  err := this.server.ListenAndServeTLS( this.CertPath, this.KeyPath )
  // After Shutdown or Close, the returned error is ErrServerClosed.
  SQLiteWeb.Logger.Infof("SQLiteWeb ListenAndServeTLS ended: %s", err)
  //this.Stop()
}
func (this *Server ) Stop( s service.Service ) error {
  if this.ticker != nil {
    this.ticker.Stop()
    this.ticker = nil
  }

  if this.server != nil {
    err := this.server.Shutdown( context.TODO() )
    if err != nil {
      // failure/timeout shutting down the server gracefully
      panic(err)
    }
    this.server = nil
  }

  // Stop should not block. Return with a few seconds.
  return nil
}

func (this *Server ) tick() {
  now := int( time.Now().Unix() )

  // every second
  if now % 1 == 0 {
    //println("tick...")
  }
  // every 2 seconds
  if now % 2 == 0 {

  }
}