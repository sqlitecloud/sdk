package sqlitecloudtest

const testConnectionString = "sqlitecloud://admin:admin@dev1.sqlitecloud.io:9960"
const testUsername = "admin"
const testPassword = "admin"

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// const testConnectionString = "sqlitecloud://admin:admin@localhost:8860"
