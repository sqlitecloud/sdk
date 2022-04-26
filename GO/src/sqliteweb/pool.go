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
  "sqlitecloud"
  "strings"
  "sync"
  "time"
  "fmt"
)

type Connection struct {
  node        string
  locked      bool
  uses        uint
  connection  *sqlitecloud.SQCloud
}

type ConnectionManager struct {
  nodeMutex   sync.Mutex
  nodes       map[string][]string

  poolMutex   sync.Mutex
  pool        map[string][]*Connection

  ticker      *time.Ticker
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

func NewConnectionManager() ( *ConnectionManager, error ) {
  pool := &ConnectionManager{
    nodes:  make( map[string][]string ),
    pool:   make( map[string][]*Connection ),
    ticker: nil,
  }
  pool.Start()
  return pool, nil
}

func ( this *ConnectionManager ) Start() error {
  if this.ticker == nil {
    this.ticker = time.NewTicker( time.Second * 10 )
    go func() {
      for {
        <-this.ticker.C
        this.tick()
    } }()
  }
  if this.ticker == nil { return errors.New( "ZZZ" ) }
  return nil
}
func ( this *ConnectionManager ) Stop() error {
  if this.ticker != nil {
    this.ticker.Stop()
    this.ticker = nil
  }
  return nil
}
func ( this *ConnectionManager ) tick() {
  // every 10 seconds...
}

////

func ( this *ConnectionManager ) getServerList( node string ) ( []string, error ) {
  switch {
  case node == "server"   : fallthrough
  case node == "core"     : fallthrough
  case node == "www"      : fallthrough
  case node == "lua"      : fallthrough
  case node == "api"      : fallthrough
  case node == "dashboard": fallthrough
  case node == "admin"    : return []string{}, errors.New( "Access denied to this node" )

  // now remains: auth and uuid

  // auth or uuid in config file?
  case GetINIString( node, "nodes", "" ) != "":
    serverList := []string{}
    for _, server := range strings.Split( GetINIString( node, "nodes", "" ), "," ) {
      server = strings.TrimSpace( server )
      if server != "" { serverList = append( serverList, server ) }
    }
    return serverList, nil

  // search for uuid in auth database
  default:
    if node != "auth" {
      query := fmt.Sprintf( "SELECT 'sqlitecloud://' || username || ':' || password || '@' || IIF( addr6, addr6, addr4 ) || IIF( port, ':' || port, '' ) AS Node FROM PROJECT JOIN NODE ON uuid == project_uuid WHERE uuid = '%s';", sqlitecloud.SQCloudEnquoteString( node ) );

      if result, err := this.ExecuteSQL( "auth", query ); result == nil {
        return []string{}, errors.New( "ERROR: Query returned no result (-1)" )
      } else {
        defer result.Free()

        switch {
        case err != nil                                     : return []string{}, err
        case result.IsError()                               : return []string{}, errors.New( result.GetErrorAsString() )
        case result.IsString()                              : return []string{ result.GetString_() }, nil
        case !result.IsRowSet()                             : return []string{}, errors.New( "ERROR: Query returned an invalid result" )
        case result.GetNumberOfColumns() != 1               : return []string{}, errors.New( "ERROR: Query returned not exactly one column" )
        default:
          stringList := []string{}
          for _, row := range result.Rows() {
            switch val, err := row.GetValue( 0 ); {
            case err != nil                                 : return []string{}, err
            case val == nil                                 : continue
            case !val.IsString()                            : continue
            case strings.TrimSpace( val.GetString() ) == "" : continue
            default                                         : stringList = append( stringList, strings.TrimSpace( val.GetString() ) )
          } }
          return stringList, nil
  } } } }
  return []string{}, fmt.Errorf( "No Node found for project: '%s'", node )
}
func ( this *ConnectionManager ) getNextServer( node string ) ( server string, err error ) {
  server = ""
  err    = nil

  this.nodeMutex.Lock()
  _, found := this.nodes[ node ]
  this.nodeMutex.Unlock()

  if !found {
    switch nodes, err := this.getServerList( node ); {
    case err != nil       : return "", err
    case len( nodes ) < 1 : return "", nil
    default               :
      this.nodeMutex.Lock()
      if _, found := this.nodes[ node ]; !found { this.nodes[ node ] = nodes }
      this.nodeMutex.Unlock()
    }
  }

  this.nodeMutex.Lock()
  nodes, found := this.nodes[ node ]
  switch {
  case !found             : break
  case len( nodes ) < 1   : break
  case len( nodes ) == 1  : server = nodes[ 0 ]
  default                 :
    server = nodes[ 0 ]
    this.nodes[ node ] = append( nodes[ 1 : ], server )
  }
  this.nodeMutex.Unlock()
  return server, nil
}

////
func ( this *ConnectionManager ) createAndAppendNewLockedConnection( node string ) ( connection *Connection, err error ) {
  connection   = &Connection{
    node       : node,
    locked     : true,
    uses       : 0,
    connection : nil,
  }

  if connectionString, err := this.getNextServer( node ); err == nil { // TODO: repeat...
    if connection.connection, err = sqlitecloud.Connect( connectionString ); err != nil {
      if connection.connection != nil {
        connection.connection.Close()
        connection.connection = nil
      }
      return nil, err
    } else {

      this.poolMutex.Lock()
      this.pool[ node ] = append( this.pool[ node ], connection )
      this.poolMutex.Unlock()

      return connection, nil
    }
  } else {
    return nil, err
  }
}
func ( this *ConnectionManager ) lockConnection( connection *Connection ) {
  if connection != nil {
    this.poolMutex.Lock()
    connection.locked = true
    this.poolMutex.Unlock()
  }
}
func ( this *ConnectionManager ) unLockConnection( connection *Connection ) {
  if connection != nil {
    this.poolMutex.Lock()
    connection.locked = false
    this.poolMutex.Unlock()
  }
}

func ( this *ConnectionManager ) closeAndRemoveLockedConnection( node string, connection *Connection ) error {
  if connection != nil {
    if connection.locked {

      if connection.connection != nil {
        connection.connection.Close()
        connection.connection = nil
      }

      for i, c := range this.pool[ node ] {
        switch {
        case c != connection: continue
        default:
          this.poolMutex.Lock()
          this.pool[ node ] = append( this.pool[ node ][ : i ], this.pool[ node ][ i + 1 : ]... )
          this.poolMutex.Unlock()
          return nil
      } }

    } else {
      errors.New( "Connection not locked" )
  } }
  return errors.New( "Connection not found" )
}

func ( this *ConnectionManager ) getFirstUnlockedConnectionAndLockIt( node string ) ( connection *Connection ) {
  connection = nil
  this.poolMutex.Lock()

  switch connections, found := this.pool[ node ]; {
  case !found                     : break
  case len( connections ) == 0    : break
  case connections[ 0 ].locked    : break
  default                         :
    connection        = connections[ 0 ]
    connection.locked = true
  }

	if connection != nil {
		fmt.Printf( "Using connection: '%s'\r\n", connection.node )
	}
	
  this.poolMutex.Unlock()
  return connection
}
func ( this *ConnectionManager ) moveConnectionToEnd( node string, connection *Connection ) ( err error ) {
  err = nil
  this.poolMutex.Lock()

  if connection == nil {
    MOVE:
    switch connections, found := this.pool[ node ]; {
    case !found                   : err = errors.New( "Connection not found" )
    case len( connections ) == 0  : break MOVE
    default                       :
      for i, c := range connections {
        switch {
        case c != connection: continue
        default:
          this.pool[ node ] = append( connections[ : i ], append( connections[ i + 1 : ], connections[ i ] )... )
          break MOVE
  } } } }

  this.poolMutex.Unlock()
  return nil
}

// Retry!!!!
func ( this *ConnectionManager ) GetConnection( node string ) ( connection *Connection, err error ) {
  if connection = this.getFirstUnlockedConnectionAndLockIt( node ); connection != nil {
    err = this.moveConnectionToEnd( node, connection )
  } else {
    connection, err = this.createAndAppendNewLockedConnection( node )
  }
  return connection, err
}
func ( this *ConnectionManager ) ReleaseConnection( node string, connection *Connection ) ( err error ) {
  err = nil

  if connection != nil {
    this.lockConnection( connection )
		maxRequests := uint( GetINIInt( "core", "maxrequests", 0 ) )
    if maxRequests != 0 && connection.uses > maxRequests {
      err = this.closeAndRemoveLockedConnection( node, connection )
    } else {
      this.unLockConnection( connection )
    }
  }
  return err
}


func ( this *ConnectionManager ) ExecuteSQL( node string, query string ) ( *sqlitecloud.Result, error ) {
  var connection *Connection  = nil
  var res *sqlitecloud.Result = nil
  var err error               = nil

                                                       maxTries := 3
  if _, nodeExists := this.nodes[ node ]; nodeExists { maxTries  = 3 + len( this.nodes[ node ] ) }

  for try := 0; try < maxTries; try++ {
    connection, err = this.GetConnection( node )
    switch {
    case err != nil                   : fallthrough
    case connection == nil            : fallthrough
    case connection.connection == nil : continue
    default                           :

      connection.uses++
      if res, err = connection.connection.Select( query ); res == nil && err == nil {
        continue
      } else if connection.connection.ErrorCode >= 100000 {
        // internal error (the SDK cannot write to or read from the connection) 
        // so remove the current failed connection and retry with a new one
        // for example: 
        // - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
        // - 100003 Internal Error: sendString (%s)
        this.closeAndRemoveLockedConnection( node, connection )
        continue
      } else {
        this.ReleaseConnection( node, connection )
        return res, err
      }
    }
    this.closeAndRemoveLockedConnection( node, connection )
  }
  return nil, nil
}