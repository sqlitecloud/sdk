//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/10/11
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Simple PSUB Test
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloudtest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	sqlitecloud "github.com/sqlitecloud/sdk"
)

const testPubsubChannelName = "TestPubsubChannel"

var testPubsubMessage = map[string]string{"msg_id": "12345", "msg_content": "this is the content"}

func TestPubsub(t *testing.T) {
	db1, err := sqlitecloud.Connect(testConnectionString)
	if err != nil {
		t.Fatal("Connect 1: ", err.Error())
	}
	defer db1.Close()

	db2, err := sqlitecloud.Connect(testConnectionString)
	if err != nil {
		t.Fatal("Connect 2: ", err.Error())
	}
	defer db2.Close()

	ch := make(chan string, 1)
	db1.Callback = func(db *sqlitecloud.SQCloud, jsonString string) {
		ch <- jsonString
	}

	if _, err := db1.ListChannels(); err != nil {
		t.Fatal("ListChannels: ", err.Error())
	}

	if err := db1.CreateChannel(testPubsubChannelName, true); err != nil {
		t.Fatal("CreateChannel: ", err.Error())
	}

	if channels, err := db1.ListChannels(); err != nil {
		t.Fatal("ListChannels: ", err.Error())
	} else if !contains(channels, testPubsubChannelName) {
		t.Fatal("ListChannels: ", fmt.Sprintf("Channel %s not found in LIST CHANNELS", testPubsubChannelName))
	}

	if err := db1.Listen(testPubsubChannelName); err != nil {
		t.Fatal("Listen: ", err.Error())
	}

	jsonStr, err := json.Marshal(testPubsubMessage)
	if err != nil {
		t.Fatal("Marshal: ", err.Error())
	}

	if err := db2.SendNotificationMessage(testPubsubChannelName, string(jsonStr)); err != nil {
		t.Fatal("SendNotificationMessage: ", err.Error())
	}

	select {
	case receivedStr := <-ch:
		var receivedMap map[string]interface{}
		json.Unmarshal([]byte(receivedStr), &receivedMap)
		payload, found := receivedMap["payload"]
		if !found {
			t.Fatal("Received notification: missing payload")
		}

		payloadStr, ok := payload.(string)
		if !ok {
			t.Fatal("Received notification: invalid payload")
		}
		var messageMap map[string]string
		json.Unmarshal([]byte(payloadStr), &messageMap)

		if !reflect.DeepEqual(messageMap, testPubsubMessage) {
			t.Fatalf("Received notification: Expected %v, Got %v", testPubsubMessage, messageMap)
		}

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout")
	}

	if err := db1.DropChannel(testPubsubChannelName); err != nil {
		t.Fatal("DropChannel: ", err.Error())
	}
}
