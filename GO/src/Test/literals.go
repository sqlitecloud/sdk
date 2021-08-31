//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/31
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Simple SQLite Cloud server
//   ////                ///  ///                     test. Creates a table,
//     ////     //////////   ///                      insertes some values, uses
//        ////            ////                        C-SDK's Dump function.
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "fmt"
import "sqlitecloud"

func main() {
  db, err := sqlitecloud.Connect( "sqlitecloud://***REMOVED***/X" )
  if err == nil {
    defer db.Close()

    // String Literal
    if res, err := db.Select( "PING" ); res != nil {
      defer res.Free()

      if err != nil {                             panic( err ) }
      if s, err := res.GetString(); err != nil {  panic( err )
      } else {                                    fmt.Printf( "String = %s\r\n", s ) }
    } else {                                      panic( err ) }

     // Integer
     if res, err := db.Select( "GET DATABASE ID" ); res != nil {
      defer res.Free()
      
      if err != nil {                             panic( err ) }
      if i, err := res.GetInt32(); err != nil {   panic( err )
     } else {                                     fmt.Printf( "Integer = %d\r\n", i ) }
    } else {                                      panic( err ) }

    // Float
    if res, err := db.Select( "GET LOAD" ); res != nil {
      defer res.Free()
      
      if err != nil {                             panic( err ) }
      if f, err := res.GetFloat32(); err != nil { panic( err )
      } else {                                    fmt.Printf( "Float = %f\r\n", f ) }
    } else {                                      panic( err ) }

    // Error
    if res, err := db.Select( "UNKNOWN COMMAND" ); res != nil {
     defer res.Free()

    } else {    
      if err != nil {                             fmt.Printf( "Error = %v\r\n", err ) }
    }
   
     return
  }
  panic( err )
}