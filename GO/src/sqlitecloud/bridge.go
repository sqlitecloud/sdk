package sqlitecloud

// #cgo CFLAGS: -Wno-multichar -I../../../C
// #cgo LDFLAGS: -L. -lsqcloud -ldl
// #include <stdlib.h>
// #include "sqcloud.h"
import "C"
import "unsafe"
//import "errors"
// import "fmt"
// import "reflect"

type SQCloud struct {
  connection 	*C.struct_SQCloudConnection
}

type SQCloudResult struct {
	result *C.struct_SQCloudResult
}

// SQCloudResType
const RESULT_OK 			= C.RESULT_OK
const RESULT_ERROR 		= C.RESULT_ERROR
const RESULT_STRING 	= C.RESULT_STRING
const RESULT_INTEGER 	= C.RESULT_INTEGER
const RESULT_FLOAT 		= C.RESULT_FLOAT
const RESULT_ROWSET 	= C.RESULT_ROWSET
const RESULT_NULL 		= C.RESULT_NULL
const RESULT_JSON 		= C.RESULT_JSON

// SQCloudValueType
const VALUE_INTEGER 	= C.VALUE_INTEGER
const VALUE_FLOAT 		= C.VALUE_FLOAT
const VALUE_TEXT 			= C.VALUE_TEXT
const VALUE_BLOB 			= C.VALUE_BLOB
const VALUE_NULL 			= C.VALUE_NULL

// SQCloudConnection *SQCloudConnect (const char *hostname, int port, SQCloudConfig *config);
func bridge_Connect( Host string, Port int, Username string, Password string, Database string, Timeout int, Family int ) *SQCloud {
	conf := C.struct_SQCloudConfigStruct{}
	conf.username = C.CString( Username )
	conf.password = C.CString( Password )
	conf.database = C.CString( Database )
	conf.timeout  = C.int( Timeout )
	conf.family   = C.int( Family )

	cHost := C.CString( Host )

	connection := SQCloud{ connection: C.SQCloudConnect( cHost, C.int( Port ), nil ) }
	
	C.free( unsafe.Pointer( cHost ) )
	C.free( unsafe.Pointer( conf.database ) )
	C.free( unsafe.Pointer( conf.password ) )
	C.free( unsafe.Pointer( conf.username ) )

	if connection.connection == nil {
		return nil
	}

  return &connection
}
// void SQCloudDisconnect (SQCloudConnection *connection);
func (this *SQCloud ) bridge_Disconnect() {
	if this.connection != nil {
		C.SQCloudDisconnect( this.connection )
		this.connection = nil
	}
}
// char *SQCloudUUID (SQCloudConnection *connection);
func (this *SQCloud ) bridge_GetCloudUUID() string {
	 return C.GoString( C.SQCloudUUID( this.connection ) )
}

//bool SQCloudIsError (SQCloudConnection *connection);
func (this *SQCloud ) bridge_IsError() bool {
	return bool( C.SQCloudIsError( this.connection ) )
}
//int SQCloudErrorCode (SQCloudConnection *connection);
func (this *SQCloud ) bridge_GetErrorCode() int {
	return int( C.SQCloudErrorCode( this.connection ) )
}
//const char *SQCloudErrorMsg (SQCloudConnection *connection);
func (this *SQCloud ) bridge_GetErrorMessage() string {
	return C.GoString( C.SQCloudErrorMsg( this.connection ) )
}

// SQCloudResult *SQCloudExec (SQCloudConnection *connection, const char *command);
func (this *SQCloud ) bridge_Exec( Command string ) *SQCloudResult {
	cCommand := C.CString( Command )
  defer C.free( unsafe.Pointer( cCommand ) )

	println( "exec ("+Command+").." )

	result := SQCloudResult{ result: C.SQCloudExec( this.connection, cCommand ) }
	if result.result == nil {
		return nil
	}
	return &result
}
// SQCloudResult *SQCloudSetPubSubOnly (SQCloudConnection *connection);
func (this *SQCloud ) bridge_SetPubSubOnly() *SQCloudResult {
	result := SQCloudResult{ result: C.SQCloudSetPubSubOnly( this.connection ) }
	
	if result.result == nil {
		return nil
	}

	return &result
}
// SQCloudResType SQCloudResultType (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetResultType() int {
	return int( C.SQCloudResultType( this.result ) )
}
// uint32_t SQCloudResultLen (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetResultLen() uint {
	return uint( C.SQCloudResultLen( this.result ) )
}
// char *SQCloudResultBuffer (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetResultBuffer() string {
	return C.GoString( C.SQCloudResultBuffer( this.result ) )
}
// void SQCloudResultFree (SQCloudResult *result);
func (this *SQCloudResult ) bridge_Free() {
	C.SQCloudResultFree( this.result )
}
// bool SQCloudResultIsOK (SQCloudResult *result);
func (this *SQCloudResult ) bridge_IsOK() bool {
	return bool( C.SQCloudResultIsOK( this.result ) )
}
// SQCloudValueType SQCloudRowsetValueType (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) bridge_GetValueType( Row uint, Column uint ) int {
	return int( C.SQCloudRowsetValueType( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// char *SQCloudRowsetColumnName (SQCloudResult *result, uint32_t col, uint32_t *len);
func (this *SQCloudResult ) bridge_GetColumnName( Row uint, Column uint ) string {
	var len C.uint32_t = 0
	return C.GoStringN( C.SQCloudRowsetColumnName( this.result, C.uint( Column ), &len ), C.int( len ) )
}
// uint32_t SQCloudRowsetRows (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetRows() uint {
	return uint( C.SQCloudRowsetRows( this.result ) )
}
// uint32_t SQCloudRowsetCols (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetColumns() uint {
	return uint( C.SQCloudRowsetCols( this.result ) )
}
// uint32_t SQCloudRowsetMaxLen (SQCloudResult *result);
func (this *SQCloudResult ) bridge_GetMaxLen() uint32 {
	return uint32( C.SQCloudRowsetMaxLen( this.result ) )
}
// char *SQCloudRowsetValue (SQCloudResult *result, uint32_t row, uint32_t col, uint32_t *len);
func (this *SQCloudResult ) bridge_GetStringValue( Row uint, Column uint ) string {
	var len C.uint32_t = 0
  return C.GoStringN( C.SQCloudRowsetValue( this.result, C.uint32_t( Row ), C.uint32_t( Column ), &len ), C.int( len ) ) // Problem: NULL Pointer in return
}
// int32_t SQCloudRowsetInt32Value (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) bridge_GetInt32Value( Row uint, Column uint ) int32 {
	return int32( C.SQCloudRowsetInt32Value( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// int64_t SQCloudRowsetInt64Value (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) bridge_GetInt64Value( Row uint, Column uint ) int64 {
	return int64( C.SQCloudRowsetInt64Value( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// float SQCloudRowsetFloatValue (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) bridge_GetFloat32Value( Row uint, Column uint ) float32 {
	return float32( C.SQCloudRowsetFloatValue( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// double SQCloudRowsetDoubleValue (SQCloudResult *result, uint32_t row, uint32_t col);
func (this *SQCloudResult ) bridge_GetFloat64Value( Row uint, Column uint ) float64 {
	return float64( C.SQCloudRowsetDoubleValue( this.result, C.uint( Row ), C.uint( Column ) ) )
}
// void SQCloudRowsetDump (SQCloudResult *result, uint32_t maxline);
func (this *SQCloudResult ) bridge_Dump( MaxLine uint ) {
	 C.SQCloudRowsetDump( this.result, C.uint( MaxLine ) )
}


// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);