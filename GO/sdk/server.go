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

type SQCloudNodeStatus int64

const (
	Leader SQCloudNodeStatus = iota
	Follower
	Candidate
	Learner
)

type SQCloudNodeProgress int64

const (
	Probe SQCloudNodeProgress = iota
	Replicate
	Snapshot
	Unknown
)

type SQCloudNode struct {
	NodeID           int64
	NodeInterface    string
	ClusterInterface string
	Status           SQCloudNodeStatus
	Progress         SQCloudNodeProgress
	Match            int64
	LastActivity     time.Time
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
	ServiceNocluster   int
	ServiceNodeID      int
	SericeMultiplexAPI string

	TLS                   string
	TLSConnVersion        string
	TLSConnCipher         string
	TLSConnCipherStrength int
	TLSConnAlpnSelected   string
	TLSConnServername     string
	TLSPeerCertProvided   int
	TLSPeerCertSubject    string
	TLSPeerCertIssuer     string
	TLSPeerCertHash       string
	TLSPeerCertNotBefore  time.Time
	TLSPeerCertNotAfter   time.Time
}

type SQCloudPlugin struct {
	Name        string
	Type        string
	Enabled     bool
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
	return this.ExecuteArray(sql, []interface{}{Node, Address, Cluster, Snapshot})
}

// RemoveNode - INTERNAL SERVER COMMAND: Removes a node to the SQLite Cloud Database Cluster.
func (this *SQCloud) RemoveNode(Node string) error {
	return this.ExecuteArray("REMOVE NODE ?", []interface{}{Node})
}

func stringToSQCloudNodeStatus(s string) (SQCloudNodeStatus, error) {
	switch strings.ToLower(s) {
	case "leader":
		return Leader, nil
	case "follower":
		return Follower, nil
	case "candidate":
		return Candidate, nil
	case "learner":
		return Learner, nil
	default:
		return -1, fmt.Errorf("Cannot convert '%s' to SQCloudNodeStatus", s)
	}
}

func stringToSQCloudNodeProgress(s string) (SQCloudNodeProgress, error) {
	switch strings.ToLower(s) {
	case "probe":
		return Probe, nil
	case "replicate":
		return Replicate, nil
	case "candidate":
		return Snapshot, nil
	case "unknown":
		return Unknown, nil
	default:
		return -1, fmt.Errorf("Cannot convert '%s' to SQCloudNodeProgress", s)
	}
}

// RemoveNode - INTERNAL SERVER COMMAND: Lists all nodes of this SQLite Cloud Database Cluster.
func (this *SQCloud) ListNodes() ([]SQCloudNode, error) {
	list := []SQCloudNode{}
	result, err := this.Select("LIST NODES")
	if err == nil {
		if result != nil {
			defer result.Free()
			if result.GetNumberOfColumns() == 7 {
				for row, rows := uint64(0), result.GetNumberOfRows(); row < rows; row++ {
					node := SQCloudNode{}
					node.NodeID, _ = result.GetInt64Value(row, 0)
					node.NodeInterface, _ = result.GetStringValue(row, 1)
					node.ClusterInterface, _ = result.GetStringValue(row, 2)
					node.Status, _ = stringToSQCloudNodeStatus(result.GetStringValue_(row, 3))
					node.Progress, _ = stringToSQCloudNodeProgress(result.GetStringValue_(row, 4))
					node.Match, _ = result.GetInt64Value(row, 5)
					node.LastActivity, _ = result.GetSQLDateTime(row, 6)
					list = append(list, node)
				}
				return list, nil
			}
			return []SQCloudNode{}, errors.New("ERROR: Query returned not 7 Columns (-1)")
		}
		return []SQCloudNode{}, nil
	}
	return []SQCloudNode{}, err
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
					connection.ClientID, _ = result.GetInt64Value(row, 0)
					connection.Address, _ = result.GetStringValue(row, 1)
					connection.Username, _ = result.GetStringValue(row, 2)
					connection.Database, _ = result.GetStringValue(row, 3)
					connection.ConnectionDate, _ = result.GetSQLDateTime(row, 4)
					connection.LastActivity, _ = result.GetSQLDateTime(row, 5)
					connectionList = append(connectionList, connection)
				}
				return connectionList, nil
			}
			return []SQCloudConnection{}, errors.New("ERROR: Query returned not 6 Columns (-1)")
		}
		return []SQCloudConnection{}, nil
	}
	return []SQCloudConnection{}, err
}

func resultToConnectionList(result *Result, err error) ([]SQCloudConnection, error) {
	connectionList := []SQCloudConnection{}
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 6 {
				for row, rows := uint64(0), result.GetNumberOfRows(); row < rows; row++ {
					connection := SQCloudConnection{}
					connection.ClientID = result.GetInt64Value_(row, 0)
					connection.Address = result.GetStringValue_(row, 1)
					connection.Username = result.GetStringValue_(row, 2)
					connection.Database = result.GetStringValue_(row, 3)
					connection.ConnectionDate = result.GetSQLDateTime_(row, 4)
					connection.LastActivity = result.GetSQLDateTime_(row, 5)
					connectionList = append(connectionList, connection)
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New("ERROR: Query returned not 6 Columns (-1)")
		}
		return []SQCloudConnection{}, nil
	}
	return []SQCloudConnection{}, err
}

// ListDatabaseConnections - INTERNAL SERVER COMMAND: Lists all connections that use the specified Database on this SQLite Cloud Database Cluster.
func (this *SQCloud) ListDatabaseConnections(Database string) ([]SQCloudConnection, error) {
	result, err := this.SelectArray("LIST DATABASE CONNECTIONS ?", []interface{}{Database})
	return resultToConnectionList(result, err)
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

// RemoveDatabase - INTERNAL SERVER COMMAND: Deletes the specified Database on this SQLite Cloud Database Cluster.
// If the given Database is not present on this Database Server or the user has not the necessary access rights,
// an error describing the problem will be returned.
// If the NoError flag is set, no error will be reported if the database does not exist.
func (this *SQCloud) RemoveDatabase(Database string, NoError bool) error {
	sql := "REMOVE DATABASE ?"
	if NoError {
		sql += " IF EXISTS"
	}
	return this.ExecuteArray(sql, []interface{}{Database})
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
				if result.GetNumberOfColumns() == 6 {
					plugin := SQCloudPlugin{}
					plugin.Name, _ = result.GetStringValue(row, 0)
					plugin.Type, _ = result.GetStringValue(row, 1)
					plugin.Enabled = result.GetInt32Value_(row, 2) != 0
					plugin.Version, _ = result.GetStringValue(row, 3)
					plugin.Copyright, _ = result.GetStringValue(row, 4)
					plugin.Description, _ = result.GetStringValue(row, 5)
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

// RemoveKey deletes the key value pair referenced with Key.
// If the Key does not exists, no error is returned.
func (this *SQCloud) RemoveKey(Key string) error {
	return this.ExecuteArray("REMOVE KEY ?", []interface{}{Key})
}

// ListKeys lists all key value pairs on the server and returns an array of SQCloudKeyValues.
func (this *SQCloud) ListKeys() (map[string]string, error) {
	return this.SelectKeyValues("LIST KEYS")
}

// ListClientKeys lists all client/connection specific keys and values and returns the data in an array of type SQCloudKeyValues.
func (this *SQCloud) ListClientKeys() (map[string]string, error) {
	return this.SelectKeyValues("LIST CLIENT KEYS")
}

// ListDatabaseKeys lists all server specific keys and values and returns an array of type SQCloudKeyValues.
func (this *SQCloud) ListDatabaseKeys(Database string) (map[string]string, error) {
	return this.SelectArrayKeyValues("LIST DATABASE ? KEYS", []interface{}{Database})
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
	//fmt.Printf("Result %v", result)
	if err == nil {
		for k, v := range result {
			switch k {
			case "sqlitecloud_version":
				info.SQCloudVersion = v
			case "sqlite_version":
				info.SQLiteVersion = v
			case "sqlitecloud_build_date":
				info.SQCloudBuildDate, _ = time.Parse("Jan 2 2006", v)
			case "sqlitecloud_git_hash":
				info.SQCloudGitHash = v
			case "os":
				info.ServerOS = v
			case "arch_bits":
				info.ServerArchitecture = v
			case "multiplexing_api":
				info.SericeMultiplexAPI = v
			case "listening_port":
				info.ServicePort, _ = strconv.Atoi(v)
			case "process_id":
				info.ServicePID, _ = strconv.Atoi(v)
			case "num_processors":
				info.ServerCPUs, _ = strconv.Atoi(v)
			case "startup_datetime":
				info.ServiceStart, _ = time.Parse("2006-01-02 15:04:05", v)
			case "current_datetime":
				info.ServerTime, _ = time.Parse("2006-01-02 15:04:05", v)
			case "nocluster":
				info.ServiceNocluster, _ = strconv.Atoi(v)
			case "nodeid":
				info.ServiceNodeID, _ = strconv.Atoi(v)
			case "tls":
				info.TLS = v
			case "tls_conn_version":
				info.TLSConnVersion = v
			case "tls_conn_cipher":
				info.TLSConnCipher = v
			case "tls_conn_cipher_strength":
				info.TLSConnCipherStrength, _ = strconv.Atoi(v)
			case "tls_conn_alpn_selected":
				info.TLSConnAlpnSelected = v
			case "tls_conn_servername":
				info.TLSConnServername = v
			case "tls_peer_cert_provided":
				info.TLSConnCipherStrength, _ = strconv.Atoi(v)
			case "tls_peer_cert_subject":
				info.TLSPeerCertSubject = v
			case "tls_peer_cert_issuer":
				info.TLSPeerCertIssuer = v
			case "tls_peer_cert_hash":
				info.TLSPeerCertHash = v
			case "tls_peer_cert_notbefore":
				info.TLSPeerCertNotBefore, _ = time.Parse("2006-01-02 15:04:05", v)
			case "tls_peer_cert_notafter":
				info.TLSPeerCertNotAfter, _ = time.Parse("2006-01-02 15:04:05", v)
			default:
			}
		}
	}
	return info, err
}

// ListTables lists all tables in the selected database and returns them in an array of strings.
// If no database was selected with SQCloud.UseDatabase(), an error is returned.
func (this *SQCloud) ListTables() ([]string, error) {
	return this.SelectStringListWithCol("LIST TABLES", 1)
}
