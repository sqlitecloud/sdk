//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/10/01
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Minimal SQLite Cloud SDK
//   ////                ///  ///                     test. Creates a table, inserts
//     ////     //////////   ///                      some values, uses DumpToWriter()
//        ////            ////                        to display all output formats.
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "fmt"
import "sqlitecloud"

func main() {
	   db  := sqlitecloud.New( "<USE INTERNAL PEM>", 10 )
  if err := db.Connect( "dev1.sqlitecloud.io", 9860, "", "", "X", 10, "NO", 0 ); err != nil { panic( err ) }

	defer db.Close()

	fmt.Printf( "PING?..." )
	switch db.Ping() {
	case nil: fmt.Printf( "PONG!\r\n" ) 
	default: 	fmt.Printf( "⚡️.\r\n" )
	}  
}