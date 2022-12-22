//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/08/16
//    ///             ///   ///  ///    Author      : Andrea Donetti
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	sqlitecloud "github.com/sqlitecloud/sdk"

	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gobwas/glob"
	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

const (
	SuccessStatus string = "success"
	ErrorStatus   string = "error"
)

// ----------------------------------------------------------------------------
// Struct definitions
// ----------------------------------------------------------------------------

type ApiResponse struct {
	Status string      `json:"status"`         // mandatory, "success" or "error"
	Id     string      `json:"id"`             // mandatory
	Type   string      `json:"type"`           // mandatory
	Data   interface{} `json:"data,omitempty"` // optional
}

type ApiResponseDataError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ApiResponseDataPAuth struct {
	UUID    string `json:"uuid"`
	Secret  string `json:"secret"`
	Channel string `json:"channel"`
}

type ApiResponseDataChannel struct {
	Channel string `json:"channel"`
}

type ApiConnection struct {
	ws         *websocket.Conn // main websocket used by the client to send command and receive command responses
	pubsubws   *websocket.Conn // secondary websocket used to send pubsub notification to the client
	pubsubwsmu sync.RWMutex    // mutex for thread-safe access to pubsubws, see GetPubsubws and SetPubsubws functions
	sqlcconn   *sqlitecloud.SQCloud
}

// thread-safe read access to pubsubws field of ApiConnection struct
func (this *ApiConnection) GetPubsubws() *websocket.Conn {
	this.pubsubwsmu.RLock()
	defer this.pubsubwsmu.RUnlock()
	return this.pubsubws
}

// thread-safe write access to pubsubws field of ApiConnection struct
func (this *ApiConnection) SetPubsubws(ws *websocket.Conn) {
	this.pubsubwsmu.Lock()
	defer this.pubsubwsmu.Unlock()
	this.pubsubws = ws
}

// ----------------------------------------------------------------------------
// ApiConnection Map with thread-safe access functions
// ----------------------------------------------------------------------------

// Global map: ApiConenctions by uuid
var m map[string]*ApiConnection // the key is the uuid from SQCloud conneciton
var mmu = sync.RWMutex{}

// Thread-safe function to get ApiConenction by uuid
func getPubsubConn(uuid string) (*ApiConnection, bool) {
	mmu.RLock()
	defer mmu.RUnlock()
	v, found := m[uuid]
	return v, found
}

// Thread-safe function to set ApiConenction for a uuid
func setPubsubConn(uuid string, conn *ApiConnection) {
	mmu.Lock()
	defer mmu.Unlock()
	m[uuid] = conn
}

// ----------------------------------------------------------------------------
// SQCloud Pubsub Callback
// ----------------------------------------------------------------------------

func pubsubCallback(conn *sqlitecloud.SQCloud, payload string) {
	uuid, _ := conn.GetPAuth()
	pubsubConn, found := getPubsubConn(uuid)
	if found {
		pubsubws := pubsubConn.GetPubsubws()
		if pubsubws != nil {
			if err := pubsubws.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
				pubsubws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			}
		}
	}
}

// ----------------------------------------------------------------------------
// Main functions
// ----------------------------------------------------------------------------

var apiWebsocketUpgrader websocket.Upgrader // use default options

func init() {
	initializeSQLiteWeb()

	apiWebsocketUpgrader = websocket.Upgrader{}
	m = make(map[string]*ApiConnection)

	originChecker := glob.MustCompile("{https://*.sqlitecloud.io,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
	localhostChecker := glob.MustCompile("{https://localhost:*,https://localhost,http://localhost:*,http://localhost}")

	apiWebsocketUpgrader.CheckOrigin = func(r *http.Request) bool {
		o := r.Header.Get("Origin")
		allowed := originChecker.Match(o)

		// TODO: only for debug purposes
		if !allowed {
			allowed = localhostChecker.Match(o)
		}
		if !allowed {
			SQLiteWeb.Logger.Debugf("CheckOrigin: not allowed origin -%s-", o)
		}
		return allowed
	}
}

func getSQCloudConnection(request *http.Request) (*sqlitecloud.SQCloud, error) {
	// get vars
	v := mux.Vars(request)
	// get projectID
	projectID, found := v["projectID"]
	if !found {
		err := fmt.Errorf("serveWebsocket: missing projectID")
		return nil, err
	}

	// get key
	qvars := request.URL.Query()
	apikeys, found := qvars["apikey"]
	if !found && len(apikeys) == 0 {
		err := fmt.Errorf("serveWebsocket: missing apikey")
		return nil, err
	}
	apikey := apikeys[0]

	SQLiteWeb.Logger.Debugf("serveWebsocket: project %s apikey %s", projectID, apikey)

	// don't use the pool, create a new connection only for this websocket
	// first get a connection url for one of the servers of the specified projectID
	connstring, err := apicm.getNextServer(projectID, false)
	if err != nil {
		err = fmt.Errorf("serveWebsocket: error on getNextServer: %s", err.Error())
		return nil, err
	}

	// remove the user from connection url, the admin user was automatically added with
	connurl, err := url.Parse(connstring)
	if err != nil {
		err = fmt.Errorf("serveWebsocket: error in connection url %s", err.Error())
		return nil, err
	}
	connurl.User = nil

	// add api key to connection url
	values := connurl.Query()
	values.Add("apikey", apikey)
	connurl.RawQuery = values.Encode()

	// try to connect with the connection url
	connection, err := sqlitecloud.Connect(connurl.String())
	if err != nil {
		if connection != nil {
			connection.Close()
			connection = nil
		}
		err = fmt.Errorf("serveWebsocket: error on connect %s", connurl.String())
		return nil, err
	}

	return connection, nil
}

func (this *Server) serveApiWebsocket(writer http.ResponseWriter, request *http.Request) {
	var connection *sqlitecloud.SQCloud = nil

	connection, err := getSQCloudConnection(request)
	if err != nil {
		SQLiteWeb.Logger.Error(err.Error())
		return
	}

	connection.Callback = pubsubCallback

	wsconn, err := apiWebsocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Error("serveWebsocket: upgrade error: ", err)
		connection.Close()
		return
	}
	defer wsconn.Close()

	uuid := ""
	for {
		var result *sqlitecloud.Result = nil
		var responsedata interface{}

		// read the command (JSON) message from the client
		messageType, message, err := wsconn.ReadMessage()
		SQLiteWeb.Logger.Debugf("serveWebsocket: ReadMessage %d %s", messageType, message)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				SQLiteWeb.Logger.Debug("serveApiWebsocket: read error: ", err)
				wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			}
			break
		}

		if messageType != websocket.TextMessage {
			SQLiteWeb.Logger.Debug("serveWebsocket: wrong message type")
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "wrong message type"), time.Now().Add(1*time.Second))
			break
		}

		var mmap map[string]interface{}
		json.Unmarshal(message, &mmap)

		// get command type
		t, ok := mmap["type"].(string)
		if !ok {
			errstring := fmt.Sprintf("wrong request: type")
			SQLiteWeb.Logger.Debugf("serveWebsocket: %s", errstring)
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, errstring), time.Now().Add(1*time.Second))
			break
		}
		id, ok := mmap["id"].(string)
		if !ok {
			errstring := fmt.Sprintf("wrong request: id")
			SQLiteWeb.Logger.Debugf("serveWebsocket: %s", errstring)
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, errstring), time.Now().Add(1*time.Second))
			break
		}

		// get command options and exec the command using the opened connection
		switch t {
		case "exec":
			command, _ := mmap["command"].(string)
			result, err = connection.Select(command)

		case "listen":
			channel, _ := mmap["channel"].(string)

			err = connection.Listen(channel)
			puuid, psecret := connection.GetPAuth()
			if puuid == "" || psecret == "" {
				err = fmt.Errorf("pubsub: error during authentication")
			} else {
				pubsubConn, found := getPubsubConn(uuid)
				if !found || pubsubConn.GetPubsubws() == nil {
					pubsubConn := ApiConnection{ws: wsconn, pubsubws: nil, sqlcconn: connection}
					setPubsubConn(puuid, &pubsubConn)
					responsedata = ApiResponseDataPAuth{UUID: puuid, Secret: psecret, Channel: channel}
				} else {
					responsedata = ApiResponseDataChannel{Channel: channel}
				}
			}
			uuid = puuid

		case "unlisten":
			channel, _ := mmap["channel"].(string)
			err = connection.Unlisten(channel)
			if err == nil {
				responsedata = ApiResponseDataChannel{Channel: channel}
			}

		case "notify":
			channel, _ := mmap["channel"].(string)
			payload, _ := mmap["payload"].(string)
			err = connection.SendNotificationMessage(channel, payload)

		default:
			SQLiteWeb.Logger.Debug("serveWebsocket: wrong json type")
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "wrong json type"), time.Now().Add(1*time.Second))
			break
		}

		if err != nil && connection.ErrorCode >= 100000 {
			// internal error (the SDK cannot write to or read from the connection)
			// for example:
			// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
			// - 100003 Internal Error: sendString (%s){
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			break
		}

		// prepare the JSON response
		response := ApiResponse{Id: id, Type: t}
		if err != nil {
			response.Status = ErrorStatus
			response.Data = ApiResponseDataError{Code: connection.ErrorCode, Message: err.Error()}
		} else {
			response.Status = SuccessStatus
			if responsedata != nil {
				response.Data = responsedata
			} else if result != nil {
				response.Data, _ = ResultToObj(result)
			}
		}

		if err = wsconn.WriteJSON(response); err != nil {
			wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			break
		}
		// if jResponse, err := json.Marshal(response); err == nil {
		// 	if err = wsconn.WriteMessage(websocket.TextMessage, jResponse); err != nil {
		// 		wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
		// 		break
		// 	}
		// } else {
		// 	wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal Error"), time.Now().Add(1*time.Second))
		// 	break
		// }
	}

	if uuid != "" {
		pubsubConn, found := getPubsubConn(uuid)
		if found {
			// there is a pubsubConn object, reset the websocket
			pubsubConn.ws = nil
			closePubsubConnIfNeeded(pubsubConn, uuid)
		} else {
			// pubsubConn not found
			connection.Close()
		}
	} else {
		// there is no pubsub websocket
		connection.Close()
	}
}

func closePubsubConnIfNeeded(pubsubConn *ApiConnection, uuid string) {
	// if both websocket has been closed
	// then close the sqlitecloud connection and destroy the pubsubConn object
	if pubsubConn.ws == nil && pubsubConn.GetPubsubws() == nil {
		pubsubConn.sqlcconn.Close()
		setPubsubConn(uuid, nil)
	}
}

func closePubsubConn(pubsubConn *ApiConnection, uuid string) {
	pubsubConn.sqlcconn.Close()
	setPubsubConn(uuid, nil)
}

func (this *Server) serveApiWebsocketPubsub(writer http.ResponseWriter, request *http.Request) {
	// get key
	qvars := request.URL.Query()
	uuids, found := qvars["uuid"]
	if !found && len(uuids) == 0 {
		SQLiteWeb.Logger.Debug("serveApiWebsocketPubsub: missing apikey")
		return
	}
	uuid := uuids[0]

	secrets, found := qvars["secret"]
	if !found && len(secrets) == 0 {
		SQLiteWeb.Logger.Debug("serveApiWebsocketPubsub: missing apikey")
		return
	}
	secret := secrets[0]

	pubsubConn, found := getPubsubConn(uuid)
	if !found {
		SQLiteWeb.Logger.Debugf("serveApiWebsocketPubsub: invalid uuid: %s, m: %s", uuid, m)
		return
	}

	_, savedsecret := pubsubConn.sqlcconn.GetPAuth()
	if !found || savedsecret != secret {
		SQLiteWeb.Logger.Debug("serveApiWebsocketPubsub: pauth failed")
		return
	}

	wsconn, err := apiWebsocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Error("serveApiWebsocketPubsub: upgrade error: ", err)
		return
	}
	defer wsconn.Close()
	pubsubConn.SetPubsubws(wsconn)

	for {
		// read the command (JSON) message from the client
		messageType, message, err := wsconn.ReadMessage()
		SQLiteWeb.Logger.Debugf("serveApiWebsocketPubsub: ReadMessage %d %s", messageType, message)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				SQLiteWeb.Logger.Debug("serveApiWebsocketPubsub: read error: ", err)
				wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			}
			break
		}

		// in a pubsub websocket the client cannot sand messages after the pubsub authentication, it is only used for notification
		err = fmt.Errorf("serveApiWebsocketPubsub: invalid message on psub websocket")
		SQLiteWeb.Logger.Debug(err.Error())
		wsconn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
	}

	// reset the websocket in the pubsubConn object
	SQLiteWeb.Logger.Debug("serveApiWebsocketPubsub: closing pubsub ws %v", pubsubConn)
	pubsubConn.SetPubsubws(nil)
	closePubsubConnIfNeeded(pubsubConn, uuid)
}

// ----------------------------------------------------------------------------
// Test Client
// ----------------------------------------------------------------------------

func apiWebsocketTestClient(w http.ResponseWriter, r *http.Request) {
	SQLiteWeb.Logger.Debugf("apiWebsocketTestClient")
	apiWebsocketTestClientTemplate.Execute(w, r.Host)
}

var apiWebsocketTestClientTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");

	var ws;
	var wsPubsub;
    var print = function(message) {
		console.log(message);
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

	var id = 1;

	function connectPubsub(response) {
		if (wsPubsub) {
			print("wsPubsub already existing");
            return false;
        }

		url = "wss://" + "{{.}}" + "/api/v1/wspsub?uuid=" + encodeURIComponent(response.uuid) + "&secret=" + encodeURIComponent(response.secret)
		
		print("PUBSUB CONNECT " + url);

        wsPubsub = new WebSocket(url);
        wsPubsub.onopen = function(evt) {
			print("PUBSUB OPEN");
        }
		
        wsPubsub.onclose = function(evt) {
			print("PUBSUB CLOSE (code:" + evt.code + ")");
            wsPubsub = null;
        }

        wsPubsub.onmessage = function(evt) {
            print("PUBSUB NOTIFICATION: " + evt.data);
			print("PUBSUB PAYLOAD: " + JSON.parse(evt.data).payload);
			print("PUBSUB PAYLOAD JSON MSG: " + JSON.parse(JSON.parse(evt.data).payload).msg);
        }

        wsPubsub.onerror = function(evt) {
            print("PUBSUB WebSocket ERROR: ");
        }
        return false;
	}

	document.getElementById("connect").onclick = function(evt) {
        if (ws) {
			print("ws already existing");
            return false;
        }
		
		var projectfield = document.getElementById('project');
		if (!projectfield.value) return false;

		var apikeyfield = document.getElementById('apikey');

		url = "wss://" + "{{.}}" + "/api/v1/" + projectfield.value + "/ws?apikey=" + apikeyfield.value
		
		print("CONNECT " + url);

        ws = new WebSocket(url);
        ws.onopen = function(evt) {
			print("OPEN");
        }
        ws.onclose = function(evt) {
			print("CLOSE (code:" + evt.code + ")");
            ws = null;
        }
        ws.onmessage = function(evt) {
			var obj = JSON.parse( evt.data );
			if (obj.type === "listen") {
				connectPubsub(obj.data)
			}
            print("RECEIVED RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("WebSocket ERROR: ");
        }
        return false;
    };

	document.getElementById("exec").onclick = function(evt) {
        if (!ws) {
            return false;
        }

		var c = document.getElementById('command');
		if (!c.value) return false;
		print("Exec: " + c.value);

		var obj = new Object();
   		obj.type = "exec";
   		obj.command  = c.value;
		obj.id = id.toString(); 
		id++;

   		var jsonString= JSON.stringify(obj);
		print("Send: " + jsonString);
		ws.send(JSON.stringify(obj))
        return false;
    };

	document.getElementById("listen").onclick = function(evt) {
        if (!ws) {
            return false;
        }

		var c = document.getElementById('channel');
		if (!c.value) return false;
		print("Listen: " + c.value);

		var obj = new Object();
   		obj.type = "listen";
   		obj.channel  = c.value;
		obj.id = id.toString(); 
		id++;

   		var jsonString= JSON.stringify(obj);
		ws.send(JSON.stringify(obj))
        return false;
    };

	document.getElementById("unlisten").onclick = function(evt) {
        if (!ws) {
            return false;
        }

		var c = document.getElementById('channel');
		if (!c.value) return false;
		print("Unlisten: " + c.value);

		var obj = new Object();
   		obj.type = "unlisten";
   		obj.channel  = c.value;
		obj.id = id.toString(); 
		id++;
   		
		var jsonString= JSON.stringify(obj);
		ws.send(JSON.stringify(obj))
        return false;
    };

	document.getElementById("notify").onclick = function(evt) {
        if (!ws) {
            return false;
        }

		var c = document.getElementById('notifychan');
		if (!c.value) return false;
		var p = document.getElementById('notifypayload');
		if (!p.value) return false;

		print("Notify: " + c.value + ", " + p.value);

		var obj = new Object();
   		obj.type = "notify";
   		obj.channel  = c.value;
		obj.payload  = p.value;
		obj.id = id.toString(); 
		id++;

   		var jsonString= JSON.stringify(obj);
		ws.send(JSON.stringify(obj))
		return false;
    };

	document.getElementById("disconnect_pubsub").onclick = function(evt) {
        if (!wsPubsub) {
            return false;
        }

		wsPubsub.close(1000);
		wsPubsub = null;

		print("Disconnect Pubsub websocket");

        return false;
    };

	document.getElementById("disconnect").onclick = function(evt) {
        if (!ws) {
			return false;
		}
		
		ws.close(1000);
		ws = null;

		print("Disconnect Main websocket");

        return false;
    };
	
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
<form>
	<div>
    	<label for="project">Project</label>
        <input type="text" id="project">
    </div>
	<div>
    	<label for="apikey">ApiKey</label>
        <input type="text" id="apikey">
    </div>
	<button id="connect">Connect</button>
</form>

<form>
    <div>
        <label for="command">Command</label>
        <input type="text" id="command">
	</div>
	<div>
    	<button id="exec">Exec</button>
	</div>
</form>

<form>
    <div>
        <label for="channel">Channel</label>
        <input type="text" id="channel">
	</div>
	<div>
    	<button id="listen">Listen</button><button id="unlisten">Unlisten</button>
	</div>
</form>

<form>
    <div>
        <label for="notifychan">Channel</label>
        <input type="text" id="notifychan">
	</div>
	<div>
		<label for="notifypayload">Payload</label>
        <input type="text" id="notifypayload">
	</div>
	<div>
	    <button id="notify">Notify</button>
	</div>
</form>

<form>
    <button id="disconnect_pubsub">Disconnect Pubsub</button>
</form>

<form>
    <button id="disconnect">Disconnect</button>
</form>

</td><td valign="top" width="50%">
</p>
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr>
</table>
</body>
</html>
`))
