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

package sqlitecloudtest

import (
	"testing"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const testDbnameLiteral = "test-gosdk-literal-db.sqlite"

func TestLiterals(t *testing.T) {
	var db *sqlitecloud.SQCloud
	var res *sqlitecloud.Result
	var err error

	// start := time.Now()
	if db, err = sqlitecloud.Connect(testConnectionString); err != nil {
		t.Fatal("CONNECT: ", err.Error())
	}
	defer db.Close()
	// fmt.Printf("CONNECT %v\n", time.Since(start))

	// start := time.Now()
	if err := db.CreateDatabase(testDbnameLiteral, "", "", true); err != nil { // Database, Key, Encoding, NoError
		t.Fatal("CREATE DATABASE: ", err.Error())
	}
	// fmt.Printf("CREATE DATABASE %v\n", time.Since(start))

	// start := time.Now()
	if err := db.UseDatabase(testDbnameLiteral); err != nil {
		t.Fatal("USE DATABASE: ", err.Error())
	}
	// fmt.Printf("USE DATABASE %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST NULL"); { // NULL
	case err != nil:
		t.Fatal("TEST NULL: ", err.Error())
	case res == nil:
		t.Fatal("TEST NULL: nil result")
	case !res.IsNULL():
		t.Fatal("TEST NULL: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST NULL %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST STRING"); { // String Literal
	case err != nil:
		t.Fatal("TEST STRING: ", err.Error())
	case res == nil:
		t.Fatal("TEST STRING: nil result")
	case !res.IsString():
		t.Fatal("TEST STRING: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST STRING %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST ZERO_STRING"); { // String Literal
	case err != nil:
		t.Fatal("TEST ZERO_STRING: ", err.Error())
	case res == nil:
		t.Fatal("TEST ZERO_STRING: nil result")
	case !res.IsString():
		t.Fatal("TEST ZERO_STRING: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST ZERO_STRING %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST ERROR"); { // Error
	case err == nil:
		t.Fatal("TEST ERROR: Unknown command returned no error")
	}
	// fmt.Printf("TEST ERROR %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST INTEGER"); { // Integer
	case err != nil:
		t.Fatal("TEST INTEGER: ", err.Error())
	case res == nil:
		t.Fatal("TEST INTEGER: nil result")
	case !res.IsInteger():
		t.Fatal("TEST INTEGER: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST INTEGER %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST FLOAT"); { // Float
	case err != nil:
		t.Fatal("TEST FLOAT: ", err.Error())
	case res == nil:
		t.Fatal("TEST FLOAT: nil result")
	case !res.IsFloat():
		t.Fatal("TEST FLOAT: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST FLOAT %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST BLOB"); { // BLOB
	case err != nil:
		t.Fatal("TEST BLOB: ", err.Error())
	case res == nil:
		t.Fatal("TEST BLOB: nil result")
	case !res.IsBLOB():
		t.Fatal("TEST BLOB: invalid type")
	default:
		if l := len(res.GetBuffer()); l != 1000 {
			t.Fatalf("TEST BLOB: invalid blob len (Expected 1000, Got %d)", l)
		}
	}
	res.Free()
	// fmt.Printf("TEST BLOB %v\n", time.Since(start))

	// TEST ROWSET_CHUNK is slow (i.e. 18s, because it sends 147 separated chunks, one for each row)
	// start := time.Now()
	switch res, err = db.Select("TEST ROWSET_CHUNK"); { // ROWSET
	case err != nil:
		t.Fatal("TEST ROWSET_CHUNK: ", err.Error())
	case res == nil:
		t.Fatal("TEST ROWSET_CHUNK: nil result")
	case !res.IsRowSet():
		t.Fatal("TEST ROWSET_CHUNK: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST ROWSET_CHUNK %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST ROWSET"); { // ROWSET
	case err != nil:
		t.Fatal("TEST ROWSET: ", err.Error())
	case res == nil:
		t.Fatal("TEST ROWSET: nil result")
	case !res.IsRowSet():
		t.Fatal("TEST ROWSET: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST ROWSET %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST JSON"); { // JSON
	case err != nil:
		t.Fatal("TEST JSON: ", err.Error())
	case res == nil:
		t.Fatal("TEST JSON: nil result")
	case !res.IsJSON():
		t.Fatal("TEST JSON: invalid type")
	}
	res.Free()
	// fmt.Printf("TEST JSON %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST COMMAND"); { // Command
	case err != nil:
		t.Fatal("TEST COMMAND: ", err.Error())
	case res == nil:
		t.Fatal("TEST COMMAND: nil result")
	case !res.IsCommand():
		t.Fatal("TEST COMMAND: invalid type")
	default:
		if res.GetString_() != "PING" { // should be ping
			t.Fatalf("TEST COMMAND: invalid command (%s)", res.GetString_())
		}
	}
	res.Free()
	// fmt.Printf("TEST COMMAND %v\n", time.Since(start))

	// start := time.Now()
	switch res, err = db.Select("TEST ARRAY"); { // ARRAY
	case err != nil:
		t.Fatal("TEST ARRAY: ", err.Error())
	case res == nil:
		t.Fatal("TEST ARRAY: nil result")
	case !res.IsArray():
		t.Fatal("TEST ARRAY: invalid type")
	default:
		if res.GetNumberOfRows() != 5 {
			t.Fatalf("TEST ARRAY: invalid number of rows (Expected 5, Got %d)", res.GetNumberOfRows())
		}
	}
	res.Free()
	// fmt.Printf("TEST ARRAY %v\n", time.Since(start))

	// start := time.Now()
	if err := db.UnuseDatabase(); err != nil {
		t.Fatal("UNUSE DATABASE: ", err.Error())
	}
	// fmt.Printf("UNUSE DATABASE %v\n", time.Since(start))

	// start := time.Now()
	if err := db.DropDatabase(testDbnameLiteral, false); err != nil {
		t.Fatal("DROP DATABASE: ", err.Error())
	}
	// fmt.Printf("DROP DATABASE %v\n", time.Since(start))
}
