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
	sqlitecloud "github.com/sqlitecloud/sdk"

	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gobwas/glob"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var dashboardWebsocketUpgrader = websocket.Upgrader{} // use default options

func init() {
	originChecker := glob.MustCompile("{https://*.sqlitecloud.io,https://*.sqlitecloud.io:*,https://sqlitecloud.io,https://sqlitecloud.io:*,https://fabulous-arithmetic-c4a014.netlify.app}")
	localhostChecker := glob.MustCompile("{https://localhost:*,https://localhost}")

	dashboardWebsocketUpgrader.CheckOrigin = func(r *http.Request) bool {
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

func (this *Server) dashboardWebsocketDownload(writer http.ResponseWriter, request *http.Request) {
	var connection *Connection = nil
	var res *sqlitecloud.Result = nil
	var err error = nil

	start := time.Now()

	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromCookie, request)
	v := mux.Vars(request)
	projectID := v["projectID"]

	projectID, _, err = verifyProjectID(id, v["projectID"], dashboardcm)
	if err != nil {
		SQLiteWeb.Logger.Error("dashboardWebsocketDownload: unauthorized: ", err)
		return
	}

	// SQLiteWeb.Logger.Debugf("dashboardWebsocketDownload: header %v", request.Header["Cookie"])

	c, err := dashboardWebsocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Error("dashboardWebsocketDownload: upgrade error: ", err)
		return
	}
	defer c.Close()
	// SQLiteWeb.Logger.Debug("dashboardWebsocketDownload: upgrade")

	query := "DOWNLOAD DATABASE ?" // , enquoteString(v["databaseName"]))
	connection, err = dashboardcm.GetConnection(projectID, false)
	switch {
	case err != nil:
		fallthrough
	case connection == nil:
		fallthrough
	case connection.connection == nil:
		SQLiteWeb.Logger.Error("dashboardWebsocketDownload: error on getConnection")
		return
	}

	if res, err = connection.connection.SelectArray(query, []interface{}{v["databaseName"]}); err != nil && connection.connection.ErrorCode >= 100000 {
		// internal error (the SDK cannot write to or read from the connection)
		// so remove the current failed connection and retry with a new one
		// for example:
		// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
		// - 100003 Internal Error: sendString (%s)
		dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
		SQLiteWeb.Logger.Debug("dashboardWebsocketDownload: Connection Error ", err)
		c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
		return
	} else if err != nil || !res.IsArray() {
		// reply must be an Array value (otherwise it is an error)

		// try to abort the current download operation
		if err1 := connection.connection.Execute("DOWNLOAD ABORT"); err1 != nil {
			SQLiteWeb.Logger.Errorf("dashboardWebsocketDownload: DOWNLOAD ABORT error (%s) closing conn: %v", err1, connection)
			dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
		} else {
			dashboardcm.ReleaseConnection(projectID, connection)
		}

		// prepare the error message
		closemsg := ""
		switch {
		case err != nil:
			closemsg = err.Error()
		case !res.IsArray():
			closemsg = fmt.Sprintf("expected array, got type %c", res.GetType())
		}
		SQLiteWeb.Logger.Errorf("dashboardWebsocketDownload: error on DOWNLOAD (%s): %s", query, closemsg)

		// send the close message
		c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, closemsg), time.Now().Add(1*time.Second))
		return
	}

	defer dashboardcm.ReleaseConnection(projectID, connection)

	dbSize, _ := res.GetInt64Value(0, 0)
	progressSize := int64(0)

	for progressSize < dbSize {
		// reply must be a BLOB value (otherwise it is an error)
		if res, err = connection.connection.Select("DOWNLOAD STEP"); err == nil && res.IsBLOB() {
			// res is BLOB, decode it
			buf := res.GetBuffer()
			datalen := len(buf)

			// execute callback (with progressSize updated)
			progressSize = progressSize + int64(datalen)
			err = c.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				SQLiteWeb.Logger.Error("dashboardWebsocketDownload: error on STEP writeMessage: ", err)
				if err1 := connection.connection.Execute("DOWNLOAD ABORT"); err1 != nil {
					SQLiteWeb.Logger.Errorf("dashboardWebsocketDownload: DOWNLOAD ABORT error (%s) closing conn: %v", err1, connection)
					dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
				}
				c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
				return
			}

			// check exit condition
			if datalen == 0 {
				break
			}
		} else {
			SQLiteWeb.Logger.Error("dashboardWebsocketDownload: error while executing download step ", err)
			if err1 := connection.connection.Execute("DOWNLOAD ABORT"); err1 != nil {
				SQLiteWeb.Logger.Errorf("dashboardWebsocketDownload: DOWNLOAD ABORT error (%s) closing conn: %v", err1, connection)
				dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
			}
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "error while executing download step"), time.Now().Add(1*time.Second))
			return
		}

		// SQLiteWeb.Logger.Debugf("dashboardWebsocketDownload: loop (progressSize: %d, dbSize: %d)", progressSize, dbSize)
	}

	c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "OK"), time.Now().Add(1*time.Second))

	t := time.Since(start)
	SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s iserr:%v", request.Method, request.URL, request.RemoteAddr, id, t, err != nil)
}

func (this *Server) dashboardWebsocketUpload(writer http.ResponseWriter, request *http.Request) {
	var connection *Connection = nil
	var res *sqlitecloud.Result = nil
	var err error = nil

	start := time.Now()

	id, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromCookie, request)
	v := mux.Vars(request)
	projectID := v["projectID"]

	// check authorization for projectID
	projectID, _, err = verifyProjectID(id, v["projectID"], dashboardcm)
	if err != nil {
		SQLiteWeb.Logger.Error("dashboardWebsocketUpload: unauthorized: ", err)
		return
	}

	c, err := dashboardWebsocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		SQLiteWeb.Logger.Error("dashboardWebsocketUpload: upgrade error: ", err)
		return
	}
	defer c.Close()

	// SQLiteWeb.Logger.Debugf("dashboardWebsocketUpload: header %v", request.Header["Cookie"])

	query := "UPLOAD DATABASE ?" // , enquoteString(v["databaseName"]))
	args := []interface{}{v["databaseName"]}
	if keys, ok := request.URL.Query()["key"]; ok && len(keys[0]) > 0 {
		query = fmt.Sprintf("%s key ?", query) // enquoteString(keys[0])
		args = append(args, keys[0])
	}

	connection, err = dashboardcm.GetConnection(projectID, false)
	switch {
	case err != nil:
		fallthrough
	case connection == nil:
		fallthrough
	case connection.connection == nil:
		SQLiteWeb.Logger.Error("dashboardWebsocketUpload: error on getConnection")
		c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Cannot connect to node"), time.Now().Add(1*time.Second))
		return
	}

	if res, err = connection.connection.SelectArray(query, args); err != nil && connection.connection.ErrorCode >= 100000 {
		// internal error (the SDK cannot write to or read from the connection)
		// so remove the current failed connection and retry with a new one
		// for example:
		// - 100001 Internal Error: SQCloud.readNextRawChunk (%s)
		// - 100003 Internal Error: sendString (%s)
		dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
		SQLiteWeb.Logger.Error("dashboardWebsocketUpload: Connection Error ", err)
		c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
		return
	} else if err != nil || !res.IsOK() {
		// reply must be an OK value (otherwise it is an error)

		// try to abort the current download operation
		if err1 := connection.connection.Execute("UPLOAD ABORT"); err1 != nil {
			SQLiteWeb.Logger.Errorf("dashboardWebsocketUpload: UPLOAD ABORT error (%s) closing conn: %v", err1, connection)
			dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
		} else {
			// SQLiteWeb.Logger.Debugf("dashboardWebsocketUpload: ReleaseConnection %v", connection)
			dashboardcm.ReleaseConnection(projectID, connection)
		}

		// prepare the error message
		closemsg := ""
		switch {
		case err != nil:
			closemsg = err.Error()
		case !res.IsOK():
			closemsg = fmt.Sprintf("expected OK, got type %c", res.GetType())
		}
		SQLiteWeb.Logger.Errorf("dashboardWebsocketUpload: error on UPLOAD (%s): %s", query, closemsg)

		// send the close message
		c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, closemsg), time.Now().Add(1*time.Second))
		return
	}

	defer dashboardcm.ReleaseConnection(projectID, connection)

	// temporarily increase the timeout, otherwise the SendBlob function would probably
	// give an "SQCloud.readNextRawChunk (Timeout)" error, expecially while waiting for
	// the response for the last empty message sent with SendBlob. After the last message,
	// the leader must transfer the database to every node, then send the load database
	// message and eventually reply with "OK", and all this process can last for minutes.
	originaltimeout := connection.connection.Timeout
	connection.connection.Timeout = time.Duration(1) * time.Hour
	defer func(connection *Connection, timeout time.Duration) {
		if connection != nil && connection.connection != nil {
			connection.connection.Timeout = originaltimeout
		}
	}(connection, originaltimeout)

	for {
		_, message, err := c.ReadMessage()
		SQLiteWeb.Logger.Debugf("dashboardWebsocketUpload: ReadMessage %d", len(message))

		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			SQLiteWeb.Logger.Debug("dashboardWebsocketUpload: read error: ", err)
			if err1 := connection.connection.Execute("UPLOAD ABORT"); err1 != nil {
				SQLiteWeb.Logger.Errorf("dashboardWebsocketUpload: UPLOAD ABORT error (%s) closing conn: %v", err1, connection)
				dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
			}
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			break
		}

		// send message, and send an 0-length message if the received message was a CloseNormalClosure close message
		err = connection.connection.SendBlob(message)
		if err != nil {
			SQLiteWeb.Logger.Error("dashboardWebsocketUpload: SendBytes error: ", err)
			if err1 := connection.connection.Execute("UPLOAD ABORT"); err1 != nil {
				SQLiteWeb.Logger.Errorf("dashboardWebsocketUpload: UPLOAD ABORT error (%s) closing conn: %v", err1, connection)
				dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
			}
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			break
		}
		SQLiteWeb.Logger.Debugf("dashboardWebsocketUpload: SendBytes %d completed", len(message))

		if len(message) == 0 {
			// empty message: end message
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "OK"), time.Now().Add(1*time.Second))
			break
		} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			// closed channel: end message
			break
		} else {
			// send back the ack message
			if err = c.WriteMessage(websocket.TextMessage, []byte("OK")); err != nil {
				if err1 := connection.connection.Execute("UPLOAD ABORT"); err1 != nil {
					SQLiteWeb.Logger.Errorf("dashboardWebsocketUpload: UPLOAD ABORT error (%s) closing conn: %v", err1, connection)
					dashboardcm.closeAndRemoveLockedConnection(projectID, connection)
				}
				c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()), time.Now().Add(1*time.Second))
			}
			SQLiteWeb.Logger.Debug("dashboardWebsocketUpload: ack")
		}
	}

	t := time.Since(start)
	SQLiteWeb.Logger.Debugf("Endpoint \"%s %s\" addr:%s user:%d exec_time:%s iserr:%v", request.Method, request.URL, request.RemoteAddr, id, t, err != nil)
}

func enquoteString(s string) string {
	enquoted := sqlitecloud.SQCloudEnquoteString(s)
	if strings.HasPrefix(enquoted, "\"") && strings.HasSuffix(enquoted, "\"") {
		enquoted = enquoted[1 : len(enquoted)-1]
	}
	return enquoted
}

func dwsTestClient(w http.ResponseWriter, r *http.Request) {
	dwsTestClientTemplate.Execute(w, r.Host)
}

var dwsTestClientTemplate = template.Must(template.New("").Parse(`
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
	var chunkBlobs = [];
	var islast = false;
	var nextstart = 0;

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

		url = "wss://" + "{{.}}" + "/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/wrongdb.sqlite/download"
		
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
		var totalsize = file.size;
		var key = (k.value && k.value.length) ? k.value : null;
		var sliceSize = 1024*1024;

		url = "wss://" + "{{.}}" + "/dashboard/v1/f9cdd1d5-7d16-454b-8cc0-548dc1712c26/database/" + name + "/upload"
		if (key != null) {
			url = url + "?key=" + key
		}

		print("UPLOAD " + url);

        ws = new WebSocket(url);
		ws.binaryType = "arraybuffer";

        ws.onopen = function(evt) {
            print("OPEN");
			islast = false;
			nextstart = 0;
			uploadLoop(file, nextstart, totalsize, sliceSize)
        }
        ws.onclose = function(evt) {
			addChunk(null)
			print("CLOSE: reason:" + evt.reason + ", code:" + evt.code);
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RECEIVED MESSAGE: " + evt.data + " size: " + evt.data.size);
			if (evt.data === "OK") {(islast) ? uploadEnd() : uploadLoop(file, nextstart, totalsize, sliceSize);}
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
			print("sending bytes " + reader.result.byteLength)
			ws.send(reader.result);
			print("bytes sent")

			// var value = Math.floor(( start / end) * 100);
			// progressSet(value);

			nextstart = start + size
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

	createCookie('sqlite-cloud-token','eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiQW5kcmVhIiwibGFzdF9uYW1lIjoiRG9uZXR0aSIsImlwYSI6IjEyNy4wLjAuMSIsImlzcyI6IndlYi5zcWxpdGVjbG91ZC5pbyIsInN1YiI6IjIiLCJhdWQiOlsid2ViLnNxbGl0ZWNsb3VkLmlvIl0sImV4cCI6MTY1ODg1NjY5MCwibmJmIjoxNjU4ODI2NjkwLCJpYXQiOjE2NTg4MjY2OTB9.0bxqtFHuEZuz1wfN8NVX3kB7uffyPEM9xc_Zg17DOx0',1)

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
