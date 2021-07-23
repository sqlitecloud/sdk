package sqlitecloud

import "fmt"
//import "os"
//import "bufio"
//import "strings"
import "errors"


func Connect( Host string, Port int, Username string, Password string, Database string, Timeout int, Family int ) (*SQCloud, error) {
	connection := bridge_Connect( Host, Port, Username, Password, Database, Timeout, Family )
	if connection != nil {
		if connection.bridge_IsError() {
			defer connection.bridge_Disconnect()
			return nil, errors.New( fmt.Sprintf( "ERROR CONNECTION TO %s: %s (%d)", Host, connection.bridge_GetErrorMessage(), connection.bridge_GetErrorCode() ) )
		} 
	}
	return connection, nil
}

func (this *SQCloud) Close() {
	this.bridge_Disconnect()
}

func (this *SQCloud ) GetCloudUUID() string {
	return this.bridge_GetCloudUUID()
}



func (this *SQCloudResult ) GetNumRows() uint {
	return this.bridge_GetRows()
}

func (this *SQCloudResult ) Free() {
	this.bridge_Free()
}

func (this *SQCloud) Execute( Command string ) (*SQCloudResult, error) {
  result := this.bridge_Exec( Command )
	if this.bridge_IsError() {
		return nil, errors.New( fmt.Sprintf( "ERROR" ) )
	}
	return result, nil
}