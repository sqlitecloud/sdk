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

package sqlitecloudtest

import (
	"testing"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const testDbnameSelectArray = "test-gosdk-selectarray-db.sqlite"

func TestSelectArray(t *testing.T) {
	// Server API test

	config, err1 := sqlitecloud.ParseConnectionString(testConnectionString)
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	db := sqlitecloud.New(*config)
	err := db.Connect()

	if err != nil {
		t.Fatalf(err.Error())
	}

	defer db.Close()

	// test select null string, it was causing a crash on the server
	if res, err := db.Select(""); err != nil {
		t.Fatal(err.Error())
	} else if !res.IsNULL() {
		t.Fatalf("Expected NULL, got %v", res.GetType())
	}

	// Checking CREATE DATABASE
	if err := db.ExecuteArray("CREATE DATABASE ? PAGESIZE ? IF NOT EXISTS", []interface{}{testDbnameSelectArray, 4096}); err != nil {
		t.Fatal(err.Error())
	}

	// Checking USE DATABASE
	if err := db.UseDatabase(testDbnameSelectArray); err != nil { // Database
		t.Fatal(err.Error())
	}

	// Creating Table
	if err := db.Execute("CREATE TABLE IF NOT EXISTS t1 (a INTEGER PRIMARY KEY, b)"); err != nil {
		t.Fatal(err.Error())
	}

	// Deleting Table content
	if err := db.Execute("DELETE FROM t1"); err != nil {
		t.Fatal(err.Error())
	}

	// Adding rows to Table with ExecuteArray
	if err := db.ExecuteArray("INSERT INTO t1 (b) VALUES (?), (?), (?), (?)", []interface{}{int(1), "text 2", 2.2, []byte("A")}); err != nil {
		t.Fatal(err.Error())
	}

	// Select Table
	if res, err := db.SelectArray("SELECT * FROM t1 WHERE a >= ?", []interface{}{1}); err != nil {
		t.Fatal(err.Error())
	} else if res.GetNumberOfRows() != 4 {
		t.Fatalf("Expected 4 rows, got %d", res.GetNumberOfRows())
	} else if res.GetNumberOfColumns() != 2 {
		t.Fatalf("Expected 2 columns, got %d", res.GetNumberOfColumns())
	} else if s, _ := res.GetStringValue(0, 1); s != "1" {
		t.Fatalf("Expected '1', got '%s'", s)
	} else if s, _ := res.GetStringValue(1, 1); s != "text 2" {
		t.Fatalf("Expected 'text 2', got '%s'", s)
	} else if s, _ := res.GetStringValue(2, 1); s != "2.2" {
		t.Fatalf("Expected '2.2', got '%s'", s)
	} else if s, _ := res.GetStringValue(3, 1); s != "A" {
		t.Fatalf("Expected 'A', got '%s'", s)
	}

	// Checking Unuse Database
	if err := db.UnuseDatabase(); err != nil {
		t.Fatal(err.Error())
	}

	// Checking DROP DATABASE
	if err := db.DropDatabase(testDbnameSelectArray, false); err != nil {
		t.Fatal(err.Error())
	}
}
