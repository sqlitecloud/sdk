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
