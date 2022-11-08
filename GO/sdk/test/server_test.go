//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andreas Pfeil
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
	"fmt"
	"strings"
	"testing"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const testDbnameServer = "test-gosdk-server-db.sqlite"

func TestServer(t *testing.T) {
	db, err := sqlitecloud.Connect(testConnectionString)
	if err != nil {
		t.Fatal("Connect: ", err.Error())
	}

	// Checking wrong AUTH
	if err := db.Auth(testUsername, "wrong password"); err == nil {
		t.Fatal("AUTH: Expected authorization failed, got authorized")
	}
	db.Close()

	// reopen the connection (it was closed because of the auth command with wrong credentials)
	db, err = sqlitecloud.Connect(testConnectionString)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	// Checking PING
	if db.Ping() != nil {
		t.Fatal(err.Error())
	}

	// Checking AUTH
	if err := db.Auth(testUsername, testPassword); err != nil {
		t.Fatal("Checking AUTH: ", err.Error())
	}

	// Checking CREATE DATABASE
	if err := db.CreateDatabase(testDbnameServer, "", "", false); err != nil { // Database, Key, Encoding, NoError
		t.Fatal("CREATE DATABASE: ", err.Error())
	}

	// Checking LIST DATABASES
	if databases, err := db.ListDatabases(); err != nil {
		t.Fatal("LIST DATABASES: ", err.Error())
	} else {
		if len(databases) == 0 || !contains(databases, testDbnameServer) {
			t.Fatal("LIST DATABASES: ", fmt.Sprintf("Database %s not found in LIST DATABASES", testDbnameServer))
		}
	}

	// USE DATABASE
	if err := db.UseDatabase(testDbnameServer); err != nil { // Database
		t.Fatal("USE DATABASES: ", err.Error())
	}

	// Checking SET KEY
	testKey := "TestKey"
	testKeyValue := "1405"
	if err := db.SetKey(testKey, testKeyValue); err != nil { // Key, Value
		t.Fatal("SET KEY: ", err.Error())
	}
	testKey = strings.ToLower(testKey)

	// LIST KEYS
	if keys, err := db.ListKeys(); err != nil {
		t.Fatal("LIST KEYS: ", err.Error())
	} else {
		if _, found := keys[testKey]; !found {
			t.Fatal("LIST KEY: ", fmt.Sprintf("Key %s not found in LIST KEYS", testKey))
		}
	}

	// GET KEY
	if val, err := db.GetKey(testKey); err != nil {
		t.Fatal("GET KEY: ", err.Error())
	} else {
		if val != testKeyValue {
			t.Fatal("GET KEY: ", fmt.Sprintf("Expected value %s, Got %s.", testKeyValue, val))
		}
	}

	// DROP KEY
	if err := db.DropKey(testKey); err != nil {
		t.Fatal("DROP KEY: ", err.Error())
	} else {
		if keys, err := db.ListKeys(); err != nil {
			t.Fatal("DROP KEY, LIST KEY: ", err.Error())
		} else {
			if _, found := keys[testKey]; found {
				t.Fatal("DROP KEY, LIST KEY: ", fmt.Sprintf("Key %s still found in LIST KEYS", testKey))
			}
		}
	}

	// LIST COMMANDS
	if commands, err := db.ListCommands(); err != nil {
		t.Fatal("LIST COMMANDS: ", err.Error())
	} else {
		if len(commands) == 0 {
			t.Fatal("LIST COMMANDS: Invalid result (0 rows).")
		}
	}

	// LIST CONNECTIONS
	if connections, err := db.ListConnections(); err != nil {
		t.Fatal("LIST CONNECTIONS: ", err.Error())
	} else {
		if len(connections) == 0 {
			t.Fatal("LIST CONNECTIONS: Invalid result (no connections).")
		}
	}

	// LIST DATABASE CONNECTIONS ()...")
	if connections, err := db.ListDatabaseConnections(testDbnameServer); err != nil {
		t.Fatal("LIST DATABASE CONNECTIONS: ", err.Error())
	} else {
		if len(connections) == 0 {
			t.Fatal("LIST DATABASE CONNECTIONS: Invalid result (0 connections).")
		}
	}

	// LIST INFO
	if _, err := db.GetInfo(); err != nil {
		t.Fatal("LIST INFO: ", err.Error())
	}

	// CREATE TABLE
	testTableServer := "TestTable"
	if err := db.Execute(fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' (a INTEGER PRIMARY KEY, b)", testTableServer)); err != nil {
		t.Fatal("CREATE TABLE: ", err.Error())
	}

	// LIST TABLES
	if tables, err := db.ListTables(); err != nil {
		t.Fatal("LIST TABLES: ", err.Error())
	} else {
		if len(tables) < 1 || !contains(tables, testTableServer) {
			t.Fatal("LIST DATABASES: ", fmt.Sprintf("Table %s not found in LIST TABLES", testTableServer))
		}
	}

	// LIST PLUGINS
	if plugins, err := db.ListPlugins(); err != nil {
		t.Fatal("LIST PLUGINS: ", err.Error())
	} else if len(plugins) == 0 {
		t.Fatal("LIST DATABASES: Invalid result")
	}

	// LIST CLIENT KEYS
	if ckeys, err := db.ListClientKeys(); err != nil {
		t.Fatal("LIST CLIENT KEYS: ", err.Error())
	} else if len(ckeys) == 0 {
		t.Fatal("LIST CLIENT KEYS: Invalid result")
	}

	// LIST NODES
	if nodes, err := db.ListNodes(); err != nil {
		t.Fatal("LIST NODES: ", err.Error())
	} else {
		if len(nodes) == 0 {
			t.Fatal("LIST NODES: Ivalid result.")
		}
	}

	// LIST DATABASE KEYS
	if _, err := db.ListDatabaseKeys(testDbnameServer); err != nil {
		t.Fatal("LIST DATABASE KEYS:", err.Error())
	}

	// UNUSE DATABASE
	if err := db.UnuseDatabase(); err != nil {
		t.Fatal("UNUSE DATABASE:", err.Error())
	}

	// DROP DATABASE
	if err := db.DropDatabase(testDbnameServer, false); err != nil { // Database, NoError
		t.Fatal("DROP DATABASES: ", err.Error())
	}

	// fmt.Printf( "Checking CLOSE CONNECTION..." )
	// if err := db.CloseConnection( "14" ); err != nil { // ConnectionID
	//  fail( err.Error() )
	// }
	// fmt.Printf( "ok.\r\n" )
}
