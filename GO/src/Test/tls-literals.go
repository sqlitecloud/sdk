//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/01
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

import (
  "errors"
  "fmt"
  "sqlitecloud"
)

func main() {
  var db    *sqlitecloud.SQCloud
  var res   *sqlitecloud.Result
  var err   error

  if db, err = sqlitecloud.Connect( "sqlitecloud://***REMOVED***/X" ); err == nil {   
    defer db.Close()

    switch res, err = db.Select( "PING" ); {            // String Literal
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "" ) )
    case !res.IsString():   panic( errors.New( "" ) )
    default:                fmt.Printf( "String = %s\r\n", res.GetString_() )
    }
    res.Free()

    switch res, err = db.Select( "GET DATABASE ID" ); { // Integer
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "" ) )
    case !res.IsInteger():  panic( errors.New( "" ) )
    default:                fmt.Printf( "Integer = %d\r\n", res.GetInt32_() )
    }
    res.Free()
      
    switch res, err = db.Select( "GET LOAD" ); {        // Float
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "" ) )
    case !res.IsFloat():    panic( errors.New( "" ) )
    default:                fmt.Printf( "Float = %f\r\n", res.GetFloat32_() )
    }
    res.Free()

    switch res, err = db.Select( "UNKNOWN COMMAND" ); { // Error
    case err == nil:        panic( errors.New( "Unknown command returned no error" ) )
    default:                fmt.Printf( "Error = %v\r\n", err )
    }
  } else { panic( err ) }
}