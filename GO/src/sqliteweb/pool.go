//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.1.0
//     //             ///   ///  ///    Date        : 2022/02/18
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

////lint:file-ignore ST1006 receiver name should be a reflection of its identity

package main

import (
	"errors"
	"fmt"
	"sqlitecloud"
	"strings"
	"sync"
	"time"
)

type Connection struct {
	project    string
	locked     bool
	uses       uint
	connection *sqlitecloud.SQCloud
}

func (c Connection) String() string {
	return fmt.Sprintf("{%p %v '%s' %s:%d %d}", c.connection, c.locked, c.project, c.connection.Host, c.connection.Port, c.uses)
}

func printPool(m map[string][]*Connection) {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	for k, v := range m {
		fmt.Printf("    "+k+": "+strings.Repeat(" ", maxLenKey-len(k))+"%v\n", v)
	}
}

type ConnectionManager struct {
	nodeMutex sync.Mutex
	nodes     map[string][]string

	poolMutex sync.Mutex
	pool      map[string][]*Connection

	ticker *time.Ticker
}

var cm *ConnectionManager = nil

func init() {
	cm, _ = NewConnectionManager()

	// c, err := NewConnectionManager()
	// fmt.Printf( "%v\r\n", err )
	//
	// r, e := c.ExecuteSQL( "fbf94289-64b0-4fc6-9c20-84083f82ee64", "LIST DATABASES" )
	// fmt.Printf( "%v\r\n", e )
	// r.DumpToScreen( 0 )
}

func NewConnectionManager() (*ConnectionManager, error) {
	pool := &ConnectionManager{
		nodes:  make(map[string][]string),
		pool:   make(map[string][]*Connection),
		ticker: nil,
	}
	pool.Start()
	return pool, nil
}

func (this *ConnectionManager) Start() error {
	if this.ticker == nil {
		this.ticker = time.NewTicker(time.Second * 10)
		go func() {
			for {
				<-this.ticker.C
				this.tick()
			}
		}()
	}
	if this.ticker == nil {
		return errors.New("ZZZ")
	}
	return nil
}
func (this *ConnectionManager) Stop() error {
	if this.ticker != nil {
		this.ticker.Stop()
		this.ticker = nil
	}
	return nil
}
func (this *ConnectionManager) tick() {
	// every 10 seconds...
}

////
func (this *ConnectionManager) getServerList(project string) ([]string, error) {
	switch {
	case project == "server":
		fallthrough
	case project == "core":
		fallthrough
	case project == "www":
		fallthrough
	case project == "lua":
		fallthrough
	case project == "api":
		fallthrough
	case project == "dashboard":
		fallthrough
	case project == "admin":
		return []string{}, errors.New("Access denied to this node")

	// now remains: auth and uuid

	// auth or uuid in config file?
	case GetINIString(project, "nodes", "") != "":
		serverList := []string{}
		for _, server := range strings.Split(GetINIString(project, "nodes", ""), ",") {
			server = strings.TrimSpace(server)
			if server != "" {
				serverList = append(serverList, server)
			}
		}
		return serverList, nil

	// search for uuid in auth database
	default:
		if project != "auth" {
			query := "SELECT 'sqlitecloud://' || admin_username || ':' || admin_password || '@' || IIF( addr6, addr6, addr4 ) || IIF( port, ':' || port, '' ) AS Node FROM Project JOIN Node ON uuid == project_uuid WHERE uuid = ?;"
			args := []interface{}{project}
			SQLiteWeb.Logger.Debugf("query %s %v\n", query, args)

			if result, err, _, _ := this.ExecuteSQLArray("auth", query, &args); result == nil || result.GetNumberOfRows() == 0 {
				return []string{}, errors.New("ERROR: Query returned no result (-1)")
			} else {
				defer result.Free()

				switch {
				case err != nil:
					return []string{}, err
				case result.IsError():
					return []string{}, errors.New(result.GetErrorAsString())
				case result.IsString():
					return []string{result.GetString_()}, nil
				case !result.IsRowSet():
					return []string{}, errors.New("ERROR: Query returned an invalid result")
				case result.GetNumberOfColumns() != 1:
					return []string{}, errors.New("ERROR: Query returned not exactly one column")
				default:
					stringList := []string{}

					for _, row := range result.Rows() {
						switch val, err := row.GetValue(0); {
						case err != nil:
							return []string{}, err
						case val == nil:
							continue
						case !val.IsString():
							continue
						case strings.TrimSpace(val.GetString()) == "":
							continue
						default:
							stringList = append(stringList, strings.TrimSpace(val.GetString()))
						}
					}
					return stringList, nil
				}
			}
		}
	}
	return []string{}, fmt.Errorf("No Node found for project: '%s'", project)
}

func (this *ConnectionManager) getNextServer(project string, reloadNodes bool) (server string, err error) {
	server = ""
	err = nil

	this.nodeMutex.Lock()
	_, found := this.nodes[project]
	this.nodeMutex.Unlock()

	if !found || reloadNodes {
		switch nodes, err := this.getServerList(project); {
		case err != nil:
			return "", err
		case len(nodes) < 1:
			return "", nil
		default:
			this.nodeMutex.Lock()
			if _, found := this.nodes[project]; !found || reloadNodes {
				this.nodes[project] = nodes
			}
			this.nodeMutex.Unlock()
		}
	}

	this.nodeMutex.Lock()
	nodes, found := this.nodes[project]
	switch {
	case !found:
		break
	case len(nodes) < 1:
		break
	case len(nodes) == 1:
		server = nodes[0]
	default:
		server = nodes[0]
		this.nodes[project] = append(nodes[1:], server)
	}
	this.nodeMutex.Unlock()

	return server, nil
}

func isProjectUuid(project string) bool {
	switch project {
	case "server", "core", "www", "lua", "api", "dashboard", "admin", "auth":
		return false

	// project uuid
	default:
		return true
	}
}

func (this *ConnectionManager) connectionWithIniParams(project string, connectionString string) (*sqlitecloud.SQCloud, error) {
	config, err := sqlitecloud.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	if isProjectUuid(project) {
		config.Timeout = GetINIDuration("pool", "timeout", time.Duration(10)*time.Second)
		config.NoBlob = GetINIBoolean("pool", "noblob", true)
		config.MaxData = GetINIInt("pool", "maxdata", 2048)
		config.MaxRowset = GetINIInt("pool", "maxrowset", 1048576)
	}

	return sqlitecloud.New(*config), nil
}

////
func (this *ConnectionManager) createAndAppendNewLockedConnection(project string, reloadNodes bool) (connection *Connection, err error) {
	connection = &Connection{
		project:    project,
		locked:     true,
		uses:       0,
		connection: nil,
	}

	if connectionString, err := this.getNextServer(project, reloadNodes); err == nil { // TODO: repeat...
		if connection.connection, err = this.connectionWithIniParams(project, connectionString); err != nil {
			return nil, err
		}

		if err = connection.connection.Connect(); err != nil {
			if connection.connection != nil {
				connection.connection.Close()
				connection.connection = nil
			}
			return nil, err
		} else {
			this.poolMutex.Lock()
			this.pool[project] = append(this.pool[project], connection)
			this.poolMutex.Unlock()

			return connection, nil
		}
	} else {
		return nil, err
	}
}
func (this *ConnectionManager) lockConnection(connection *Connection) {
	if connection != nil {
		this.poolMutex.Lock()
		connection.locked = true
		this.poolMutex.Unlock()
	}
}
func (this *ConnectionManager) unLockConnection(connection *Connection) {
	if connection != nil {
		this.poolMutex.Lock()
		connection.locked = false
		this.poolMutex.Unlock()
	}
}

func (this *ConnectionManager) closeAndRemoveLockedConnection(project string, connection *Connection) error {
	if connection != nil {
		if connection.locked {

			if connection.connection != nil {
				connection.connection.Close()
				connection.connection = nil
			}

			for i, c := range this.pool[project] {
				switch {
				case c != connection:
					continue
				default:
					this.poolMutex.Lock()
					this.pool[project] = append(this.pool[project][:i], this.pool[project][i+1:]...)
					this.poolMutex.Unlock()
					return nil
				}
			}

		} else {
			errors.New("Connection not locked")
		}
	}
	return errors.New("Connection not found")
}

func (this *ConnectionManager) getFirstUnlockedConnectionAndLockIt(project string) (connection *Connection) {
	connection = nil
	this.poolMutex.Lock()

	switch connections, found := this.pool[project]; {
	case !found:
		break
	case len(connections) == 0:
		break
	case connections[0].locked:
		break
	default:
		connection = connections[0]
		connection.locked = true
	}

	this.poolMutex.Unlock()
	return connection
}
func (this *ConnectionManager) moveConnectionToEnd(project string, connection *Connection) (err error) {
	err = nil
	this.poolMutex.Lock()

	if connection != nil {
	MOVE:
		switch connections, found := this.pool[project]; {
		case !found:
			err = errors.New("Connection not found")
		case len(connections) == 0:
			break MOVE
		default:
			for i, c := range connections {
				switch {
				case c != connection:
					continue
				default:
					this.pool[project] = append(connections[:i], append(connections[i+1:], connections[i])...)
					break MOVE
				}
			}
		}
	}

	this.poolMutex.Unlock()
	return nil
}

func (this *ConnectionManager) GetPoolLen(project string) (length int) {
	length = 0
	switch connections, found := this.pool[project]; {
	case !found:
		break
	default:
		length = len(connections)
	}

	return length
}

// Retry!!!!
func (this *ConnectionManager) GetConnection(project string, reloadNodes bool) (connection *Connection, err error) {
	if connection = this.getFirstUnlockedConnectionAndLockIt(project); connection != nil {
		err = this.moveConnectionToEnd(project, connection)
		SQLiteWeb.Logger.Debugf("(%s) Reusing connection %v, pool len %d", project, connection, this.GetPoolLen(project))
	} else {
		connection, err = this.createAndAppendNewLockedConnection(project, reloadNodes)
		if err != nil {
			SQLiteWeb.Logger.Errorf("(%s) Creating new connection %v, err:%s", project, connection, err)
		} else {
			SQLiteWeb.Logger.Debugf("(%s) Creating new connection %v, pool_len:%d", project, connection, this.GetPoolLen(project))
		}
	}

	// SQLiteWeb.Logger.Debugf("[%s] Using connection %v\n", project, connection)
	// SQLiteWeb.Logger.Debugf("pool:")
	// printPool(this.pool)

	return connection, err
}

func (this *ConnectionManager) ReleaseConnection(project string, connection *Connection) (err error) {
	err = nil

	if connection != nil {
		this.lockConnection(connection)
		maxRequests := uint(GetINIInt("pool", "maxrequests", 0))
		if maxRequests != 0 && connection.uses > maxRequests {
			err = this.closeAndRemoveLockedConnection(project, connection)
		} else {
			this.unLockConnection(connection)
		}
	}
	return err
}

func (this *ConnectionManager) ExecuteSQL(project string, query string) (*sqlitecloud.Result, error, int, int) {
	return this.ExecuteSQLArray(project, query, nil)
}

func (this *ConnectionManager) ExecuteSQLArray(project string, query string, args *[]interface{}) (*sqlitecloud.Result, error, int, int) {
	var connection *Connection = nil
	var res *sqlitecloud.Result = nil
	var err error = nil
	var reloadNodes bool = false

	maxTries := 3
	if _, nodeExists := this.nodes[project]; nodeExists {
		maxTries = 3 + len(this.nodes[project])
	}

	for try := 0; try < maxTries; try++ {
		connection, err = this.GetConnection(project, reloadNodes)
		switch {
		case err != nil:
			fallthrough
		case connection == nil:
			fallthrough
		case connection.connection == nil:
			reloadNodes = true
			continue
		default:
			reloadNodes = false
			connection.uses++

			start := time.Now()
			errCode := int(0)
			extErrCode := int(0)
			if args != nil {
				res, err = connection.connection.SelectArray(query, *args)
			} else {
				res, err = connection.connection.Select(query)
			}

			if res == nil && err == nil {
				continue
			} else if connection.connection.ErrorCode >= 100000 {
				// internal error (the SDK cannot write to or read from the connection)
				// so remove the current failed connection and retry with a new one
				// for example:
				// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
				// - 100003 Internal Error: sendString (%s)
				errCode = connection.connection.ErrorCode
				extErrCode = connection.connection.ExtErrorCode
				this.closeAndRemoveLockedConnection(project, connection)
				SQLiteWeb.Logger.Debugf("(%s) ExecuteSQL connection error %d", project, errCode)
				reloadNodes = true
				continue
			} else {
				errCode = connection.connection.ErrorCode
				extErrCode = connection.connection.ExtErrorCode
				this.ReleaseConnection(project, connection)
				t := time.Since(start)
				SQLiteWeb.Logger.Debugf("(%s) ExecuteSQL query:\"%s\" time:%s", project, query, t)
				return res, err, errCode, extErrCode
			}
		}
		this.closeAndRemoveLockedConnection(project, connection)
	}
	return nil, nil, 0, 0
}
