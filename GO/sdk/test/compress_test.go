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

package sqlitecloudtest

import (
	"fmt"
	"math/rand"
	"testing"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const testDbnameCompress = "test-gosdk-compress-db.sqlite"
const testCompressArg = "compress=LZ4"

func TestCompress(t *testing.T) {
	connstring := fmt.Sprintf("%s/?%s", testConnectionString, testCompressArg)
	// log.Printf("TestCompress %s", connstring)

	db, err := sqlitecloud.Connect(connstring)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	if err = db.CreateDatabase(testDbnameCompress, "", "UTF-8", true); err != nil {
		t.Fatal(err)
	}

	if err := db.UseDatabase(testDbnameCompress); err != nil {
		t.Fatal(err.Error())
	}

	if err = db.Execute(`CREATE TABLE IF NOT EXISTS "TestCompress" (ID INTEGER PRIMARY KEY AUTOINCREMENT, Dummy TEXT(200) )`); err != nil {
		t.Fatal(err)
	}

	if err = db.Execute(`DELETE FROM TestCompress`); err != nil {
		t.Fatal(err)
	}

	if res, err := db.Select("GET CLIENT KEY COMPRESSION"); err != nil {
		t.Fatal(err)
	} else if res.GetString_() != "1" {
		res.Dump()
		t.Fatalf("Expected COMPRESSION = 1, got %s", res.GetString_())
	}

	rndStr := ""
	for i := 0; i < 10; i++ {
		rndStr = fmt.Sprintf("%s%d", rndStr, rand.Int())
	}

	nrows := 1000
	sql := "BEGIN; "
	for rows := 0; rows < nrows; rows++ {
		sql = sql + fmt.Sprintf("INSERT INTO TestCompress ( Dummy ) VALUES( '%s' ); ", rndStr)
	}
	sql = sql + "COMMIT;"

	// log.Printf("Execute INSERT: start")
	if err = db.Execute(sql); err != nil {
		t.Fatal(err)
	}
	// log.Printf("Execute INSERT: end")

	if err = db.Compress("lz4"); err != nil {
		t.Fatal(err)
	}

	res1, err := db.Select("SELECT * FROM TestCompress")
	if err != nil {
		t.Fatal(err)
	} else if !res1.IsRowSet() {
		t.Fatalf("Expected RowSet, got %v", res1.GetType())
	} else if res1.GetNumberOfRows() != uint64(nrows) {
		t.Fatalf("Expected %d rows, got %d", nrows, res1.GetNumberOfRows())
	}

	if err = db.Compress("NO"); err != nil {
		t.Fatal(err)
	}

	if res2, err := db.Select("SELECT * FROM TestCompress"); err != nil {
		t.Fatal(err)
	} else if !res2.IsRowSet() {
		t.Fatalf("Expected RowSet, got %v", res2.GetType())
	} else if res2.GetNumberOfRows() != uint64(nrows) {
		t.Fatalf("Expected %d rows, got %d", nrows, res2.GetNumberOfRows())
	} else if res1.GetStringValue_(uint64(nrows)-1, 1) != res2.GetStringValue_(uint64(nrows)-1, 1) {
		t.Fatalf("Expected value '%s', got '%s'", res1.GetStringValue_(uint64(nrows), 1), res2.GetStringValue_(uint64(nrows), 1))
	}

	// Checking Unuse Database
	if err := db.UnuseDatabase(); err != nil {
		t.Fatal(err.Error())
	}

	// Checking DROP DATABASE
	if err := db.DropDatabase(testDbnameCompress, false); err != nil {
		t.Fatal(err.Error())
	}
}
