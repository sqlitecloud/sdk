import { logThis } from './utils.js'
//CONFIG
var SQLiteCloud = window.SQLiteCloud;
var config = {
  PROJECT_ID: 'f9cdd1d5-7d16-454b-8cc0-548dc1712c26',
  API_KEY: 'B24tAXTnXFYatN7mSXTPTIRXEEJRiH1lawMEdxmRrps'
};
//DEFINE CALLBACKS FUNCTION PASSED TO WEBSOCKET AND REGISTERED ON ERROR AND CLOSE EVENTS
var onErrorCallback = function (event, msg) {
  var errorCallbackResult = document.getElementById("error-callback-result");
  errorCallbackResult.innerHTML = msg;
  logThis("WebSocket onError callback:" + msg);
  console.log(event);
}
var onCloseCallback = function (msg) {
  var closeCallbackResult = document.getElementById("close-callback-result");
  closeCallbackResult.innerHTML = msg;
  logThis("WebSocket OnClose callback:" + msg);
}
//INIT SQLITECLOUD CLIENT
var client = new SQLiteCloud(config.PROJECT_ID, config.API_KEY, onErrorCallback, onCloseCallback);
//SET REQUEST TIMEOUT
client.setRequestTimeout(5000);
//DECIDED IF FILTER OR NOT SENT MESSAGE
client.setFilterSentMessages(false);
//CONNECTION OPEN
var connectButton = document.getElementById("connect");
var connectResult = document.getElementById("connect-result");
var connect = async function () {
  //try to establish websocket connection
  var connectionResponse = await client.connect();
  //check how websocket connection completed  
  connectResult.innerHTML = connectionResponse.data.message;
  if (connectionResponse.status == 'success') {
    //successful websocket connection
    logThis(connectionResponse.data.message);
  } else {
    //error on websocket connection
    logThis(connectionResponse.data.message);
  }
}
connectButton.addEventListener("click", connect);
//CONNECTION CLOSE
var closeAllButton = document.getElementById("close-all");
var closeOnlyMainButton = document.getElementById("close-only-main");
var closeResult = document.getElementById("close-result");
var close = function (closeAll) {
  //try to close websocket connection
  var closeResponse = client.close(closeAll);
  //check how websocket close completed  
  console.log(closeResponse);
  closeResult.innerHTML = closeResponse.data.message;
  if (closeResponse.status == 'success') {
    //successful websocket close
    logThis(closeResponse.data.message);
  } else {
    //error on websocket close
    logThis(closeResponse.data.message);
  }
}
closeAllButton.addEventListener("click", function () { close(true) });
closeOnlyMainButton.addEventListener("click", function () { close(false) });
//MAIN WEBSOCKET AND PUBSUB WEBSOCKET STATE
var mainWebSocketState = document.getElementById("main-websocket-state");
var pubSubWebSocketState = document.getElementById("pubsub-websocket-state");
setInterval(function () {
  mainWebSocketState.innerHTML = client.connectionState;
  pubSubWebSocketState.innerHTML = client.pubSubState;
}, 500)
//LIST CHANNELS
var listChannelsButton = document.getElementById("list-channels");
var listChannelsResult = document.getElementById("list-channels-result");
var listChannels = async function () {
  //try to request channels
  var listChannelsResponse = await client.listChannels();
  //check how channels request completed  
  if (listChannelsResponse.status == 'success') {
    //successful channels request connection
    listChannelsResult.innerHTML = "";
    logThis("received channels list");
    console.log(listChannelsResponse.data.rows);
    var channels = listChannelsResponse.data.rows;
    for (var i = 0; i < channels.length; i++) {
      logThis("ch. name: " + channels[i].chname);
      var li = document.createElement("li");
      li.innerText = channels[i].chname;
      listChannelsResult.appendChild(li);
    }
  } else {
    //error on channels request
    listChannelsResult.innerHTML = listChannelsResponse.data.message;
    logThis(listChannelsResponse.data.message);
  }
}
listChannelsButton.addEventListener("click", listChannels);
//LIST CHANNELS
var listTablesButton = document.getElementById("list-tables");
var listTablesResult = document.getElementById("list-tables-result");
var listTables = async function () {
  //try to request database tables
  var dbName = "chinook-enc.sqlite";
  var execMessage = `USE DATABASE ${dbName}; LIST TABLES PUBSUB`
  var listTablesResponse = await client.exec(execMessage);
  //check how request database tables completed  
  if (listTablesResponse.status == 'success') {
    //successful database tables completed
    listTablesResult.innerHTML = "";
    logThis("received channels list");
    console.log(listTablesResponse.data.rows);
    var channels = listTablesResponse.data.rows;
    for (var i = 0; i < channels.length; i++) {
      logThis("ch. name: " + channels[i].chname);
      var li = document.createElement("li");
      li.innerText = channels[i].chname;
      listTablesResult.appendChild(li);
    }
  } else {
    //error on database tables completed
    listTablesResult.innerHTML = listTablesResponse.data.message;
    logThis(listTablesResponse.data.message);
  }
}
listTablesButton.addEventListener("click", listTables);
