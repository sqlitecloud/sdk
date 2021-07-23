package main

import "os"
import "fmt"
import "sqlitecloud"


func main() {
	fmt.Printf( "Simple API test...\r\n")

	db, err := sqlitecloud.Connect( "***REMOVED***", 8860, "", "", "Test", 10, 0 )
	if err == nil {
		defer db.Close()

		fmt.Printf( "UUID = '%s'\r\n", db.GetCloudUUID() )

		db.Use( "Test" )
		res, err := db.Execute( "SELECT * FROM Dummy" )
		if( err == nil ) {
			fmt.Printf( "Num Rows = %d\r\n", res.GetNumRows() )
			res.Dump( 80 )
			res.Free()
			os.Stderr.WriteString("\r\n")
			return
		}
		panic( err )
	}
}