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

//import "os"
//import "io"
import "fmt"
import "time"
import "context"
import "errors"
//import "strings"
//import "strconv"
//import "sqlitecloud"

import "github.com/kardianos/service"

import "net/http"
import "github.com/gorilla/mux"
// import "github.com/gorilla/websocket"

type Server struct{
	Address       string
	Port          int   

	Hostname			string
	CertPath			string
	KeyPath 			string

	Auth 					AuthServer

	WWWPath       string
	APIPath				string

	server  *http.Server
	router  *mux.Router
	ticker 	*time.Ticker
}

var SQLiteWeb *Server = nil

func initializeSQLiteWeb() {
	if SQLiteWeb == nil {
		SQLiteWeb = &Server{
			Address: 		"127.0.0.1",
			Port: 			8433,

			Hostname:		"",
			CertPath:   "",
			KeyPath:    "",

			Auth:				AuthServer{
				JWTSecret:  []byte( "" ),
				JWTTTL:     0,	
				Tokens:		map[string]TokenInfo{},
			},

			WWWPath: 		"",
			APIPath: 		"",
			server: 		nil,
			router: 		mux.NewRouter(),
			ticker: 		nil,
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

	fmt.Printf( "%s, starting at: %s:%d\r\n", long_name, this.Address, this.Port )
	go this.run()
	return nil
}
func ( this *Server ) run() {
		if this.server.ListenAndServeTLS( this.CertPath, this.KeyPath ) != http.ErrServerClosed { // Diese Zeile hängt jetzt bis Shutdown aufgerufen wird...
		// ich wurde friedlich per Shutdown beendet...
		println( "Drin" )
	}
	println( "Ende" )
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