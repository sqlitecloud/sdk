package sqlitecloud

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "time"
import "errors"
import "strconv"

type SQCloudConnection struct {
	ClientID 				int64
	Address 				string
	Username 				string
	Database 				string
	ConnectionDate 	time.Time
	LastActivity 		time.Time
}

type SQCloudInfo struct {
	SQLiteVersion 			string
	SQCloudVersion			string
	SQCloudBuildDate		time.Time
	SQCloudGitHash 			string

	ServerTime 					time.Time
	ServerCPUs 					int
	ServerOS 						string
	ServerArchitecture 	string

	ServicePID 					int
	ServiceStart 				time.Time
	ServicePort 				int
	SericeMultiplexAPI 	string
}

type SQCloudPlugin struct {
	Name 				string
	Type 				string
	Version 		string
	Copyright 	string
	Description string
}

// Node Functions

func (this *SQCloud) AddNode( Node string, Address string, Cluster string, Snapshot string, Learner bool ) error {
	sql := "ADD"
	if Learner {
		sql += " LEARNER"
	}
	sql += fmt.Sprintf( " NODE %s ADDRESS %s CLUSTER %s SNAPSHOT %s", SQCloudEnquoteString( Node ), SQCloudEnquoteString( Address ), SQCloudEnquoteString( Cluster ), SQCloudEnquoteString( Snapshot ) )
	return this.Execute( sql )
}
func (this *SQCloud) RemoveNode( Node string ) error {
	return this.Execute( fmt.Sprintf( "REMOVE NODE %s", SQCloudEnquoteString( Node ) ) )
}
func (this *SQCloud) ListNodes() ( []string, error ) {
	return this.SelectStringList( "LIST NODES")
}

// Connection Functions

func (this *SQCloud) CloseConnection( Connection string ) error {
	return this.Execute( fmt.Sprintf( "CLOSE CONNECTION %s", SQCloudEnquoteString( Connection ) ) )
}
func (this *SQCloud) ListConnections() ( []SQCloudConnection, error ) {
	connectionList := []SQCloudConnection{}
	result, err := this.Select( "LIST CONNECTIONS" )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 6 {
				rows :=result.GetNumberOfRows()
				for row := uint( 1 ); row < rows; row++ {
					connectionList = append( connectionList, SQCloudConnection{ 
						ClientID: result.GetInt64Value( row, 1 ), 
						Address:  result.CGetStringValue( row, 2 ),
						Username: result.CGetStringValue( row, 3 ),
						Database: result.CGetStringValue( row, 4 ),
						ConnectionDate: result.GetSQLDateTime( row, 5 ),
						LastActivity: result.GetSQLDateTime( row, 6 ),
					} )
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New( "ERROR: Query returned not 6 Columns (-1)" )
		}
		return []SQCloudConnection{}, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return []SQCloudConnection{}, err
}
func (this *SQCloud) ListDatabaseConnections( Database string ) ( []SQCloudConnection, error ) {
	connectionList := []SQCloudConnection{}
	result, err := this.Select( fmt.Sprintf( "LIST DATABASE CONNECTIONS %s", SQCloudEnquoteString( Database ) ) )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 2 {
				rows :=result.GetNumberOfRows()
				for row := uint( 1 ); row < rows; row++ {
					connectionList = append( connectionList, SQCloudConnection{ 
						ClientID: result.GetInt64Value( row, 1 ), 
						Address:  "",
						Username: "",
						Database: Database,
						ConnectionDate: time.Unix( 0, 0 ),
						LastActivity: result.GetSQLDateTime( row, 2 ),
					} )
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New( "ERROR: Query returned not 2 Columns (-1)" )
		}
		return []SQCloudConnection{}, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return []SQCloudConnection{}, err
}
func (this *SQCloud) ListDatabaseConnectionIds( Id uint ) ( []SQCloudConnection, error ) {
	connectionList := []SQCloudConnection{}
	result, err := this.Select( fmt.Sprintf( "LIST DATABASE CONNECTIONS ID %d", Id ) )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 2 {
				rows :=result.GetNumberOfRows()
				for row := uint( 1 ); row < rows; row++ {
					connectionList = append( connectionList, SQCloudConnection{ 
						ClientID: result.GetInt64Value( row, 1 ), 
						Address:  "",
						Username: "",
						Database: "",
						ConnectionDate: time.Unix( 0, 0 ),
						LastActivity: result.GetSQLDateTime( row, 2 ),
					} )
				}
				result.Free()
				return connectionList, nil
			}
			result.Free()
			return []SQCloudConnection{}, errors.New( "ERROR: Query returned not 2 Columns (-1)" )
		}
		return []SQCloudConnection{}, errors.New( "ERROR: Query returned no result (-1)" )
	}
	return []SQCloudConnection{}, err
}

// Auth Functions

func (this *SQCloud) Auth( Username string, Password string ) error {
	sql := "AUTH"
	if strings.TrimSpace( Username ) != "" {
		sql += fmt.Sprintf( " USER %s", SQCloudEnquoteString( Username ) )
	}
	if strings.TrimSpace( Password ) != "" {
		sql += fmt.Sprintf( " PASS %s", SQCloudEnquoteString( Password ) )
	}
	return this.Execute( sql )
}

// Database funcitons

func (this *SQCloud) CreateDatabase( Database string, Key string, Encoding string, NoError bool ) error {
	sql := fmt.Sprintf( "CREATE DATABASE %s", SQCloudEnquoteString( Database ) )
	if strings.TrimSpace( Key ) != "" {
		sql += fmt.Sprintf( " KEY", SQCloudEnquoteString( Key ) )
	}
	if strings.TrimSpace( Encoding ) != "" {
		sql += fmt.Sprintf( " ENCODING", SQCloudEnquoteString( Encoding ) )
	}
	if NoError {
		sql += " IF NOT EXISTS"
	}
	return this.Execute( sql )
}
func (this *SQCloud) DropDatabase( Database string, NoError bool ) error {
	sql := fmt.Sprintf( "DROP DATABASE %s", SQCloudEnquoteString( Database ) )
	if NoError {
		sql += " IF EXISTS"
	}
	return this.Execute( sql )
}
func (this *SQCloud) ListDatabases() ( []string, error ) {
	return this.SelectStringList( "LIST DATABASES" )
}
func (this *SQCloud) GetDatabase() ( string, error ) {
	return this.SelectSingleString( "GET DATABASE " )
}
func (this *SQCloud) GetDatabaseID() ( int64, error ) {
	return this.SelectSingleInt64( "GET DATABASE ID" )
}
func (this *SQCloud) UseDatabase( Database string ) error {
	this.Database = Database
	return this.Execute( fmt.Sprintf( "USE DATABASE %s", Database ) )
}
func (this *SQCloud) UnuseDatabase() error {
	this.Database = ""
	return this.Execute( "UNUSE DATABASE" )
}

// Plugin Functions

func (this *SQCloud) EnablePlugin( Plugin string ) error {
	return this.Execute( fmt.Sprintf( "ENABLED PLUGIN %s", SQCloudEnquoteString( Plugin ) ) )
}
func (this *SQCloud) DisablePlugin( Plugin string ) error {
	return this.Execute( fmt.Sprintf( "DISABLE PLUGIN %s", SQCloudEnquoteString( Plugin  )) )
}
func (this *SQCloud) ListPlugins() ( []SQCloudPlugin, error ) {
	pluginList := []SQCloudPlugin{}
	result, err := this.Select( "LIST PLUGINS" )
	if err == nil {
		if result != nil {
			if result.GetNumberOfColumns() == 6 {
				rows :=result.GetNumberOfRows()
				for row := uint( 1 ); row < rows; row++ {
					pluginList = append( pluginList, SQCloudPlugin{ 
						Name:        result.CGetStringValue( row, 1 ), 
						Type:        result.CGetStringValue( row, 2 ),
						Version:     result.CGetStringValue( row, 3 ),
						Copyright:   result.CGetStringValue( row, 4 ),
						Description: result.CGetStringValue( row, 5 ),
					} )
				}
				result.Free()
				return pluginList, nil
			}
			result.Free()
			return []SQCloudPlugin{}, errors.New( "ERROR: Query returned not 6 Columns (-1)" )
		}
		return []SQCloudPlugin{}, errors.New( "ERROR: Query returned no results (-1)" )
	}
	return []SQCloudPlugin{}, err
}

// Key / Value Pair functions

func (this *SQCloud) SetKey( Key string, Value string ) error {
	return this.Execute( fmt.Sprintf( "SET KEY %s TO %s", SQCloudEnquoteString( Key ), SQCloudEnquoteString( Value ) ) )
}
func (this *SQCloud) GetKey( Key string ) ( string, error ) {
	return this.SelectSingleString( fmt.Sprintf( "GET KEY %s", Key ) )
}
func (this *SQCloud) DropKey( Key string ) error {
	return this.Execute( fmt.Sprintf( "DROP KEY %s", SQCloudEnquoteString( Key ) ) )
}
func (this *SQCloud) ListKeys() ( []SQCloudKeyValues, error ) {
	return this.SelectKeyValues( "LIST KEYS" )
}
func (this *SQCloud) ListClientKeys() ( []SQCloudKeyValues, error ) {
	return this.SelectKeyValues( "LIST CLIENT KEYS" )
}
func (this *SQCloud) ListDatabaseKeys() ( []SQCloudKeyValues, error ) {
	return this.SelectKeyValues( "LIST DATABASE KEYS" )
}

// Channel Communication

func (this *SQCloud) Listen( Channel string ) error { // add a call back function...
	return this.Execute( fmt.Sprintf( "LISTEN %s", SQCloudEnquoteString( Channel ) ) )
}
func (this *SQCloud) Notify( Channel string ) error {
	return this.Execute( fmt.Sprintf( "NOTIFY %s", SQCloudEnquoteString( Channel ) ) )
}
func (this *SQCloud) SendNotificationMessage( Channel string, Message string ) error {
	return this.Execute( fmt.Sprintf( "NOTIFY %s %s", SQCloudEnquoteString( Channel ), SQCloudEnquoteString( Message ) ) )
}
func (this *SQCloud) Unlisten( Channel string ) error {
	return this.Execute( fmt.Sprintf( "UNLISTEN %s", SQCloudEnquoteString( Channel ) ) )
}

/// Misc functions

func (this *SQCloud) Ping() error {
	result, err := this.Select( "PING" )
	if err == nil {
		if result != nil {
			defer result.Free()
			if result.GetType() == RESULT_STRING {
				if result.CGetResultBuffer() == "PONG" {
					return nil
				}
			}
		}
		return errors.New( "ERROR: Unexpected result (-1)" )
	}
	return err
}
func (this *SQCloud) ListCommands() ( []string, error ) {
	return this.SelectStringList( "LIST COMMANDS" )
}
func (this *SQCloud) ListInfo() ( SQCloudInfo, error ) {
	info := SQCloudInfo{
		SQLiteVersion: 			"0.0.0",
		SQCloudVersion:			"0.0.0",
		SQCloudBuildDate:		time.Unix( 0, 0 ),
		SQCloudGitHash: 		"N/A",	
		ServerTime: 				time.Unix( 0, 0 ),	
		ServerCPUs: 				0,	
		ServerOS: 					"N/A",	
		ServerArchitecture: "N/A",	
		ServicePID: 				0	,
		ServiceStart: 			time.Unix( 0, 0 ),	
		ServicePort: 				0,
		SericeMultiplexAPI: "N/A",	
	}

	result, err := this.SelectKeyValues( "LIST INFO" )
	if err == nil {
		for _, pair := range result {
			switch pair.Key {
			  case "sqlitecloud_version": 		info.SQCloudVersion = pair.Value
			  case "sqlite_version": 					info.SQLiteVersion = pair.Value
			  case "sqlitecloud_build_date": 	info.SQCloudBuildDate, _ = time.Parse( "Jan 2 2006", pair.Value )
			  case "sqlitecloud_git_hash": 		info.SQCloudGitHash = pair.Value
			  case "os": 											info.ServerOS = pair.Value
			  case "arch_bits": 							info.ServerArchitecture = pair.Value
			  case "multiplexing_api": 				info.SericeMultiplexAPI = pair.Value
			  case "tcp_port": 								info.ServicePort, _ = strconv.Atoi( pair.Value )
			  case "process_id": 							info.ServicePID, _ = strconv.Atoi( pair.Value )
			  case "num_processors": 					info.ServerCPUs, _ = strconv.Atoi( pair.Value )
			  case "startup_datetime": 				info.ServiceStart, _ = time.Parse( "2006-01-02 15:04:05", pair.Value )
			  case "current_datetime": 				info.ServerTime, _ = time.Parse( "2006-01-02 15:04:05", pair.Value )
			default:
			}
		}
	}
	return info, err
}
func (this *SQCloud) ListTables() ( []string, error ) {
	return this.SelectStringList( "LIST TABLES" )
}