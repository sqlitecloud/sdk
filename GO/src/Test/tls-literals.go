//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.2.0
//     //             ///   ///  ///    Date        : 2021/10/14
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

    switch res, err = db.Select( "TEST NULL" ); {             // NULL
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsNULL():     panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "NULL = %v\r\n", err )
    }
    res.Free()

    switch res, err = db.Select( "TEST STRING" ); {           // String Literal
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsString():   panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "String = %s\r\n", res.GetString_() )
    }
    res.Free()

    switch res, err = db.Select( "TEST ZERO_STRING" ); {      // String Literal
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsString():   panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "String = %s\r\n", res.GetString_() )
    }
    res.Free()

    switch res, err = db.Select( "TEST ERROR" ); {            // Error
    case err == nil:        panic( errors.New( "Unknown command returned no error" ) )
    default:                fmt.Printf( "Error = %v\r\n", err )
    }
        
    switch res, err = db.Select( "TEST INTEGER" ); {          // Integer
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsInteger():  panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "Integer = %d\r\n", res.GetInt32_() )
    }
    res.Free()
      
    switch res, err = db.Select( "TEST FLOAT" ); {            // Float
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsFloat():    panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "Float = %f\r\n", res.GetFloat32_() )
    }
    res.Free()

    switch res, err = db.Select( "TEST BLOB" ); {            // BLOB
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsBLOB():     panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "BLOB LEN = %d\r\n", len( res.GetBuffer() ) )
    }
    res.Free()

    switch res, err = db.Select( "TEST ROWSET" ); {          // ROWSET
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsRowSet():   panic( errors.New( "invalid type" ) )
    default:                res.DumpToScreen( 0 )
    }
    res.Free()

    switch res, err = db.Select( "TEST JSON" ); {            // JSON
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsJSON():     panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "JSON = %s\r\n", res.GetString_() )
    }
    res.Free()

    switch res, err = db.Select( "TEST COMMAND" ); {        // Command
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsCommand():  panic( errors.New( "invalid type" ) )
    default:                fmt.Printf( "COMMAND = %v\r\n", res.GetString_() ) // should be pong
    }
    res.Free()

    switch res, err = db.Select( "TEST ARRAY" ); {          // ARRAY
    case err != nil:        panic( err )
    case res == nil:        panic( errors.New( "nil result" ) )
    case !res.IsArray():    panic( errors.New( "This is not an array" ) )
    default:
      for row := uint64( 0 ); row < res.GetNumberOfRows(); row++ {
        fmt.Printf( "[%d] %s\r\n", row, res.GetStringValue_( row, 0 ) )
      }
      // res.DumpToScreen( 0 )
    }
    res.Free()




  } else { panic( err ) }
}