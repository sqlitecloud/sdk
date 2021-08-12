package main

import "sqlitecloud"

func main() {
	db, err := sqlitecloud.Connect( "sqlitecloud://dev1.sqlitecloud.io/X" )
	if err == nil {
		defer db.Close()

		db.CreateDatabase( "X", "", "UTF-8", true )
		db.Execute( `CREATE TABLE IF NOT EXISTS "Dummy" (ID INTEGER PRIMARY KEY AUTOINCREMENT, FirstName TEXT(20), LastName TEXT(20), ZIP INTEGER, City TEXT, Address TEXT)` )
		db.Execute( `DELETE FROM Dummy` )
		db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Some', 'One', 96450, 'Coburg', "Mohrenstrasse 1" )` )
		db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'Someone', 'Else', 96145, 'Sesslach', 'Raiffeisenstrasse 6' )` )
		db.Execute( `INSERT INTO Dummy ( FirstName, LastName, ZIP, City, Address ) VALUES( 'One', 'More', 91099, 'Poxdorf', "Langholzstr. 4" )` )

		if res, err := db.Select( "SELECT * FROM Dummy" ); res != nil {
			defer res.Free()

			if err == nil {
				res.Dump()
				return

	} } }
	panic( err )
}