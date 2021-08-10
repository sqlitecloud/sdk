package main

import "os"
import "fmt"
//import "strings"
import "sqlitecloud"
import "encoding/json"


func main() {
	fmt.Printf( "Simple API test...\r\n")

	db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" )

	//db := sqlitecloud.New()
	//err := db.Connect( "dev1.sqlitecloud.io", 8860, "", "", "X", 10, 0 )


	if err == nil {
		defer db.Close()


		if db.Ping() == nil {
			fmt.Println( "PONG." )
		}

		db.UseDatabase( "X" )

		commands, _ := db.ListCommands()
		jCommands, _ := json.Marshal( commands )
		fmt.Printf( "LIST COMMANDS:\r\n%v\r\n", string( jCommands ) )
	
		info, _ := db.ListInfo()
		jInfo, _ := json.Marshal( info )
		fmt.Printf( "LIST INFO:\r\n%v\r\n", string( jInfo ) )

		tables, _ := db.ListTables()
		jTables, _ := json.Marshal( tables )
		fmt.Printf( "LIST TABLES:\r\n%v\r\n", string( jTables ) )

		fmt.Printf( "UUID = '%s'\r\n", db.GetUUID() )

		
		res, err := db.Select( "SELECT * FROM Dummy" )
		if( err == nil ) {
			fmt.Printf( "Num Rows = %d\r\n", res.GetNumberOfRows() )
			res.Dump()
			res.Free()
			os.Stderr.WriteString("\r\n")
			return
		}
		panic( err )
	}
}