//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andrea Donetti
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

import (
	"fmt"
	"os"
	"sqlitecloud"
)

const dbname = "test-selectarray-db.sqlite"

func main() {
	fmt.Printf("Server API test...\r\n")

	config, err1 := sqlitecloud.ParseConnectionString("sqlitecloud://admin:admin@localhost:8850")
	if err1 != nil {
		fail(err1.Error())
	}

	db := sqlitecloud.New(*config)
	err := db.Connect()

	if err == nil {
		defer db.Close()

		// test select null string, it was causing a crash on the server
		if res, err :=  db.Select(""); err != nil {
			fail(err.Error())
		} else {
			res.DumpToScreen(0)
		}

		fmt.Printf("Checking CREATE DATABASE...")
		if err := db.ExecuteArray("CREATE DATABASE ? PAGESIZE ? IF NOT EXISTS", []interface{}{dbname, 4096}); err != nil {
			fail(err.Error())
		}
		fmt.Printf("ok.\r\n")

		fmt.Printf("Checking USE DATABASE...")
		if err := db.UseDatabase(dbname); err != nil { // Database
			fail(err.Error())
		}
		fmt.Printf("ok.\r\n")

		fmt.Printf("Creating Table...")
		if err := db.Execute("CREATE TABLE IF NOT EXISTS t1 (a INTEGER PRIMARY KEY, b)"); err != nil {
			fail(err.Error())
		}
		fmt.Printf("ok.\r\n")

		fmt.Printf("Deleting Table content...")
		if err := db.Execute("DELETE FROM t1"); err != nil {
			fail(err.Error())
		}
		fmt.Printf("ok.\r\n")

		fmt.Printf("Adding rows to Table with ExecuteArray...")
		if err := db.ExecuteArray("INSERT INTO t1 (b) VALUES (?), (?), (?), (?)", []interface{}{int(1), "text 2", 2.2, []byte("A")}); err != nil {
			fail(err.Error())
		}
		fmt.Printf("ok.\r\n")

		fmt.Printf("Select Table...")
		if res, err := db.SelectArray("SELECT * FROM t1 WHERE a >= ?", []interface{}{1}); err != nil {
			fail(err.Error())
		} else {
			res.Dump()
		}
		fmt.Printf("ok.\r\n")

		// fmt.Printf("Checking DROP DATABASE...")
		// if err := db.DropDatabase(dbname, false); err != nil { // Database, NoError
		// 	fail(err.Error())
		// }
		// fmt.Printf("ok.\r\n")

		// fmt.Printf( "Checking CLOSE CONNECTION..." )
		// if err := db.CloseConnection( "14" ); err != nil { // ConnectionID
		//  fail( err.Error() )
		// }
		// fmt.Printf( "ok.\r\n" )

		fmt.Printf("Success.\r\n")
		os.Exit(0)
	}
	fail(err.Error())
}

func fail(Message string) {
	fmt.Printf("failed: %s\r\n", Message)
	os.Exit(1)
}
