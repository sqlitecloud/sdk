// Package sqlitecloud provides an easy to use GO driver for connecting to and using the SQLite Cloud Database server.
package sqlitecloud

// #include <stdlib.h>
// #include "sqcloud.h"
import "C"

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "errors"
//mport "time"
import "strconv"
import "net"

import "github.com/xo/dburl"

type SQCloud struct {
  connection    *C.struct_SQCloudConnection

  Host          string
  Port          int
  Username      string
  Password      string
  Database      string
  Timeout       int
  Family        int

  UUID          string

  ErrorCode     int
  ErrorMessage  string
}

// init registers the sqlitecloud scheme in the connection steing parser.
func init() {
  dburl.Register( dburl.Scheme {
    Driver: "sc", // sqlitecloud
    Generator: dburl.GenFromURL("sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=disabled&sslmode=disabled"),
    Transport: dburl.TransportTCP,
    Opaque: false,
    Aliases: []string{ "sqlitecloud" },
    Override: "",
  } )
}

// Helper functions

// ParseConnectionString parses the given connection string and returns it's components.
// An empty string ("") or -1 is returned as the corresponding return value, if a component of the connection string was not present.
// No plausibility checks are done (see: CheckConnectionParameter).
// If the connection string could not be parsed, an error is returned.
func ParseConnectionString( ConnectionString string ) ( Host string, Port int, Username string, Password string, Database string, Timeout int, Compress string, err error ) {
  u, err := dburl.Parse( ConnectionString ) // sqlitecloud://dev1.sqlitecloud.io/X?timeout=14&compress=LZ4
  if err == nil {

    host      := u.Hostname()
    user      := u.User.Username()
    pass, _   := u.User.Password()
    database  := strings.TrimPrefix( u.Path, "/" )
    timeout   := 0
    compress  := "NO"
    port      := 0
    if port, err = strconv.Atoi( u.Port() ); err != nil {
      port = -1
    }

    for key, values := range u.Query() {
      lastValue := values[ len( values ) - 1 ]
      switch strings.ToLower( strings.TrimSpace( key ) ) {
      case "timeout":
        if timeout, err = strconv.Atoi( lastValue ); err != nil {
          timeout = -1
        } 

      case "compress":
        compress = strings.ToUpper( lastValue )
      
      default: // Ignore
      }
    }

    fmt.Printf( "NO ERROR: Host=%s, Port=%d, User=%s, Password=%s, Database=%s, Timeout=%d, Compress=%s\r\n", host, port, user, pass, database, timeout, compress )
    return host, port, user, pass, database, timeout, compress, nil
  }
  return "", -1, "", "", "", -1, "", err
}

// CheckConnectionParameter checks the given connection arguments for validly.
// Host is either a resolve able hostname or an IP address.
// Port is an unsigned int number between 1 and 65535. 
// Timeout must be 0 (=no timeout) or a positive number. 
// Compress must be "NO" or "LZ4".
// Username, Password and Database are ignored.
// If a given value does not fulfill the above criteria's, an error is returned.
func (this *SQCloud) CheckConnectionParameter( Host string, Port int, Username string, Password string, Database string, Timeout int, Compress string ) error {
  
  if strings.TrimSpace( Host ) == ""  {
    return errors.New( "Invalid Hostname" )
  }

  ip := net.ParseIP( Host )
  if ip == nil {
    if _, err := net.LookupHost( Host ); err != nil {
      return errors.New( "Can't resolve Hostname" )
    }
  }

  if Port < 1 || Port >= 0xFFFF {
    return errors.New( "Invalid Port" )
  }

  if Timeout < 0 {
    return errors.New( "Invalid Timeout" )
  }

  switch Compress {
  case "NO", "LZ4":
  default:
    return errors.New( "Invalid Compression Method" )
  }

  return nil
}

// Creation

// reset resets all Connection attributes.
func (this *SQCloud) reset() {
  this.Close()
  this.resetError()
  this.Host       = ""
  this.Port       = -1
  this.Username   = ""
  this.Password   = ""
  this.Database   = ""
  this.Timeout    = -1
  this.Family     = -1
  this.UUID       = ""
}

// New creates an empty connection and resets it.
// A pointer to the newly created connection is returned (see: Connect).
func New() *SQCloud {
  connection := SQCloud{ connection: nil }
  connection.reset()
  return &connection
}

// Connect creates a new connection and tries to connect to the server using the given connection string.
// The given connection string is parsed and checked for correct parameters.
// Nil and an error is returned if the connection string had invalid values or a connection to the server could not be established,
// otherwise, a pointer to the newly established connection is returned.
func Connect( ConnectionString string ) (*SQCloud, error) {
  Host, Port, Username, Password, Database, Timeout, Compress, err := ParseConnectionString( ConnectionString )

  if err != nil {
    return nil, err
  }

  connection := New()
  err = connection.Connect( Host, Port, Username, Password, Database, Timeout, Compress, 0 )
  if err != nil {
    connection.Close()
    return nil, err
  }

  switch( Compress ) {
  case "LZ4": connection.Compress( true )
  default:    connection.Compress( false )
  }
  
  return connection, nil
}

// Connection Functions

// Connect connects to a SQLite Cloud server instance using the given arguments.
// If Connect is called on an already established connection, the old connection is closed first.
// All arguments are checked for valid values (see: CheckConnectionParameter). An error is returned if the protocol Family was not '0', 
// invalid argument values where given or the connection could not be established. 
func (this *SQCloud) Connect( Host string, Port int, Username string, Password string, Database string, Timeout int, Compression string, Family int ) error {
  this.reset()

  if Family != 0 {
    return errors.New( "Invalid Protocol Family" )
  }

  if err := this.CheckConnectionParameter( Host, Port, Username, Password, Database, Timeout, Compression ); err != nil {
    return err
  }

  this.Host     = Host
  this.Port     = Port
  this.Username = Username
  this.Password = Password
  this.Database = Database
  this.Timeout  = Timeout
  this.Family   = Family

  return this.reconnect()
}

// reconnect closes and then reopens a connection to the SQLite Cloud database server.
func (this *SQCloud) reconnect() error {
  this.Close()
  this.resetError()

  this.connection = CConnect( this.Host, this.Port, this.Username, this.Password, this.Database, this.Timeout, this.Family )

  if this.connection != nil {
    this.ErrorCode    = this.CGetErrorCode()
    this.ErrorMessage = this.CGetErrorMessage()
 
    if !this.IsError() {
      if strings.TrimSpace( this.Database ) != "" {
        this.UseDatabase( this.Database )
      }
    }

  } else {

    this.ErrorCode    = 666
    this.ErrorMessage = "Not enoght memory to allocate a SQCloudConnection."
  }

  if this.IsError() {
    err := errors.New( fmt.Sprintf( "ERROR CONNECTION TO %s: %s (%d)", this.Host, this.ErrorMessage, this.ErrorCode ) )
    this.Close()
    return err
  }

  return nil
}

// Close closes the connection to the SQLite Cloud Database server.
// The connection can later be reopened (see: reconnect)
func (this *SQCloud) Close() {
  if this.IsConnected() {
    this.CDisconnect()
  }
  this.connection = nil
  this.resetError()
}

// Compress enabled or disables data compression for this connection.
// If enabled, the data is compressed with the LZ4 compression algorithm, otherwise no compression is applied the data.
func (this *SQCloud) Compress( Enabled bool ) error {
  switch Enabled {
    case false: return this.Execute( "SET KEY CLIENT_COMPRESSION TO 0" )
    default:    return this.Execute( "SET KEY CLIENT_COMPRESSION TO 1" )
  }
}


// IsConnected checks the connection to the SQLite Cloud database server by sending a PING command.
// true is returned, if the connection is established and actually working, false otherwise.
func (this *SQCloud) IsConnected() bool {
  if( this.connection == nil ) {
    return false
  }
  if( this.Ping() != nil ) {
    return false
  }
  return true
}

// Error Methods

// resetError resets the error code and message of the last run command.
func (this *SQCloud) resetError() {
  this.ErrorCode    = 0
  this.ErrorMessage = ""
}

// IsError checks the successful execution of the last method call / command.
// true is returned if the last command resulted in an error, false otherwise.
func (this *SQCloud) IsError() bool {
  return this.ErrorCode != 0
}

// GetErrorCode returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud ) GetErrorCode() int {
  return this.ErrorCode
}

// GetErrorMessage returns the error message of the last unsuccuessful command as an error.
// nil is returned if the last command run successful.
func (this *SQCloud ) GetErrorMessage() error {
  if this.IsError() {
    return errors.New( this.ErrorMessage )
  }
  return nil
}

// GetError returned the error code and message of the last unsuccessful command.
// 0 and nil is returned if the last command run successful.
func (this *SQCloud ) GetError() ( int, error ) {
  return this.GetErrorCode(), this.GetErrorMessage()
}


// Data Access Functions

// Select executes a query on an open SQLite Cloud database connection.
// If an error occurs during the execution of the query, nil and an error describing the problem is returned.
// On successful execution, a pointer to the result is returned. 
func (this *SQCloud) Select( SQL string ) (*SQCloudResult, error) {
  this.resetError()

  result           := this.CExec( SQL )
  this.ErrorCode    = this.CGetErrorCode()
  this.ErrorMessage = this.CGetErrorMessage()

  if result != nil {
    result.Rows           = result.CGetRows()
    result.Columns        = result.CGetColumns()
    result.ColumnWidth    = make( []uint, result.Columns )
    result.HeaderWidth    = make( []uint, result.Columns )
    result.MaxHeaderWidth = 0  

    result.Type           = result.CGetResultType()
    result.ErrorCode      = this.ErrorCode
    result.ErrorMessage   = this.ErrorMessage
    
    for c := uint( 0 ); c < result.Columns ; c++ {
      result.HeaderWidth[ c ] = uint( len( result.GetColumnName( c ) ) )
      result.ColumnWidth[ c ] = result.CGetMaxColumnLenght( c )
      if result.ColumnWidth[ c ] < result.HeaderWidth[ c ] {
        result.ColumnWidth[ c ] = result.HeaderWidth[ c ]
      }
      if result.MaxHeaderWidth < result.HeaderWidth[ c ] {
        result.MaxHeaderWidth = result.HeaderWidth[ c ]
      }
    }
  } else {
    return nil, errors.New( "ERROR: Could not execute SQL command (-1)" )
  }

  if this.IsError() {
    result.Free()
    return nil, errors.New( fmt.Sprintf( "ERROR: %s (%d)", this.CGetErrorMessage(), this.CGetErrorCode() ) )
  }
  return result, nil // nil, nil or *result, nil
}

// Execute executes the given query.
// If the execution was not successful, an error describing the reason of the failure is returned.
func (this *SQCloud) Execute( SQL string ) error {
  result, err := this.Select( SQL )
  if result != nil {
    

    isOK := result.IsOK()
    result.Free()
    if !isOK {
      return errors.New( "ERROR: Unexpected Result Set (-1)")
    }
  }
  return err
}