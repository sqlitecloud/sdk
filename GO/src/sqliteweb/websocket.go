//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/07/04
//    ///             ///   ///  ///    Author      : Andreas Donetti
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
	"fmt"
	"html/template"
	"net/http"
	"sqlitecloud"
	"strings"
	"time"

	"github.com/gobwas/glob"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func init() {
	originChecker := glob.MustCompile("{https://*.sqlitecloud.io,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
	localhostChecker := glob.MustCompile("{https://localhost:*,https://localhost}")

	upgrader.CheckOrigin = func(r *http.Request) bool {
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

func (this *Server) websocketDownload(writer http.ResponseWriter, request *http.Request) {
	var connection *Connection = nil
	var res *sqlitecloud.Result = nil
	var err error = nil

	start := time.Now()

	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromCookie, request)
	v := mux.Vars(request)
	projectID := v["projectID"]

	projectID, _, err = verifyProjectID(int(id), v["projectID"])
	if err != nil {
		SQLiteWeb.Logger.Debug("websocketDownload: unauthorized: ", err)
		return
	}

	// SQLiteWeb.Logger.Debugf("websocketDownload: header %v", request.Header["Cookie"])

	query := fmt.Sprintf("DOWNLOAD DATABASE %s", v["databaseName"])
	connection, err = cm.GetConnection(projectID, false)
	switch {
	case err != nil:
		fallthrough
	case connection == nil:
		fallthrough
	case connection.connection == nil:
		SQLiteWeb.Logger.Debug("websocketDownload: error on getConnection: ", err)
		return
	}

	if res, err = connection.connection.Select(query); err != nil && connection.connection.ErrorCode >= 100000 {
		// internal error (the SDK cannot write to or read from the connection)
		// so remove the current failed connection and retry with a new one
		// for example:
		// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
		// - 100003 Internal Error: sendString (%s)
		cm.closeAndRemoveLockedConnection(projectID, connection)
		SQLiteWeb.Logger.Debug("websocketDownload: Connection Error ", err)
		return
	} else if err != nil || !res.IsArray() {
		// reply must be an Array value (otherwise it is an error)
		cm.ReleaseConnection(projectID, connection)
		SQLiteWeb.Logger.Debug("websocketDownload: error on DOWNLOAD select ", err)
		return
	}

	defer cm.ReleaseConnection(projectID, connection)

	dbSize, _ := res.GetInt64Value(0, 0)
	progressSize := int64(0)

	c, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Debug("websocketDownload: upgrade error: ", err)
		return
	}
	defer c.Close()
	// SQLiteWeb.Logger.Debug("websocketDownload: upgrade")

	for progressSize < dbSize {
		// reply must be a BLOB value (otherwise it is an error)
		if res, err = connection.connection.Select("DOWNLOAD STEP"); err == nil && res.IsBLOB() {
			// res is BLOB, decode it
			buf := res.GetBuffer()
			datalen := len(buf)

			// execute callback (with progressSize updated)
			progressSize = progressSize + int64(datalen)
			c.WriteMessage(websocket.BinaryMessage, buf)

			// check exit condition
			if datalen == 0 {
				break
			}
		} else {
			SQLiteWeb.Logger.Debug("websocketDownload: error while executing download step ", err)
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "error while executing download step"), time.Now().Add(1*time.Second))
			return
		}

		// SQLiteWeb.Logger.Debugf("websocketDownload: loop (progressSize: %d, dbSize: %d)", progressSize, dbSize)
	}

	c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "OK"), time.Now().Add(1*time.Second))

	t := time.Since(start)
	SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s", request.Method, request.URL, request.RemoteAddr, id, t)
}

func (this *Server) websocketUpload(writer http.ResponseWriter, request *http.Request) {
	var connection *Connection = nil
	var res *sqlitecloud.Result = nil
	var err error = nil

	start := time.Now()

	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromCookie, request)
	v := mux.Vars(request)
	projectID := v["projectID"]

	// check authorization for projectID
	projectID, _, err = verifyProjectID(int(id), v["projectID"])
	if err != nil {
		SQLiteWeb.Logger.Debug("websocketUpload: unauthorized: ", err)
		return
	}

	// SQLiteWeb.Logger.Debugf("websocketUpload: header %v", request.Header["Cookie"])

	query := fmt.Sprintf("UPLOAD DATABASE %s", v["databaseName"])
	if keys, ok := request.URL.Query()["key"]; ok && len(keys[0]) > 0 {
		query = fmt.Sprintf("%s key %s", query, keys[0])
	}

	connection, err = cm.GetConnection(projectID, false)
	switch {
	case err != nil:
		fallthrough
	case connection == nil:
		fallthrough
	case connection.connection == nil:
		SQLiteWeb.Logger.Debug("websocketUpload: error on getConnection: ", err)
		return
	}

	if res, err = connection.connection.Select(query); err != nil && connection.connection.ErrorCode >= 100000 {
		// internal error (the SDK cannot write to or read from the connection)
		// so remove the current failed connection and retry with a new one
		// for example:
		// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
		// - 100003 Internal Error: sendString (%s)
		cm.closeAndRemoveLockedConnection(projectID, connection)
		SQLiteWeb.Logger.Debug("websocketUpload: Connection Error ", err)
		return
	} else if err != nil || !res.IsOK() {
		// reply must be an Array value (otherwise it is an error)
		cm.ReleaseConnection(projectID, connection)
		SQLiteWeb.Logger.Debug("websocketUpload: error on UPLOAD select", err)
		return
	}

	defer cm.ReleaseConnection(projectID, connection)

	c, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Debug("websocketUpload: upgrade error: ", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		// SQLiteWeb.Logger.Debugf("websocketUpload: ReadMessage %d %d", mt, len(message))

		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			SQLiteWeb.Logger.Debug("websocketUpload: read error: ", err)
			break
		}

		// send message, and send an 0-length message the received message was a CloseNormalClosure close message
		_, err = connection.connection.SendBytes(message)
		if err != nil {
			SQLiteWeb.Logger.Debug("websocketUpload: SendBytes error : ", err)
			break
		}
		// SQLiteWeb.Logger.Debugf("websocketUpload: SendBytes %d", len(message))

		if len(message) == 0 || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			// end message
			break
		}
	}

	t := time.Since(start)
	SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s", request.Method, request.URL, request.RemoteAddr, id, t)
}

func enquoteString(s string) string {
	enquoted := sqlitecloud.SQCloudEnquoteString(s)
	if strings.HasPrefix(enquoted, "\"") && strings.HasSuffix(enquoted, "\"") {
		enquoted = enquoted[1 : len(enquoted)-1]
	}
	return enquoted
}

func verifyProjectID(userID int, projectUUID string) (string, int, error) {
	query := fmt.Sprintf("SELECT uuid FROM User JOIN Company ON User.company_id = Company.id JOIN Project ON Company.id = Project.company_id WHERE User.enabled=1 AND Company.enabled = 1 AND User.id=%d AND Project.uuid = '%s';", userID, enquoteString(projectUUID))
	res, err, errCode, _ := cm.ExecuteSQL("auth", query)

	if res == nil {
		return "", 503, fmt.Errorf("Service Unavailable")
	}
	if err != nil || errCode != 0 {
		SQLiteWeb.Logger.Debug("verifyProjectID: error ", err)
		return "", 502, fmt.Errorf("Bad Gateway")
	}
	if res.GetNumberOfColumns() != 1 {
		SQLiteWeb.Logger.Debug("verifyProjectID: error on number of columns")
		return "", 502, fmt.Errorf("Bad Gateway")
	}
	if res.GetNumberOfRows() < 1 {
		return "", 404, fmt.Errorf("Project Not Found")
	}
	if res.GetNumberOfRows() > 1 {
		SQLiteWeb.Logger.Debug("verifyProjectID: error on number of rows")
		return "", 502, fmt.Errorf("Bad Gateway")
	}

	return res.GetStringValue_(0, 0), 0, nil
}

func wsTestClient(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, r.Host)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
		console.log(message);
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

	var currentChunk = 0;
	var mime = 'application/octet-binary';
	var finalBlob = null;
	var chunkBlobs =[];

	function addChunk(data) {
		// chunkBlobs[currentChunk] = new Blob([data], {type: mime});

		if (data && data.size > 0) {
			chunkBlobs[currentChunk] = data
			console.log('added chunk ', currentChunk, ' with size ', data.size);
			currentChunk++;
		} else {
			console.log('all chunks completed');

			if (currentChunk > 0) {
				finalBlob = new Blob(chunkBlobs, {type: mime});
				// document.getElementById('completedFileLink').href = URL.createObjectURL(finalBlob);
			
				var a = document.createElement("a"),
            	url = URL.createObjectURL(finalBlob);
        		a.href = url;
        		a.download = "chinook.sqlite";
        		document.body.appendChild(a);
        		a.click();
        		setTimeout(function() {
         			document.body.removeChild(a);
            		window.URL.revokeObjectURL(url);  
        		}, 0); 

				currentChunk = 0
			}
		} 
	}

    document.getElementById("download").onclick = function(evt) {
        if (ws) {
            return false;
        }

		url = "wss://" + "{{.}}" + "/ws/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/bigdb.sqlite/download"
		
		print("DOWNLOAD " + url);

        ws = new WebSocket(url);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
			addChunk(null)
			print("CLOSE (code:" + evt.code + ")");
            ws = null;
        }
        ws.onmessage = function(evt) {
			addChunk(evt.data)
            print("RECEIVED CHUNK: " + evt.data + " size: " + evt.data.size);
        }
        ws.onerror = function(evt) {
            print("WebSocket ERROR: ");
        }
        return false;
    };

	document.getElementById("upload").onclick = function(evt) {
        if (ws) {
            return false;
        }

		print("upload ... ");

		var f = document.getElementById('datafile');
		var k = document.getElementById('enckey');
		
		if (!f.value) return false;
		print("upload value " + f.value);

		var name = f.value.split(/(\\|\/)/g).pop();
		var file = f.files[0];
		var size = file.size;
		var key = (k.value && k.value.length) ? k.value : null;
		var sliceSize = 1024*1024;

		url = "wss://" + "{{.}}" + "/ws/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/" + name + "/upload"
		if (key != null) {
			url = url + "?key=" + key
		}

		print("UPLOAD " + url);

        ws = new WebSocket(url);
		ws.binaryType = "arraybuffer";

        ws.onopen = function(evt) {
            print("OPEN");
			uploadLoop(file, 0, size, sliceSize)
        }
        ws.onclose = function(evt) {
			addChunk(null)
			print("CLOSE (code:" + evt.code + ")");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RECEIVED MESSAGE: " + evt.data + " size: " + evt.data.size);
        }
        ws.onerror = function(evt) {
            print("WebSocket ERROR: ");
        }
		
        return false;
    };


	function slice(file, start, end) {
		var slice = file.mozSlice ? file.mozSlice :
					file.webkitSlice ? file.webkitSlice :
  				  	file.slice ? file.slice : noop;
		return slice.bind(file)(start, end);
	}

	function uploadLoop (file, start, end, size) {
		// compute local values
		var islast = false;
		var len = size;
		if (start + len > end) {len = end - size; islast = true;}
		if (len < 0) len = end;
	
		// compute chunk to send
		var chunk = slice(file, start, start + len);
	
		// chunk is now a Blob that can only be read async
		const reader = new FileReader();
		reader.onloadend = function () {
			print("onloadend length " + reader.result.byteLength)
			ws.send(reader.result);
			
			var value = Math.floor(( start / end) * 100);
			// progressSet(value);
			(islast) ? uploadEnd() : uploadLoop(file, start + size, end, size);
		}
		reader.readAsArrayBuffer(chunk);
	
		return true;
	}

	function uploadEnd () {
		print("uploadEnd")
		ws.close(1000);
		ws = null;
	}

	function createCookie(name,value,days) {
		if (days) {
			var date = new Date();
			date.setTime(date.getTime()+(days*24*60*60*1000));
			var expires = "; expires="+date.toGMTString();
		}
		else var expires = "";
		document.cookie = name+"="+value+expires+"; secure; samesite=lax; path=/";
	}

	createCookie('sqlite-cloud-token','eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY1Nzg0MTU1OSwibmJmIjoxNjU3ODExNTU5LCJpYXQiOjE2NTc4MTE1NTl9.kZN1V6AorVhA7ZGr4raHmcBux_WvDoB-NM31sTS5jvk',1)

    // document.getElementById("send").onclick = function(evt) {
    //     if (!ws) {
    //         return false;
    //     }
    //     print("SEND: " + input.value);
    //     ws.send(input.value);
    //     return false;
    // };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
<form>
<button id="download">Download</button>
</form>

<form>
    <p>Please select an SQLite database and click "Upload" to continue.</p>   
        <div>
            <label for="enckey">Encryption key (optional)</label>
            <input type="text" id="enckey">
        </div>
        <div>
            <input type="file" id="datafile">
        </div>
        <button id="upload">Upload</button>
    </form>

</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr>
</table>
</body>
</html>
`))
