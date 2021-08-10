package main

import "os"
import "fmt"
//import "strings"
import "sqlitecloud"
//import "encoding/json"


func main() {
	fmt.Printf( "Server API test...\r\n")

	db := sqlitecloud.New()

	err := db.Connect( "dev1.sqlitecloud.io", 8860, "", "", "X", 10, "NO", 0 ) // Host, Port, Username, Password, Database, Timeout, Compression, Family
	if err == nil {
		defer db.Close()

		fmt.Printf( "Checking PING..." )
		if db.Ping() != nil {
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking AUTH..." )
		if err := db.Auth( "pfeil", "secret"); err != nil { //  Username, Password
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )





		fmt.Printf( "Checking CREATE DATABASE..." )
		if err := db.CreateDatabase( "xyz", "", "", false ); err != nil { // Database, Key, Encoding, NoError
			fail( err.Error() )
		} 
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST DATABASES..." )
		if databases, err := db.ListDatabases(); err != nil {
			fail( err.Error() )
		} else {
			if len( databases ) == 0 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking DROP DATABASE..." )
		if err := db.DropDatabase( "xyz", false ); err != nil { // Database, NoError
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )
		
		fmt.Printf( "Checking USE DATABASE..." )
		if err := db.UseDatabase( "X" ); err != nil { // Database
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )


		fmt.Printf( "Checking SET KEY..." )
		if err := db.SetKey( "A", "1405" ); err != nil { // Key, Value
			fail( err.Error() )
		} 
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST KEYS..." )
		if keys, err := db.ListKeys() ; err != nil { 
			fail( err.Error() )
		} else {
			if len( keys ) != 1 {
				fail( "Invalid amount of keys found." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking GET KEY..." )
		if val, err := db.GetKey( "A" ) ; err != nil { 
			fail( err.Error() )
		} else {
			if val != "1405" {
				fail( "Invalid value." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking DROP KEY..." )
		if err := db.DropKey( "A" ); err != nil {
			fail( err.Error() )
		} else {
			if keys, err := db.ListKeys() ; err != nil {
				fail( err.Error() )
			} else {
				if len( keys ) != 0 {
					fail( "Key could not be deleted." )
				}
			}
			fmt.Printf( "ok.\r\n" )
		}


		


		fmt.Printf( "Checking LIST COMMANDS..." )
		if commands, err := db.ListCommands(); err != nil {
			fail( err.Error() )
		} else {
			if len( commands ) != 27 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )


		fmt.Printf( "Checking LIST CONNECTIONS..." )
		if connections, err := db.ListConnections(); err != nil {
			fail( err.Error() )
		} else {
			if len( connections ) == 0 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST DATABASE CONNECTIONS (X)..." )
		if connections, err := db.ListDatabaseConnections( "X" ); err != nil {
			fail( err.Error() )
		} else {
			if len( connections ) == 0 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST DATABASE CONNECTIONS ID(3)..." )
		if connections, err := db.ListDatabaseClientConnectionIds( 3 ); err != nil {
			fail( err.Error() )
		} else {
			if len( connections ) == 0 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST INFO..." )
		if _, err := db.GetInfo(); err != nil {
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST TABLES..." )
		if tables, err := db.ListTables(); err != nil {
			fail( err.Error() )
		} else {
			if len( tables ) != 1 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST PLUGINS..." )
		if _, err := db.ListPlugins(); err != nil {
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )
	

		fmt.Printf( "Checking LIST CLIENT KEYS..." )
		if _, err := db.ListClientKeys(); err != nil {
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )




		fmt.Printf( "Checking LIST NODES..." )
		if nodes, err := db.ListNodes(); err != nil {
			fail( err.Error() )
		} else {
			if len( nodes ) == 0 {
				fail( "Ivalid result." )
			}
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking LIST DATABASE KEYS..." )
		if _, err := db.ListDatabaseKeys(); err != nil {
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )




		fmt.Printf( "Checking UNUSE DATABASE..." )
		if err := db.UnuseDatabase(); err != nil { // Database
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )

		fmt.Printf( "Checking CLOSE CONNECTION..." )
		if err := db.CloseConnection( "14" ); err != nil { // ConnectionID
			fail( err.Error() )
		}
		fmt.Printf( "ok.\r\n" )
		

		fmt.Printf( "Success.\r\n" )
		os.Exit( 0 )
	}
	fail( err.Error() )
}

func fail( Message string ) {
	fmt.Printf( "failed: %s\r\n", Message )
	os.Exit( 1 )
}