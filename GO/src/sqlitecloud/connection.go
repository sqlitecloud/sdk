//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.2
//     //             ///   ///  ///    Date        : 2021/10/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Go Methods related to the
//   ////                ///  ///                     SQCloud class for managing
//     ////     //////////   ///                      the connection and executing
//        ////            ////                        queries.
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
  "crypto/tls"
  "crypto/x509"
  "errors"
  "fmt"
  "io/ioutil"
  "net"
  "strconv"
  "strings"
  "time"

  "github.com/xo/dburl"
)

type SQCloud struct {
  sock          *net.Conn

  psub          *SQCloud
  Callback      func( string )
  
  Host          string
  Port          int
  Username      string
  Password      string
  Database      string
  cert          *tls.Config
  Timeout       time.Duration
  Family        int

  uuid          string // 36 runes -> remove maybe????
  secret        string // 36 runes -> remove maybe????

  ErrorCode     int
  ExtErrorCode  int
  ErrorMessage  string
}

// init registers the sqlitecloud scheme in the connection steing parser.
func init() {
  dburl.Register( dburl.Scheme {
    Driver: "sc", // sqlitecloud
    Generator: dburl.GenFromURL( "sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=NO&tls=INTERN" ),
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
func ParseConnectionString( ConnectionString string ) ( Host string, Port int, Username string, Password string, Database string, Timeout int, Compress string, Pem string, err error ) {
  u, err := dburl.Parse( ConnectionString ) // sqlitecloud://dev1.sqlitecloud.io/X?timeout=14&compress=LZ4&tls=INTERN
  if err == nil {

    host      := u.Hostname()
    user      := u.User.Username()
    pass, _   := u.User.Password()
    database  := strings.TrimPrefix( u.Path, "/" )
    timeout   := 0
    compress  := "NO"
    pem       := "INTERN"
    port      := 0
    sPort     := strings.TrimSpace( u.Port() )
    if len( sPort ) > 0 {
      if port, err = strconv.Atoi( sPort ); err != nil {
        return "", -1, "", "", "", -1, "", "", err
      }
    }

    for key, values := range u.Query() {
      lastLiteral := strings.TrimSpace( values[ len( values ) - 1 ] )
      switch strings.ToLower( strings.TrimSpace( key ) ) {
      case "timeout":
        if timeout, err = strconv.Atoi( lastLiteral ); err != nil {
          return "", -1, "", "", "", -1, "", "", err
        }

      case "compress":  compress = strings.ToUpper( lastLiteral )
      case "tls":       pem      = parsePEMString( lastLiteral )
      }
    }

    // fmt.Printf( "NO ERROR: Host=%s, Port=%d, User=%s, Password=%s, Database=%s, Timeout=%d, Compress=%s\r\n", host, port, user, pass, database, timeout, compress )
    return host, port, user, pass, database, timeout, compress, pem, nil
  }
  return "", -1, "", "", "", -1, "", "", err
}

func parsePEMString( Pem string ) string {
  switch strings.ToUpper( strings.TrimSpace( Pem ) ) {
  case "", "0", "N", "NO",  "FALSE", "OFF", "DISABLE", "DISABLED":                                return ""
  case     "1", "Y", "YES", "TRUE",  "ON",  "ENABLE",  "ENABLED", "INTERN", "<USE INTERNAL PEM>": return "<USE INTERNAL PEM>"
  default:                                                                                        return strings.TrimSpace( Pem )
  }
}

// CheckConnectionParameter checks the given connection arguments for validly.
// Host is either a resolve able hostname or an IP address.
// Port is an unsigned int number between 1 and 65535.
// Timeout must be 0 (=no timeout) or a positive number.
// Compress must be "NO" or "LZ4".
// Username, Password and Database are ignored.
// If a given value does not fulfill the above criteria's, an error is returned.
func (this *SQCloud) CheckConnectionParameter( Host string, Port int, Username string, Password string, Database string, Timeout int, Compress string, Pem string ) error {

  if strings.TrimSpace( Host ) == ""  {
    return errors.New( fmt.Sprintf( "Invalid hostname (%s)", Host ) )
  }

  ip := net.ParseIP( Host )
  if ip == nil {
    if _, err := net.LookupHost( Host ); err != nil {
      return errors.New( fmt.Sprintf( "Can't resolve hostname (%s)", Host ) )
    }
  }

  if Port < 1 || Port >= 0xFFFF {
    return errors.New( fmt.Sprintf( "Invalid Port (%d)", Port ) )
  }

  if Timeout < 0 {
    return errors.New( fmt.Sprintf( "Invalid Timeout (%d)", Timeout ) )
  }

  switch strings.ToUpper( Compress ) {
  case "NO", "LZ4":
  default: return errors.New( fmt.Sprintf( "Invalid compression method (%s)", Compress ) )
  }

  switch trimmed := parsePEMString( Pem ); trimmed {
  case "", "<USE INTERNAL PEM>":
  default:
    if _, err := ioutil.ReadFile( trimmed ); err != nil {
      return errors.New( fmt.Sprintf( "Could not open PEM file in '%s'", trimmed ) )
    }
  }

  return nil
}

// Creation

// reset resets all Connection attributes.
func (this *SQCloud) reset() {
  this.Close()
  this.Host       = ""
  this.Port       = -1
  this.Username   = ""
  this.Password   = ""
  this.Database   = ""
  this.Family     = -1
  this.uuid       = ""
  this.secret     = ""
  this.resetError()
}

// New creates an empty connection and resets it.
// A pointer to the newly created connection is returned (see: Connect).
func New( Certificate string, TimeOut uint ) *SQCloud {
  connection := SQCloud{
    sock        : nil,

    psub        : nil,
    Callback    : func( json string ) {}, // empty call back function

    Host        : "",
    Port        : -1,
    Username    : "",
    Password    : "",
    Database    : "",
    cert        : nil,
    Timeout     : time.Duration( TimeOut ) * time.Second,
    Family      : 1,
    uuid        : "",
    secret      : "",
    ErrorCode   : 0,
    ErrorMessage: "",
  }

  switch trimmed := parsePEMString( Certificate ); trimmed {
  case "":                    break             // unencrypted connection
  case "<USE INTERNAL PEM>":  Certificate = PEM // use internal Certificate
  default:
    switch pem, err := ioutil.ReadFile( trimmed ); {
    case err != nil:   return nil
    default:           Certificate = string( pem )
    }
  }

  if len( Certificate ) != 0 {
    pool := x509.NewCertPool()
    if !pool.AppendCertsFromPEM( []byte( Certificate ) ) { return nil }

    connection.cert = &tls.Config {
      RootCAs:            pool,
      InsecureSkipVerify: true,
    }
  }

  return &connection
}

// Connect creates a new connection and tries to connect to the server using the given connection string.
// The given connection string is parsed and checked for correct parameters.
// Nil and an error is returned if the connection string had invalid values or a connection to the server could not be established,
// otherwise, a pointer to the newly established connection is returned.
func Connect( ConnectionString string ) ( *SQCloud, error ) {
  Host, Port, Username, Password, Database, Timeout, Compress, Pem, err := ParseConnectionString( ConnectionString )

  if err != nil   { return nil, err }
  if Port == 0    { Port    = 8860  }
  if Timeout == 0 { Timeout = 10    }

  connection := New( Pem, uint( Timeout ) ) // allways works

  if err = connection.Connect( Host, Port, Username, Password, Database, uint( Timeout ), Compress, 0 ); err != nil {
    connection.Close()
    return nil, err
  } else {
    connection.Compress( Compress )
    return connection, nil
  }
}

// Connection Functions

// Connect connects to a SQLite Cloud server instance using the given arguments.
// If Connect is called on an already established connection, the old connection is closed first.
// All arguments are checked for valid values (see: CheckConnectionParameter). An error is returned if the protocol Family was not '0',
// invalid argument values where given or the connection could not be established.
func (this *SQCloud) Connect( Host string, Port int, Username string, Password string, Database string, Timeout uint, Compression string, Family int ) error {
  this.reset() // also closes an open connection

  switch err := this.CheckConnectionParameter( Host, Port, Username, Password, Database, int( Timeout ), Compression, "" );  {
  case Family != 0: return errors.New( "Invalid Protocol Family" )
  case err != nil:  return err
  default:
    this.Host     = Host
    this.Port     = Port
    this.Username = Username
    this.Password = Password
    this.Database = Database
    this.Timeout  = time.Duration( Timeout ) * time.Second
    this.Family   = Family

    return this.reconnect()
  }
}

// reconnect closes and then reopens a connection to the SQLite Cloud database server.
func (this *SQCloud) reconnect() error {
  if this.sock != nil { return nil }

  this.resetError()

  var dialer       = net.Dialer{}
  dialer.Timeout   = this.Timeout
  dialer.DualStack = true

  switch {
  case this.cert != nil:
    if tls_c, err := tls.DialWithDialer( &dialer, "tcp", net.JoinHostPort(this.Host, strconv.Itoa(this.Port)), this.cert ); err != nil {
      this.ErrorCode    = -1
      this.ErrorMessage = err.Error()
      return err
    } else {
      c := net.Conn( tls_c )
      this.sock = &c
    }
  default:
    // todo: use the dialer...
    if c, err := net.DialTimeout( "tcp", net.JoinHostPort(this.Host, strconv.Itoa(this.Port)), this.Timeout ); err != nil {
      this.ErrorCode    = -1
      this.ErrorMessage = err.Error()
      return err
    } else {
      this.sock = &c
    }
  }

  if strings.TrimSpace(this.Username) != "" {
	if err := this.Auth(this.Username, this.Password); err != nil {
		this.ErrorCode = -1
		this.ErrorMessage = err.Error()
		return err
	}
  }

  if strings.TrimSpace( this.Database ) != "" {
    this.UseDatabase( this.Database )
  }

  return nil
}

// Close closes the connection to the SQLite Cloud Database server.
// The connection can later be reopened (see: reconnect)
func (this *SQCloud) Close() error {
  var err_sock, err_psub error

  err_psub = this.psubClose()

  if this.sock != nil { err_sock = ( *this.sock ).Close() }
  this.sock = nil

  this.resetError()

  if err_sock != nil {
    this.setError( -1, err_sock.Error() )
    return err_sock
  }

  if err_psub != nil { 
    this.setError( -1, err_psub.Error() )
    return err_psub
  }
  return nil
}

// Compress enabled or disables data compression for this connection.
// If enabled, the data is compressed with the LZ4 compression algorithm, otherwise no compression is applied the data.
func (this *SQCloud) Compress( CompressMode string ) error {
  switch compression := strings.ToUpper( CompressMode ); {
  case this.sock == nil:     return errors.New( "Not connected" )
  case compression == "NO":  return this.Execute( "SET CLIENT KEY COMPRESSION TO 0" )
  case compression == "LZ4": return this.Execute( "SET CLIENT KEY COMPRESSION TO 1" )
  default:                   return errors.New( fmt.Sprintf( "Invalid method (%s)", CompressMode ) )
  }
}

// IsConnected checks the connection to the SQLite Cloud database server by sending a PING command.
// true is returned, if the connection is established and actually working, false otherwise.
func (this *SQCloud) IsConnected() bool {
  switch {
  case this.sock == nil:    return false
  case this.Ping() != nil:  return false
  default:                  return true
  }
}

// Error Methods

func (this *SQCloud) setError( ErrorCode int, ErrorMessage string ) {
  this.ErrorCode    = ErrorCode
  this.ErrorMessage = ErrorMessage
}

// resetError resets the error code and message of the last run command.
func (this *SQCloud) resetError() { this.setError( 0, "" ) }

// GetErrorCode returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud ) GetErrorCode() int { return this.ErrorCode }

// GetExtErrorCode returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud ) GetExtErrorCode() int { return this.ExtErrorCode }

// IsError checks the successful execution of the last method call / command.
// true is returned if the last command resulted in an error, false otherwise.
func (this *SQCloud) IsError() bool { return this.GetErrorCode() != 0 }

// GetErrorMessage returns the error message of the last unsuccuessful command as an error.
// nil is returned if the last command run successful.
func (this *SQCloud ) GetErrorMessage() error {
  switch this.IsError() {
  case true: return errors.New( this.ErrorMessage )
  default:   return nil
  }
}

// GetError returned the error code and message of the last unsuccessful command.
// 0 and nil is returned if the last command run successful.
func (this *SQCloud ) GetError() ( int, int, error ) { return this.GetErrorCode(), this.GetExtErrorCode(), this.GetErrorMessage() }


// Data Access Functions

// Select executes a query on an open SQLite Cloud database connection.
// If an error occurs during the execution of the query, nil and an error describing the problem is returned.
// On successful execution, a pointer to the result is returned.
func ( this *SQCloud ) Select( SQL string ) ( *Result, error ) {
  this.resetError()

  if _, err := this.sendString( SQL ); err != nil { return nil, err }

  switch result, err := this.readResult(); {
  case result == nil: return nil, errors.New( "nil" )

  case result.IsError():
    this.ErrorCode, this.ExtErrorCode, this.ErrorMessage, _ = result.GetError()
    result.Free()
    return nil, errors.New( this.ErrorMessage )

  case err != nil:
    this.ErrorCode, this.ExtErrorCode, this.ErrorMessage = 100000, 0, err.Error()
    result.Free()
    return nil, err

  default:            return result, nil
  }
}

// Execute executes the given query.
// If the execution was not successful, an error describing the reason of the failure is returned.
func (this *SQCloud) Execute( SQL string ) error {
  if result, err := this.Select( SQL ); result != nil {
    defer result.Free()

    if !result.IsOK() {
      return errors.New( "ERROR: Unexpected Result (-1)")
    }
    return err
  } else { return err }
}
