//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/10/12
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///     
//   ///     //////////   ///  ///      Description : GO Methods related to the
//   ////                ///  ///                     SQCloud class for handling
//     ////     //////////   ///                      asynchronous communication.
//        ////            ////          
//          ////     /////              
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2


/*
    REM LISTEN CHANNEL:
    REM LISTEN TabName  = Listen on WRITES on Table "TabName" in this database  (execute 1...n times)
    REM LISTEN *        = Listen on WRITES on All Tables in this database       (execute 1...n times)
    REM LISTEN ChanName = Listen on NOTIFYs on the Channel ChanName             (execute 1...n times)
    
    REM UNLISTEN ChanName|TabName = Unregisteres a previous registration
    REM UNLISTEN *                = Unregisteres ALL (=TabName,*,ChanName) registrations
    
    REM NOTIFY ChanName           = NOTIFY ChanName ""
    REM NOTIFY ChanName <STRING-PAYLOAD>
    
    REM LISTEN
    10 SEND "LISTEN *"
    20 RECEIVE "|79 PAUTH a365efef-cfb7-4672-8ed4-45a489ddb194 9230b8d8-93dc-4edc-bcaf-cc118fe32d4d"
    30 IF NO 2.socket IS THERE: OPEN 2.socket
    40 SEND "PAUTH a365efef-cfb7-4672-8ed4-45a489ddb194 9230b8d8-93dc-4edc-bcaf-cc118fe32d4d"
    50 RECEIVE "OK"
    60 START 2.thread
    2.10 IF 2.socket LOST CONNECTION: CLOSE 2.socket and TERMINATE 2.thread
    2.20 RECEIVE "#LEN json"
    2.30 CALL callback_function WITH json
    2.40 GOTO 2.10
    
    ON_CLOSE EVENT: 
    10 IF 2.socket IS CONNECTED: CLOSE 2.socket
    20 IF main.socket IS CONNECTED: CLOSE main.socket
*/

package sqlitecloud

import "fmt"

// GetUUID returns the UUID as string
func (this *SQCloud ) GetUUID() string {
  return this.uuid // this.CGetCloudUUID()
}

// psubClose closes the PSUB connection to the SQLite Cloud Database server.
func (this *SQCloud) psubClose() error {
  var err error = nil

  if this.psub != nil  { err = ( *this.psub ).Close() }
  this.psub = nil

  return err
}


// Listen subscribes this connection to the specified Channel.
func (this *SQCloud) Listen( Channel string ) error { // add a call back function...
  return this.Execute( fmt.Sprintf( "LISTEN %s", SQCloudEnquoteString( Channel ) ) )
}

// Notify sends a wakeup call to the channel Channel
func (this *SQCloud) Notify( Channel string ) error {
  return this.Execute( fmt.Sprintf( "NOTIFY %s", SQCloudEnquoteString( Channel ) ) )
}

// SendNotificationMessage sends the message Message to the channel Channel
func (this *SQCloud) SendNotificationMessage( Channel string, Message string ) error {
  return this.Execute( fmt.Sprintf( "NOTIFY %s %s", SQCloudEnquoteString( Channel ), SQCloudEnquoteString( Message ) ) )
}

// Unlisten unsubsribs this connection from the specified Channel.
func (this *SQCloud) Unlisten( Channel string ) error {
  return this.Execute( fmt.Sprintf( "UNLISTEN %s", SQCloudEnquoteString( Channel ) ) )
}