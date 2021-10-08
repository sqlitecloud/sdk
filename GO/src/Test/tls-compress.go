//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : SQLite Cloud server test
//   ////                ///  ///                     Creates a table, inserts many
//     ////     //////////   ///                      values, enables compress and
//        ////            ////                        selectes all values.
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import "fmt"
import "math/rand"
import "sqlitecloud"

func main() {
  db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X?&compress=LZ4" )
  if err == nil {
    defer db.Close()

    db.CreateDatabase( "X", "", "UTF-8", true )
    db.Execute( `CREATE TABLE IF NOT EXISTS "CompressTest" (ID INTEGER PRIMARY KEY AUTOINCREMENT, Dummy TEXT(200) )` )
    //db.Execute( `DELETE FROM CompressTest`)

    rndStr := ""
    for i := 0; i < 10; i++ { rndStr = fmt.Sprintf( "%s%d", rndStr, rand.Int() ) }

    sql := ""
    for rows := 0; rows < 1000; rows++ {
      sql = sql + fmt.Sprintf( "INSERT INTO CompressTest ( Dummy ) VALUES( '%s' ); ", rndStr )
    }
    // db.Execute( sql )

    db.Compress( "lz4" )

    if res, err := db.Select( "SELECT * FROM CompressTest" ); res != nil {
      defer res.Free()

      if err == nil {
        res.Dump()
        return

  } } }
  panic( err )
}