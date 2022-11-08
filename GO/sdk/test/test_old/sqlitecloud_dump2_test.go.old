//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/18
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : Advanced SQLite Cloud server
//   ////                ///  ///                     test. Creates a table, inserts
//     ////     //////////   ///                      some values, uses DumpToWriter()
//        ////            ////                        to display all output formats.
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "os"
import "bufio"
import "sqlitecloud"

func main() {
  out := bufio.NewWriter( os.Stdout )

  db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" )
  if err == nil {
    defer db.Close()

    db.CreateDatabase( "X", "", "UTF-8", true )
    db.Execute( `CREATE TABLE IF NOT EXISTS "Dummy" (ID INTEGER PRIMARY KEY AUTOINCREMENT, FirstName TEXT(20), LastName TEXT(20), ZIP FLOAT, City TEXT, Address TEXT)` )
    db.Execute( `DELETE FROM Dummy` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Some', 'One', 96.450, 'Coburg', "Mohrenstraße 1" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Someone', 'Else', 96145, 'Sesslach', 'Raiffeisenstraße 6' )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'One', 'More', 91099, 'Poxdorf', "Langholzstr. 4" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Quotation', 'Test', 12345, '&"<>', "'Straße 0'" )` )

    sql := "SELECT * FROM Dummy"
    if res, err := db.Select( sql ); res != nil {
      defer res.Free()

      if err == nil {
        for index, FormatName := range []string{ "LIST", "CSV", "QUOTE", "TABS", "LINE", "JSON", "HTML", "MARKDOWN", "TABLE", "BOX" } {
          Format, _ := sqlitecloud.GetOutputFormatFromString( FormatName )
          res.DumpToWriter( out, Format, false, "<AUTO>", "NULL", "\r\n", 25 + uint( index * 5 ), false )
        }
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_XML, false, sql, "NULL", "\r\n", 25, false )
        return
  } } }
  panic( err )
}