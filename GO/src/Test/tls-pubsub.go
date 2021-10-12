//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/10/11
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Simple PSUB Test
//   ////                ///  ///       
//     ////     //////////   ///        
//        ////            ////          
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "fmt"
import "time"
import "sqlitecloud"

func main() {
  fmt.Printf( "Simple PSUB test...\r\n")
  
  if db, _ := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" ); db != nil {
    defer db.Close()

    db.Callback = func( JSON string ) {
      println( "The json string was: " + JSON )
    }
      
    db.Listen( "Channel" )

    db.Notify( "Channel" )
    db.SendNotificationMessage( "Channel", "Hello dear fellows!" )
    
    println( "Waiting for Messages, press CTL+C to exit..." )
    time.Sleep( 10000 * time.Second )
} }