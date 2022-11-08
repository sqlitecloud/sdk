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

import "sqlitecloud"

func main() {
  db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" )
  if err == nil {
    defer db.Close()

    db.CreateDatabase( "X", "", "UTF-8", true )
    db.Execute( `CREATE TABLE IF NOT EXISTS "Dummy" (ID INTEGER PRIMARY KEY AUTOINCREMENT, FirstName TEXT(20), LastName TEXT(20), ZIP INTEGER, City TEXT, Address TEXT)` )
    db.Execute( `DELETE FROM Dummy` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Some', 'One', 96450, 'Coburg', "Mohrenstraße 1" )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Someone', 'Else', 96145, 'Sesslach', 'Raiffeisenstraße 6' )` )
    db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'One', 'More', 91099, 'Poxdorf', "Langholzstr. 4" )` )

    sql := "SELECT * FROM Dummy"
    if res, _ := db.Select( sql ); res != nil {
      defer res.Free()
      if err == nil {
        res.Dump()
        return
  } } }
  panic( err )
}