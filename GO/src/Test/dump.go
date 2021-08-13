//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/08/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Simple SQLite Cloud server
//   ////                ///  ///                     test. Creates a table, inserts
//     ////     //////////   ///                      some values, uses DumpToWriter()
//        ////            ////                        to display all output formats.
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "os"
//import "io"
//import "fmt"
//import "strings"
import "sqlitecloud"
//import "encoding/json"
import "bufio"


func main() {
  out := bufio.NewWriter( os.Stdout )

  db, err := sqlitecloud.Connect( "sqlitecloud://***REMOVED***/X" )
  if err == nil {
    defer db.Close()

    db.CreateDatabase( "X", "", "UTF-8", true )
    db.Execute( `CREATE TABLE IF NOT EXISTS "Dummy" (ID INTEGER PRIMARY KEY AUTOINCREMENT, FirstName TEXT(20), LastName TEXT(20), ZIP INTEGER, City TEXT, Address TEXT)` )
    db.Execute( `DELETE FROM Dummy` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Some', 'One', 96450, 'Coburg', "Mohrenstrasse 1" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Someone', 'Else', 96145, 'Sesslach', 'Raiffeisenstrasse 6' )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'One', 'More', 91099, 'Poxdorf', "Langholzstr. 4" )` )

    sql := "SELECT * FROM Dummy"
    if res, err := db.Select( sql ); res != nil {
      defer res.Free()

      if err == nil {
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_LIST, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_CSV, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_QUOTE, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_TABS, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_LINE, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_JSON, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_HTML, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_MARKDOWN, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_TABLE, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_BOX, "|", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_XML, sql, 0, false )
        
        // out.Flush()
        return

  } } }
  panic( err )
}