package sqlitecloud

import "fmt"
//import "os"
//import "bufio"
import "strings"
import "errors"
import "time"
//import "strconv"

import "github.com/xo/dburl"

type SQCloudKeyValues struct {
	Key 				string
	Value 			string
}

func init() {
	dburl.Register( dburl.Scheme {
		Driver: "sc", // sqlitecloud
		Generator: dburl.GenFromURL("sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&uuid=12342&compress=disabled&sslmode=disabled"),
		Transport: dburl.TransportTCP,
		Opaque: false,
		Aliases: []string{ "sqlitecloud" },
		Override: "",
	} )
}

// Creation

func New() *SQCloud {
	connection := SQCloud{ connection: nil }
	connection.reset()
	return &connection
}

// sqlitecloud://user:pass@host.com:port/dbname?timeout=10&uuid=12342&compress=disabled&sslmode=disabled&family=1
func Connect( ConnectionString string ) (*SQCloud, error) {
	u, err := dburl.Parse( ConnectionString ) // "postgresql://user:pass@localhost/mydatabase/?sslmode=disable")
	if err == nil {

		host, _	  := parseString( u.Hostname(), "localhost" )
		port, err	:= parseInt( u.Port(), 8860, 0x0001, 0xFFFF )
		if err != nil {
			return nil, errors.New( "ERROR: Invalid Port number" )
		}
		user 			:= strings.TrimSpace( u.User.Username() )
		pass, _ 	:= u.User.Password()
		database 	:= strings.TrimSpace( strings.TrimPrefix( u.Path, "/" ) )

		timeout 	:= 10
		uuid      := ""
		compress  := true
		ssl       := true
		family    := 0

		for key, value := range u.Query() {
			switch strings.ToLower( strings.TrimSpace( key ) ) {
			case "compress":
				if compress, err = parseBool( value[ 0 ], true ); err != nil {
					return nil, errors.New( "ERROR: Value for Argument 'compress' is not a boolean value" )
				}
			case "sslmode":
				if ssl, err = parseBool( value[ 0 ], true ); err != nil {
					return nil, errors.New( "ERROR: Value for Argument 'sslmode' is not a boolean value" )
				}
			case "timeout":
				if timeout, err = parseInt( value[ 0 ], 10, 0x0000, 0xFFFF ); err != nil {
					return nil, errors.New( "ERROR: Value for Argument 'timeout' is invalid" )
				}
			case "family":
				if family, err = parseInt( value[ 0 ], 0, 0, 0 ); err != nil {
					return nil, errors.New( "ERROR: Value for Argument 'family' is invalid" )
				}
			case "uuid":
				if uuid, err = parseString( value[ 0 ], "" ); err != nil {
					return nil, errors.New( "ERROR: Value for Argument 'uuid' must not be empty" )
				}
			default: // Ignore
			}
		}

		connection := New()
		err = connection.Connect( host, port, user, pass, database, timeout, family )
		if err != nil {
			connection.Close()
			return nil, err
		}

		ssl = ssl
		connection.Compress( compress )
		connection.SetUUID( uuid )	

		return connection, nil
	}
	return nil, err
}

// Connection Functions

func (this *SQCloud) reset() {
	this.Close()
	this.resetError()
	this.Host       = ""
	this.Port     	= -1
	this.Username 	= ""
	this.Password 	= ""
	this.Database 	= ""
  this.Timeout  	= -1
	this.Family   	= -1
	this.UUID       = ""
}
func (this *SQCloud) resetError() {
	this.ErrorCode    = 0
	this.ErrorMessage = ""
}

func (this *SQCloud) Connect( Host string, Port int, Username string, Password string, Database string, Timeout int, Family int ) error {
	this.reset()

	Host = strings.TrimSpace( Host )
	if Host == ""  {
		return errors.New( "ERROR: Invalid Hostname" )
	}

	if Port < 1 || Port > 0xFFFF {
		return errors.New( "ERROR: Invalid Port Number" )
	}
	if Timeout < 0 {
		return errors.New( "ERROR: Invalid Timout" )
	}
	if Family < 0 {
		return errors.New( "ERROR: Invalid Family" )
	}

	this.Host     = Host
	this.Port     = Port
	this.Username = Username
	this.Password = Password
	this.Database = Database
  this.Timeout  = Timeout
	this.Family   = Family

	return this.Reconnect()
}

func (this *SQCloud) Reconnect() error {
	this.Close()
	this.resetError()

	// fmt.Printf( "Reconnect with: %s, %d, %s, %s, %s, %d, %d\r\n",  this.Host, this.Port, this.Username, this.Password, this.Database, this.Timeout, this.Family )

	this.connection = CConnect( this.Host, this.Port, this.Username, this.Password, this.Database, this.Timeout, this.Family )

	if this.connection != nil {
		this.ErrorCode  	= this.CGetErrorCode()
		this.ErrorMessage = this.CGetErrorMessage()

		if !this.IsError() {
			if strings.TrimSpace(	this.UUID ) != "" {
				this.SetUUID( this.UUID )
			}
		}
	
		if !this.IsError() {
			if strings.TrimSpace(	this.Database ) != "" {
				this.UseDatabase( this.Database )
			}
		}

	} else {

		this.ErrorCode 		= 666
		this.ErrorMessage = "Not enoght memory to allocate a SQCloudConnection."
	}

	if this.IsError() {
		err := errors.New( fmt.Sprintf( "ERROR CONNECTION TO %s: %s (%d)", this.Host, this.ErrorMessage, this.ErrorCode ) )
		this.Close()
		return err
	}

	return nil
}

func (this *SQCloud) Close() {
	if this.IsConnected() {
		this.CDisconnect()
	}
	this.connection = nil
	this.resetError()
}

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

func (this *SQCloud) IsError() bool {
	return this.ErrorCode != 0
}
func (this *SQCloud ) GetErrorCode() int {
	return this.ErrorCode
}
func (this *SQCloud ) GetErrorMessage() error {
	if this.IsError() {
		return errors.New( this.ErrorMessage )
	}
	return nil
}
func (this *SQCloud ) GetError() ( int, error ) {
	return this.GetErrorCode(), this.GetErrorMessage()
}

// Connection Info Methods

func (this *SQCloud ) GetUUID() string {
	return this.UUID // this.CGetCloudUUID()
}
func (this *SQCloud) SetUUID( UUID string ) error {
	this.UUID = UUID
	return this.Execute( fmt.Sprintf( "SET CLIENT UUID TO %s", SQCloudEnquoteString( UUID ) ) )
}



// ResultSet Methods

func (this *SQCloudResult ) GetType() int {
	return this.CGetResultType()
}
func (this *SQCloudResult ) GetBuffer() string {
	return this.CGetResultBuffer()
}
func (this *SQCloudResult ) GetLength() uint {
	return this.CGetResultLen()
}
func (this *SQCloudResult ) GetMaxLength() uint32 {
	return this.CGetMaxLen()
}
func (this *SQCloudResult ) Free() {
	this.CFree()
}
func (this *SQCloudResult ) IsOK() bool {
	return this.Type == RESULT_OK
}
func (this *SQCloudResult ) GetNumberOfRows() uint {
	return this.CGetRows()
}
func (this *SQCloudResult ) GetNumberOfColumns() uint {
	return this.CGetColumns()
}
func (this *SQCloudResult ) GetValueType( Row uint, Column uint ) int {
	return this.CGetValueType( Row, Column )
}
func (this *SQCloudResult ) GetColumnName( Row uint, Column uint ) string {
	return this.CGetColumnName( Row, Column )
}
func (this *SQCloudResult ) GetStringValue( Row uint, Column uint ) string {
	return this.CGetStringValue( Row, Column )
} 
func (this *SQCloudResult ) GetInt32Value( Row uint, Column uint ) int32 {
	return this.CGetInt32Value( Row, Column )
} 
func (this *SQCloudResult ) GetInt64Value( Row uint, Column uint ) int64 {
	return this.CGetInt64Value( Row, Column )
} 
func (this *SQCloudResult ) GetFloat32Value( Row uint, Column uint ) float32 {
	return this.CGetFloat32Value( Row, Column )
} 
func (this *SQCloudResult ) GetFloat64Value( Row uint, Column uint ) float64 {
	return this.CGetFloat64Value( Row, Column )
}
func (this *SQCloudResult ) Dump( MaxLine uint ) {
	this.CDump( MaxLine )
}

// Additional ResultSet Methods

func (this *SQCloudResult ) GetSQLDateTime( Row uint, Column uint ) time.Time {
	datetime, _ := time.Parse( "2006-01-02 15:04:05", this.CGetStringValue( Row, Column ) )
	return datetime
} 

// Data Access Functions

func (this *SQCloud) Select( SQL string ) (*SQCloudResult, error) {
	this.resetError()

  result           := this.CExec( SQL )
	this.ErrorCode  	= this.CGetErrorCode()
	this.ErrorMessage = this.CGetErrorMessage()

	if result != nil {
		result.Rewind()
		result.Type 				= result.GetType()
		result.ErrorCode 		= this.ErrorCode
		result.ErrorMessage = this.ErrorMessage
		result.Rows         = result.GetNumberOfRows()
	} else {
		return nil, errors.New( "ERROR: Could not execute SQL command (-1)" )
	}

	if this.IsError() {
		result.Free()
		return nil, errors.New( fmt.Sprintf( "ERROR: %s (%d)", this.CGetErrorMessage(), this.CGetErrorCode() ) )
	}
	return result, nil // nil, nil or *result, nil
}

// Additional Data Access Functions

func (this *SQCloudResult ) IsError() bool {
	return this.Type == RESULT_ERROR
}
func (this *SQCloudResult ) IsNull() bool {
	return this.Type == RESULT_NULL
}
func (this *SQCloudResult ) IsJson() bool {
	return this.Type == RESULT_JSON
}
func (this *SQCloudResult ) IsString() bool {
	return this.Type == RESULT_STRING
}
func (this *SQCloudResult ) IsInteger() bool {
	return this.Type == RESULT_INTEGER
}
func (this *SQCloudResult ) IsFloat() bool {
	return this.Type == RESULT_FLOAT
}
func (this *SQCloudResult ) IsRowSet() bool {
	return this.Type == RESULT_ROWSET
}
func (this *SQCloudResult ) IsTextual() bool {
	return this.IsJson() || this.IsString() || this.IsInteger() || this.IsFloat()
}

func (this *SQCloudResult ) Rewind() {
	this.row = 0
}
func (this *SQCloudResult ) GetNextRow() bool {
	if !this.IsEOF() {
		this.row++
		return true
	}
	return false
}
func (this *SQCloudResult ) IsEOF() bool {
	return this.row >= this.Rows
}

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





// Pub/Sub

func (this *SQCloud ) SetPubSubOnly() *SQCloudResult {
	return this.CSetPubSubOnly()
}
// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);






// Helper functions

func SQCloudEnquoteString( Token string ) string {
	Token = strings.Replace( Token, "\"", "\"\"", -1 )
	if strings.Contains( Token, " " ) {
		return fmt.Sprintf( "\"%s\"", Token )
	}
	return Token
}