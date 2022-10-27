//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.1
//     //             ///   ///  ///    Date        : 2021/08/31
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Go Methods related to the
//   ////                ///  ///                     SQCloud class for using
//     ////     //////////   ///                      the internal server
//        ////            ////                        commands.
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SQCloudConnection struct {
	ClientID       int64
	Address        string
	Username       string
	Database       string
	ConnectionDate time.Time
	LastActivity   time.Time
}

type SQCloudInfo struct {
	SQLiteVersion    string
	SQCloudVersion   string
	SQCloudBuildDate time.Time
	SQCloudGitHash   string

	ServerTime         time.Time
	ServerCPUs         int
	ServerOS           string
	ServerArchitecture string

	ServicePID         int
	ServiceStart       time.Time
	ServicePort        int
	SericeMultiplexAPI string
}

type SQCloudPlugin struct {
	Name        string
	Type        string
	Version     string
	Copyright   string
	Description string
}

// Node Functions

// AddNode - INTERNAL SERVER COMMAND: Adds a node to the SQLite Cloud Database Cluster.
func (this *SQCloud) AddNode(Node string, Address string, Cluster string, Snapshot string, Learner bool) error {
	sql := "ADD"
	if Learner {
		sql += " LEARNER"
	}
	sql += " NODE ? ADDRESS ? CLUSTER ? SNAPSHOT ?"
	return this.ExecuteArray(sql, []interface{}{Node, Address, Cluster, Snapshot} )
}

// RemoveNode - INTERNAL SERVER COMMAND: Removes a node to the SQLite Cloud Database Cluster.
func (this *SQCloud) RemoveNode(Node string) error {
	return this.ExecuteArray("REMOVE NODE ?", []interface{}{Node})
}

// RemoveNode - INTERNAL SERVER COMMAND: Lists all nodes of this SQLite Cloud Database Cluster.
func (this *SQCloud) ListNodes() ([]string, error) {
	return this.SelectStringList("LIST NODES")
}

// Connection Functions

// CloseConnection - INTERNAL SERVER COMMAND: Closes the specified connection.
func (this *SQCloud) CloseConnection(ConnectionID string) error {
	return this.ExecuteArray("CLOSE CONNECTION ?", []interface{}{ConnectionID})
}

// ListConnections - INTERNAL SERVER COMMAND: Lists all connections of this SQLite Cloud Database Cluster.
func (this *SQCloud) ListConnections() ([]SQCloudConnection, error) {
	connectionList := []SQCloudConnection{}
	result, err := this.Select("LIST CONNECTIONS")
	if err == nil {
		if result != nil {
			defer result.Free()
			if result.GetNumberOfColumns() == 6 {
				for row, rows := uint64(0), result.GetNumberOfRows(); row < rows; row++ {
					connection := SQCloudConnection{}
					connection.ClientID, _ = result.GetInt64Value(row, 1)
					connection.Address, _ = result.GetStringValue(row, 2)
					connection.Username, _ = result.GetStringValue(row, 3)
					connection.Database, _ = result.GetStringValue(row, 4)
					connection.ConnectionDate, _ = result.GetSQLDateTime(row, 4)
					connection.LastActivity, _ = result.GetSQLDateTime(row, 6)
					connectionList = append(connectionList, connection)
				}
				return connectionList, nil
			}
			return []SQCloudConnection{}, errors.New("ERROR: Query returned not 6 Columns (-1)")
		}
		return []SQCloudConnection{}, errors.New("ERROR: Query returned no result (-1)")
	}
	return []SQCloudConnection{}, err
}

// ListDatabaseConnections - INTERNAL SERVER COMMAND: Lists all connections that use the specified Database on this SQLite Cloud Database Cluster.
func (this *SQCloud) ListDatabaseConnections(Database string) ([]SQCloudConnection, error) {
	connectionList := []SQCloudConnection{}
	result, err := this.SelectArray("LIST DATABASE CONNECTIONS ?", []interface{}{Database})
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 2 {
				for row, rows := uint64(0), result.GetNumberOfRows(); row < rows; row++ {
					connection := SQCloudConnection{}
					connection.ClientID, _ = result.GetInt64Value(row, 1)
					connection.Address = ""
					connection.Username = ""
					connection.Database = Database
					connection.ConnectionDate = time.Unix(0, 0)
					connection.LastActivity, _ = result.GetSQLDateTime(row, 2)
					connectionList = append(connectionList, connection)
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New("ERROR: Query returned not 2 Columns (-1)")
		}
		return []SQCloudConnection{}, errors.New("ERROR: Query returned no result (-1)")
	}
	return []SQCloudConnection{}, err
}

// ListDatabaseClientConnectionIds - INTERNAL SERVER COMMAND: Lists all connections with the specified DatabaseId on this SQLite Cloud Database Cluster.
// GUB(andreas) Start index = 1
func (this *SQCloud) ListDatabaseClientConnectionIds(DatabaseID uint64) ([]SQCloudConnection, error) {
	connectionList := []SQCloudConnection{}
	result, err := this.Select(fmt.Sprintf("LIST DATABASE CONNECTIONS ID %d", DatabaseID))
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 2 {
				for row, rows := uint64(1), result.GetNumberOfRows(); row < rows; row++ {
					connection := SQCloudConnection{}
					connection.ClientID, _ = result.GetInt64Value(row, 1)
					connection.Address = ""
					connection.Username = ""
					connection.Database = ""
					connection.ConnectionDate = time.Unix(0, 0)
					connection.LastActivity, _ = result.GetSQLDateTime(row, 2)
					connectionList = append(connectionList, connection)
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New("ERROR: Query returned not 2 Columns (-1)")
		}
		return []SQCloudConnection{}, errors.New("ERROR: Query returned no result (-1)")
	}
	return []SQCloudConnection{}, err
}

// Auth Functions

func authCommand(Username string, Password string) (string, []interface{}) {
	return "AUTH USER ? PASSWORD ?;", []interface{}{Username, Password}
}

// Auth - INTERNAL SERVER COMMAND: Authenticates User with the given credentials.
func (this *SQCloud) Auth(Username string, Password string) error {
	return this.ExecuteArray(authCommand(Username, Password))
}

func authWithKeyCommand(Key string) (string, []interface{}) {
	return "AUTH APIKEY ?;", []interface{}{Key}
}

// Auth - INTERNAL SERVER COMMAND: Authenticates User with the given API KEY.
func (this *SQCloud) AuthWithKey(Key string) error {
	return this.ExecuteArray(authWithKeyCommand(Key))
}

// Database funcitons

// CreateDatabase - INTERNAL SERVER COMMAND: Creates a new Database on this SQLite Cloud Database Cluster.
// If the Database already exists on this Database Server, an error is returned except the NoError flag is set.
// Encoding specifies the character set Encoding that should be used for the new Database - for example "UFT-8".
func (this *SQCloud) CreateDatabase(Database string, Key string, Encoding string, NoError bool) error {
	sql := "CREATE DATABASE ?"
	args := []interface{}{Database}
	if strings.TrimSpace(Key) != "" {
		sql += " KEY ?"
		args = append(args, Key)
	}
	if strings.TrimSpace(Encoding) != "" {
		sql += " ENCODING ?"
		args = append(args, Encoding)
	}
	if NoError {
		sql += " IF NOT EXISTS"
	}
	// println( sql )
	return this.ExecuteArray(sql, args)
}

// DropDatabase - INTERNAL SERVER COMMAND: Deletes the specified Database on this SQLite Cloud Database Cluster.
// If the given Database is not present on this Database Server or the user has not the necessary access rights,
// an error describing the problem will be returned.
// If the NoError flag is set, no error will be reported if the database does not exist.
func (this *SQCloud) DropDatabase(Database string, NoError bool) error {
	sql := "DROP DATABASE ?"
	if NoError {
		sql += " IF EXISTS"
	}
	return this.ExecuteArray(sql, []interface{}{Database} )
}

// ListDatabases - INTERNAL SERVER COMMAND: Lists all Databases that are present on this SQLite Cloud Database Cluster and returns the Names of the databases in an array of strings.
func (this *SQCloud) ListDatabases() ([]string, error) {
	return this.SelectStringList("LIST DATABASES")
}

// GetDatabase - INTERNAL SERVER COMMAND: Gets the name of the previously selected Database as string. (see: *SQCloud.UseDatabase())
// If no database was selected, an error describing the problem is returned.
func (this *SQCloud) GetDatabase() (string, error) {
	result, err := this.Select("GET DATABASE")
	if result != nil {
		defer result.Free()
		if err != nil {
			return "", err
		}
		return result.GetString()
	}
	return "", err
}

// GetDatabaseID - INTERNAL SERVER COMMAND: Gets the ID of the previously selected Database as uint64. (see: *SQCloud.UseDatabase())
// If no database was selected, an error describing the problem is returned.
func (this *SQCloud) GetDatabaseID() (uint64, error) {
	result, err := this.Select("GET DATABASE ID")
	if result != nil {
		defer result.Free()
		if err != nil {
			return 0, err
		}
		val, err := result.GetInt64()
		return uint64(val), err
	}
	return 0, err
}

func useDatabaseCommand(Database string) (string, []interface{}) {
	return "USE DATABASE ?;", []interface{}{Database}
}

// UseDatabase - INTERNAL SERVER COMMAND: Selects the specified Database for usage.
// Only if a database was selected, SQL Commands can be sent to this specific Database.
// An error is returned if the specified Database was not found or the user has not the necessary access rights to work with this Database.
func (this *SQCloud) UseDatabase(Database string) error {
	this.Database = Database
	return this.ExecuteArray(useDatabaseCommand(Database))
}

// UseDatabase - INTERNAL SERVER COMMAND: Releases the actual Database.
// Any further SQL commands will result in an error before selecting a new Database. (see: *SQCloud.UseDatabase())
func (this *SQCloud) UnuseDatabase() error {
	this.Database = ""
	return this.Execute("UNUSE DATABASE")
}

// Plugin Functions

// EnablePlugin enables the SQLite Plugin on the SQlite Cloud Database server.
func (this *SQCloud) EnablePlugin(Plugin string) error {
	return this.ExecuteArray("ENABLED PLUGIN ?", []interface{}{Plugin})
}

// DisablePlugin disables the SQLite Plugin on the SQlite Cloud Database server.
func (this *SQCloud) DisablePlugin(Plugin string) error {
	return this.ExecuteArray("DISABLE PLUGIN ?", []interface{}{Plugin})
}

// ListPlugins list all available Plugins at the SQlite Cloud Database server and returns an array of SQCloudPlugin.
func (this *SQCloud) ListPlugins() ([]SQCloudPlugin, error) {
	pluginList := []SQCloudPlugin{}
	result, err := this.Select("LIST PLUGINS")
	if err == nil {
		if result != nil {
			for row, rows := uint64(1), result.GetNumberOfRows(); row < rows; row++ {
				if result.GetNumberOfColumns() == 5 {
					plugin := SQCloudPlugin{}
					plugin.Name, _ = result.GetStringValue(row, 0)
					plugin.Type, _ = result.GetStringValue(row, 1)
					plugin.Version, _ = result.GetStringValue(row, 2)
					plugin.Copyright, _ = result.GetStringValue(row, 3)
					plugin.Description, _ = result.GetStringValue(row, 4)
					pluginList = append(pluginList, plugin)

				} else {
					result.Free()
					return []SQCloudPlugin{}, errors.New("ERROR: Query returned not 5 Columns (-1)")
				}
			}
			result.Free()
			return pluginList, nil
		}
		return []SQCloudPlugin{}, nil
	}
	return []SQCloudPlugin{}, err
}

// Key / Value Pair functions

// SetKey set the provided key value pair with the key Key to the string value Value.
func (this *SQCloud) SetKey(Key string, Value string) error {
	return this.ExecuteArray("SET KEY ? TO ?", []interface{}{Key, Value})
}

// GetKey gets the Value of the key Key and returns it as a string value.
// If the Key was not found an error is returned.
// BUG(andreas): If key is not set, DB returns NULL -> does not work with current implementation
func (this *SQCloud) GetKey(Key string) (string, error) {
	result, err := this.SelectArray("GET KEY ?", []interface{}{Key})
	if result != nil {
		defer result.Free()
		if err != nil {
			return "", err
		}
		return result.GetString()
	}
	return "", err
}

// DropKey deletes the key value pair referenced with Key.
// If the Key does not exists, no error is returned.
func (this *SQCloud) DropKey(Key string) error {
	return this.ExecuteArray("DROP KEY ?", []interface{}{Key})
}

// ListKeys lists all key value pairs on the server and returns an array of SQCloudKeyValues.
func (this *SQCloud) ListKeys() ([]SQCloudKeyValues, error) {
	return this.SelectKeyValues("LIST KEYS")
}

// ListClientKeys lists all client/connection specific keys and values and returns the data in an array of type SQCloudKeyValues.
func (this *SQCloud) ListClientKeys() ([]SQCloudKeyValues, error) {
	return this.SelectKeyValues("LIST CLIENT KEYS")
}

// ListDatabaseKeys lists all server specific keys and values and returns an array of type SQCloudKeyValues.
func (this *SQCloud) ListDatabaseKeys() ([]SQCloudKeyValues, error) {
	return this.SelectKeyValues("LIST DATABASE KEYS")
}

/// Misc functions

// Ping sends the PING command to the SQLite Cloud Database Server and returns nil if it got a PONG answer.
// If no PONG was received or a timeout occurred, an error describing the problem is retuned.
func (this *SQCloud) Ping() error {
	if result, err := this.Select("PING"); result != nil {
		defer result.Free()

		if err == nil {
			if retval, err := result.GetString(); retval == "PONG" {
				return err // should be nil on success...
			} else {
				return errors.New("ERROR: Unexpected result (-1)")
			}
		}
		return err
	} else {
		if err != nil {
			return err
		}
		return errors.New("Got no result on Ping")
	}
}

// ListCommands lists all available server commands and returns them in an array of strings.
func (this *SQCloud) ListCommands() ([]string, error) {
	return this.SelectStringList("LIST COMMANDS")
}

// GetInfo fetches all SQLite Cloud Database server specific runtime informations and returns a SQCloudInfo structure.
func (this *SQCloud) GetInfo() (SQCloudInfo, error) {
	info := SQCloudInfo{
		SQLiteVersion:      "0.0.0",
		SQCloudVersion:     "0.0.0",
		SQCloudBuildDate:   time.Unix(0, 0),
		SQCloudGitHash:     "N/A",
		ServerTime:         time.Unix(0, 0),
		ServerCPUs:         0,
		ServerOS:           "N/A",
		ServerArchitecture: "N/A",
		ServicePID:         0,
		ServiceStart:       time.Unix(0, 0),
		ServicePort:        0,
		SericeMultiplexAPI: "N/A",
	}

	result, err := this.SelectKeyValues("LIST INFO")
	if err == nil {
		for _, pair := range result {
			switch pair.Key {
			case "sqlitecloud_version":
				info.SQCloudVersion = pair.Value
			case "sqlite_version":
				info.SQLiteVersion = pair.Value
			case "sqlitecloud_build_date":
				info.SQCloudBuildDate, _ = time.Parse("Jan 2 2006", pair.Value)
			case "sqlitecloud_git_hash":
				info.SQCloudGitHash = pair.Value
			case "os":
				info.ServerOS = pair.Value
			case "arch_bits":
				info.ServerArchitecture = pair.Value
			case "multiplexing_api":
				info.SericeMultiplexAPI = pair.Value
			case "tcp_port":
				info.ServicePort, _ = strconv.Atoi(pair.Value)
			case "process_id":
				info.ServicePID, _ = strconv.Atoi(pair.Value)
			case "num_processors":
				info.ServerCPUs, _ = strconv.Atoi(pair.Value)
			case "startup_datetime":
				info.ServiceStart, _ = time.Parse("2006-01-02 15:04:05", pair.Value)
			case "current_datetime":
				info.ServerTime, _ = time.Parse("2006-01-02 15:04:05", pair.Value)
			default:
			}
		}
	}
	return info, err
}

// ListTables lists all tables in the selected database and returns them in an array of strings.
// If no database was selected with SQCloud.UseDatabase(), an error is returned.
func (this *SQCloud) ListTables() ([]string, error) {
	return this.SelectStringList("LIST TABLES")
}
