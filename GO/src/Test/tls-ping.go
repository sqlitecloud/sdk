//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/18
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
  err := db.Connect( "dev1.sqlitecloud.io", 9860, "", "", "X", 10, "NO", 0 )

  if err == nil {
    defer db.Close()

		fmt.Printf( "PING..." )
		if db.Ping() == nil { 
			fmt.Printf( "PONG!\r\n" ) 
		}	else {
			fmt.Printf( "⚡️.\r\n" )
		}
		return
	}
  panic( err )
}