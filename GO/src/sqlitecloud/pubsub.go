//
//                    ////              SQLite Cloud
//        ////////////  ///             
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 0.0.1
//     //             ///   ///  ///    Date        : 2021/08/13
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

package sqlitecloud

//import "fmt"
//import "os"
//import "bufio"
//import "strings"
//import "errors"
//import "time"
//import "strconv"

// Pub/Sub

// Connection Info Methods

// GetUUID returns the UUID as string
func (this *SQCloud ) GetUUID() string {
  return this.UUID // this.CGetCloudUUID()
}

// SetPubSubOnly
// BUG(andreas): TODO, postponed by Marco
func (this *SQCloud ) SetPubSubOnly() *SQCloudResult {
  return this.CSetPubSubOnly()
}


// void SQCloudSetPubSubCallback (SQCloudConnection *connection, SQCloudPubSubCB callback, void *data);
