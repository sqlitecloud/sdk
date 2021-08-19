//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/08/18
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
import "bufio"
import "sqlitecloud"

func main() {
  out := bufio.NewWriter( os.Stdout )

  db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" )
  if err == nil {
    defer db.Close()

    db.CreateDatabase( "X", "", "UTF-8", true )
    db.Execute( `CREATE TABLE IF NOT EXISTS "Dummy" (ID INTEGER PRIMARY KEY AUTOINCREMENT, FirstName TEXT(20), LastName TEXT(20), ZIP INTEGER, City TEXT, Address TEXT)` )
    db.Execute( `DELETE FROM Dummy` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Some', 'One', 96450, 'Coburg', "Mohrenstraße 1" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Someone', 'Else', 96145, 'Sesslach', 'Raiffeisenstraße 6' )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'One', 'More', 91099, 'Poxdorf', "Langholzstr. 4" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Quotation', 'Test', 12345, '&"<>', "'Straße 0'" )` )

    sql := "SELECT * FROM Dummy"
    if res, err := db.Select( sql ); res != nil {
      defer res.Free()
      
      if err == nil {
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_LIST, false, "|", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_CSV, false, ",", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_QUOTE, false, ",", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_TABS, false, "\t", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_LINE, false, "", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_JSON, false, ",", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_HTML, false, "", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_MARKDOWN, false, "|", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_TABLE, false, "|", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_BOX, false, "│", "NULL", "\r\n", 0, false )
        res.DumpToWriter( out, sqlitecloud.OUTFORMAT_XML, false, sql, "NULL", "\r\n", 0, false )

        return
  } } }
  panic( err )
}