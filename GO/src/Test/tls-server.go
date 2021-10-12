//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Test program for the 
//   ////                ///  ///                     SQLite Cloud internal 
//     ////     //////////   ///                      server commands.
//        ////            ////          
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "os"
import "fmt"
import "sqlitecloud"

func main() {
  fmt.Printf( "Server API test...\r\n")

  db := sqlitecloud.New( "<USE INTERNAL PEM>", 10 )
  err := db.Connect( "***REMOVED***", 8860, "user with space", "password with space", "X", 10, "NO", 0 ) // Host, Port, Username, Password, Database, Timeout, Compression, Family
  if err == nil {
    defer db.Close()

    fmt.Printf( "Checking PING..." )
    if db.Ping() != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking AUTH (without space)..." )
    if err := db.Auth( "pfeil", "secret"); err != nil { //  Username, Password
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking AUTH (with space)..." )
    if err := db.Auth( "user with space", "password with space"); err != nil { //  Username, Password
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking CREATE DATABASE..." )
    if err := db.CreateDatabase( "xyz", "", "", false ); err != nil { // Database, Key, Encoding, NoError
      fail( err.Error() )
    } 
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST DATABASES..." )
    if databases, err := db.ListDatabases(); err != nil {
      fail( err.Error() )
    } else {
      if len( databases ) == 0 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking DROP DATABASE..." )
    if err := db.DropDatabase( "xyz", false ); err != nil { // Database, NoError
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )
    


    fmt.Printf( "Checking USE DATABASE..." )
    db.CreateDatabase( "X", "", "", true )
    if err := db.UseDatabase( "X" ); err != nil { // Database
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking GET DATABASE ID..." )
    id, err := db.GetDatabaseID()
    if err != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )


    fmt.Printf( "Checking SET KEY..." )
    if err := db.SetKey( "A", "1405" ); err != nil { // Key, Value
      fail( err.Error() )
    } 
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST KEYS..." )
    if keys, err := db.ListKeys() ; err != nil { 
      fail( err.Error() )
    } else {
      if len( keys ) != 1 {
        fail( "Invalid amount of keys found." )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking GET KEY..." )
    if val, err := db.GetKey( "A" ) ; err != nil { 
      fail( err.Error() )
    } else {
      if val != "1405" {
        fail( fmt.Sprintf( "Invalid value (%s).", val ) )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking DROP KEY..." )
    if err := db.DropKey( "A" ); err != nil {
      fail( err.Error() )
    } else {
      if keys, err := db.ListKeys() ; err != nil {
        fail( err.Error() )
      } else {
        if len( keys ) != 0 {
          fail( "Key could not be deleted." )
        }
      }
      fmt.Printf( "ok.\r\n" )
    }


    


    fmt.Printf( "Checking LIST COMMANDS..." )
    if commands, err := db.ListCommands(); err != nil {
      fail( err.Error() )
    } else {
      if len( commands ) != 28 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )


    fmt.Printf( "Checking LIST CONNECTIONS..." )
    if connections, err := db.ListConnections(); err != nil {
      fail( err.Error() )
    } else {
      if len( connections ) == 0 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST DATABASE CONNECTIONS (X)..." )
    if connections, err := db.ListDatabaseConnections( "X" ); err != nil {
      fail( err.Error() )
    } else {
      if len( connections ) == 0 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )


  

    fmt.Printf( "Checking LIST DATABASE CONNECTIONS ID %d...", id )
    if connections, err := db.ListDatabaseClientConnectionIds( id ); err != nil {
      fail( err.Error() )
    } else {
      if len( connections ) == 0 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST INFO..." )
    if _, err := db.GetInfo(); err != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST TABLES..." )
    if tables, err := db.ListTables(); err != nil {
      fail( err.Error() )
    } else {
      if len( tables ) < 1 {
        fail( "Ivalid result." )
      }
    }
    fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST PLUGINS..." )
    if _, err := db.ListPlugins(); err != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )
  

    fmt.Printf( "Checking LIST CLIENT KEYS..." )
    if _, err := db.ListClientKeys(); err != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )




    // fmt.Printf( "Checking LIST NODES..." )
    // if nodes, err := db.ListNodes(); err != nil {
    //  fail( err.Error() )
    // } else {
    //  if len( nodes ) == 0 {
    //    fail( "Ivalid result." )
    //  }
    // }
    // fmt.Printf( "ok.\r\n" )

    fmt.Printf( "Checking LIST DATABASE KEYS..." )
    if _, err := db.ListDatabaseKeys(); err != nil {
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )




    fmt.Printf( "Checking UNUSE DATABASE..." )
    if err := db.UnuseDatabase(); err != nil { // Database
      fail( err.Error() )
    }
    fmt.Printf( "ok.\r\n" )

    // fmt.Printf( "Checking CLOSE CONNECTION..." )
    // if err := db.CloseConnection( "14" ); err != nil { // ConnectionID
    //  fail( err.Error() )
    // }
    // fmt.Printf( "ok.\r\n" )
    

    fmt.Printf( "Success.\r\n" )
    os.Exit( 0 )
  }
  fail( err.Error() )
}

func fail( Message string ) {
  fmt.Printf( "failed: %s\r\n", Message )
  os.Exit( 1 )
}