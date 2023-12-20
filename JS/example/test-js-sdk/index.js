import { logThis } from './utils.js'
//CONFIG
var SQLiteCloud = window.SQLiteCloud;
var config = {
  PROJECT_ID: 'f9cdd1d5-7d16-454b-8cc0-548dc1712c26',
  API_KEY: 'AkNW407R8oKYslzCcdMdCTbAOA8oClRpVYlZLGHZfIs'
};
//DEFINE CALLBACKS FUNCTIONS PASSED TO WEBSOCKET AND REGISTERED ON ERROR AND CLOSE EVENTS
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
  //try to establish main websocket connection
  logThis("start main websocket connection");
  var connectionResponse = await client.connect();
  logThis("end main websocket connection");
  logThis(connectionResponse.status);
  logThis(connectionResponse.data.message);
  connectResult.innerHTML = connectionResponse.data.message;
  if (connectionResponse.status == "success") {
    //do your stuff 
  }
  if (connectionResponse.status == "error") {
    //error handling
  }
}
connectButton.addEventListener("click", connect);
//CONNECTION CLOSE
var closeAllButton = document.getElementById("close-all");
var closeOnlyMainButton = document.getElementById("close-only-main");
var closeResult = document.getElementById("close-result");
var close = function (closeAll) {
  //try to close based on closeAll value 
  //both main websocket and PUB/SUB websocket
  //or only main websocket
  if (closeAll) {
    logThis("start closing main websocket and pub/sub websocket");
  } else {
    logThis("start closing only main websocket");
  }
  var closeResponse = client.close(closeAll);
  if (closeAll) {
    logThis("end closing main websocket and pub/sub websocket");
  } else {
    logThis("end closing only main websocket");
  }
  logThis(closeResponse.status);
  logThis(closeResponse.data.message);
  closeResult.innerHTML = closeResponse.data.message;
  if (closeResponse.status == "success") {
    //do your stuff 
  }
  if (closeResponse.status == "error") {
    //error handling
  }
}
closeAllButton.addEventListener("click", function () { close(true) });
closeOnlyMainButton.addEventListener("click", function () { close(false) });
//MAIN WEBSOCKET AND PUB/SUB WEBSOCKET STATE AND MAIN WEBSOCKET PENDING REQUEST
var mainWebSocketState = document.getElementById("main-websocket-state");
var pubSubWebSocketState = document.getElementById("pubsub-websocket-state");
var mainWebsocketPendingRequests = document.getElementById("main-websocket-pending-requests");
setInterval(function () {
  mainWebSocketState.innerHTML = "state: " + client.connectionState.state + " | " + client.connectionState.description;
  pubSubWebSocketState.innerHTML = "state: " + client.pubSubState.state + " | " + client.pubSubState.description;
  var pendingRequests = client.requestsStackState;
  mainWebsocketPendingRequests.innerText = "";
  for (var requestID of pendingRequests.keys()) {
    var li = document.createElement("li");
    li.innerText = "request ID: " + requestID;
    mainWebsocketPendingRequests.appendChild(li);
  }
}, 500)
//LIST CHANNELS
var listChannelsButton = document.getElementById("list-channels");
var listChannelsResult = document.getElementById("list-channels-result");
var listChannels = async function () {
  //try to request channels
  logThis("start channels list request");
  var listChannelsResponse = await client.listChannels();
  logThis("end channels list request");
  logThis(listChannelsResponse.status);
  if (listChannelsResponse.status == 'success') {
    listChannelsResult.innerHTML = "";
    logThis("received channels list name");
    var channels = listChannelsResponse.data.rows;
    for (var i = 0; i < channels.length; i++) {
      logThis(channels[i].chname);
      var li = document.createElement("li");
      li.innerText = channels[i].chname;
      listChannelsResult.appendChild(li);
    }
  }
  if (listChannelsResponse.status == 'error') {
    logThis(listChannelsResponse.data.message);
    listChannelsResult.innerHTML = listChannelsResponse.data.message;
  }
}
listChannelsButton.addEventListener("click", listChannels);
//CREATE CHANNEL
var createChannelButton = document.getElementById("create-channel");
var createChannelNameInput = document.getElementById("create-channel-name");
var createChannelResult = document.getElementById("create-channel-result");
var createChannel = async function () {
  //try to create a channel
  var newChannelName = createChannelNameInput.value;
  logThis("start creation of channel: " + newChannelName);
  var createChannelsResponse = await client.createChannel(newChannelName);
  logThis("end creation of channel: " + newChannelName);
  logThis(createChannelsResponse.status);
  logThis(createChannelsResponse.data.message);
  createChannelResult.innerHTML = "";
  createChannelResult.innerHTML = createChannelsResponse.data.message;
  if (createChannelsResponse.status == 'success') {
    //do your stuff
  }
  if (createChannelsResponse.status == 'error') {
    //error handling
  }
}
createChannelButton.addEventListener("click", createChannel);
//REMOVE CHANNEL
var removeChannelButton = document.getElementById("remove-channel");
var removeChannelNameInput = document.getElementById("remove-channel-name");
var removeChannelResult = document.getElementById("remove-channel-result");
var removeChannel = async function () {
  //try to remove a channel
  var removeChannelName = removeChannelNameInput.value;
  logThis("start remove of channel: " + removeChannelName);
  var removeChannelsResponse = await client.removeChannel(removeChannelName);
  logThis("end remove of channel: " + removeChannelName);
  logThis(removeChannelsResponse.status);
  logThis(removeChannelsResponse.data.message);
  removeChannelResult.innerHTML = "";
  removeChannelResult.innerHTML = removeChannelsResponse.data.message;
  if (removeChannelsResponse.status == 'success') {
    //do your stuff 
  }
  if (removeChannelsResponse.status == 'error') {
    //error handling
  }
}
removeChannelButton.addEventListener("click", removeChannel);
//EXEC COMMAND
var execCommandButton = document.getElementById("exec-command");
var commandInput = document.getElementById("command");
var execCommandResult = document.getElementById("exec-command-result");
var execCommand = async function () {
  //try to exec command
  logThis("start exec command: " + commandInput.value);
  var execCommandResponse = await client.exec(commandInput.value);
  logThis("end exec command: " + commandInput.value);
  logThis(execCommandResponse.status);
  if (execCommandResponse.status == 'success') {
    execCommandResult.innerHTML = JSON.stringify(execCommandResponse.data);
    logThis(JSON.stringify(execCommandResponse.data));
  }
  if (execCommandResponse.status == 'error') {
    logThis(execCommandResponse.data.message);
    execCommandResult.innerHTML = execCommandResponse.data.message;
  }
}
execCommandButton.addEventListener("click", execCommand);
//NOTIFY MESSAGE
var notifyButton = document.getElementById("notify");
var notifyChannelNameInput = document.getElementById("notify-channel-name");
var notifyMessageInput = document.getElementById("notify-message");
var notifyResult = document.getElementById("notify-result");
var notify = async function () {
  //try to send notification to a channel
  var payload = { message: notifyMessageInput.value };
  logThis("start notify with message: " + JSON.stringify(payload));
  var notificationResponse = await client.notify(notifyChannelNameInput.value, payload);
  logThis("end notify with message: " + JSON.stringify(payload));
  logThis(notificationResponse.status);
  logThis(notificationResponse.data.message);
  notifyResult.innerHTML = "";
  notifyResult.innerHTML = notificationResponse.data.message;
  if (notificationResponse.status == 'success') {
    //do your stuff
  }
  if (notificationResponse.status == 'error') {
    //error handling
  }
}
notifyButton.addEventListener("click", notify);
//LISTEN CHANNEL
var listenChannelButton = document.getElementById("listen-channel");
var listenChannelNameInput = document.getElementById("listen-channel-name");
var listenChannelResult = document.getElementById("listen-channel-result");
var listenChannelIncomingMessage = document.getElementById("listen-channel-incoming-message");
var listenChannelCallback = function (incomingMessage) {
  logThis("incoming message on channel: " + incomingMessage.channel);
  logThis("incoming message payload: " + JSON.stringify(incomingMessage));
  listenChannelResult.innerHTML = "received message on " + incomingMessage.channel;
  listenChannelIncomingMessage.innerHTML = incomingMessage.payload.message;
};
var listenChannel = async function () {
  //try to listen channel
  var listenChannelName = listenChannelNameInput.value;
  logThis("start listen on channel: " + listenChannelName);
  var listenChannelResponse = await client.listenChannel(listenChannelName, listenChannelCallback);
  logThis("end listen on channel: " + listenChannelName);
  logThis(listenChannelResponse.status);
  if (listenChannelResponse.status == 'success') {
    logThis("success on listening to " + listenChannelName);
    listenChannelResult.innerHTML = "";
    listenChannelResult.innerHTML = "listening on " + listenChannelName;
    listenChannelIncomingMessage.innerHTML = "";
  }
  if (listenChannelResponse.status == 'error') {
    logThis(listenChannelResponse.data.message);
    listenChannelResult.innerHTML = listenChannelResponse.data.message;
  }
}
listenChannelButton.addEventListener("click", listenChannel);
//UNLISTEN CHANNEL
var unlistenChannelButton = document.getElementById("unlisten-channel");
var unlistenChannelNameInput = document.getElementById("unlisten-channel-name");
var unlistenChannelResult = document.getElementById("unlisten-channel-result");
var unlistenChannel = async function () {
  //try to unlisten channel
  var unlistenChannelName = unlistenChannelNameInput.value;
  logThis("start unlisten on channel: " + unlistenChannelName);
  var unlistenChannelResponse = await client.unlistenChannel(unlistenChannelName);
  logThis("end unlisten on channel: " + unlistenChannelName);
  logThis(unlistenChannelResponse.status);
  if (unlistenChannelResponse.status == 'success') {
    logThis("success on unlistening to channel: " + unlistenChannelName);
    unlistenChannelResult.innerHTML = "unlistening on " + unlistenChannelResponse.data.channel;
  }
  if (unlistenChannelResponse.status == 'error') {
    logThis(unlistenChannelResponse.data.message);
    unlistenChannelResult.innerHTML = unlistenChannelResponse.data.message;
  }
}
unlistenChannelButton.addEventListener("click", unlistenChannel);
//LISTEN TABLE
var listenTableButton = document.getElementById("listen-table");
var listenTableNameInput = document.getElementById("listen-table-name");
var listenTableResult = document.getElementById("listen-table-result");
var listenTableIncomingMessage = document.getElementById("listen-table-incoming-message");
var listenTableCallback = function (incomingMessage) {
  logThis("incoming message on table: " + incomingMessage.channel);
  logThis("incoming message payload: " + JSON.stringify(incomingMessage));
  listenTableResult.innerHTML = "received message on " + incomingMessage.channel;
  listenTableIncomingMessage.innerHTML = JSON.stringify(incomingMessage.payload);
};
var listenTable = async function () {
  //try to listen table
  var listenTableName = listenTableNameInput.value;
  logThis("start listen on table: " + listenTableName);
  var listenTableResponse = await client.listenTable(listenTableName, listenTableCallback);
  logThis("end listen on table: " + listenTableName);
  logThis(listenTableResponse.status);
  if (listenTableResponse.status == 'success') {
    logThis("success on listening to " + listenTableName);
    listenTableResult.innerHTML = "";
    listenTableResult.innerHTML = "listening on " + listenTableName;
    listenTableIncomingMessage.innerHTML = "";
  }
  if (listenTableResponse.status == 'error') {
    logThis(listenTableResponse.data.message);
    listenTableIncomingMessage.innerHTML = listenTableResponse.data.message;
  }
}
listenTableButton.addEventListener("click", listenTable);
//UNLISTEN TABLE
var unlistenTableButton = document.getElementById("unlisten-table");
var unlistenTableNameInput = document.getElementById("unlisten-table-name");
var unlistenTableResult = document.getElementById("unlisten-table-result");
var unlistenTable = async function () {
  //try to unlisten table
  var unlistenTableName = unlistenTableNameInput.value;
  logThis("start unlisten on table: " + unlistenTableName);
  var unlistenTableResponse = await client.unlistenTable(unlistenTableName);
  logThis("end unlisten on table: " + unlistenTableName);
  logThis(unlistenTableResponse.status);
  if (unlistenTableResponse.status == 'success') {
    logThis("success on unlistening to table " + unlistenTableResponse.data.channel);
    unlistenTableResult.innerHTML = "unlistening on table " + unlistenTableResponse.data.channel;
  }
  if (unlistenTableResponse.status == 'error') {
    logThis(unlistenTableResponse.data.message);
    unlistenTableResult.innerHTML = unlistenTableResponse.data.message;
  }
}
unlistenTableButton.addEventListener("click", unlistenTable);

//ACTUAL SUBSCRIPTIONS STATE
var pubSubSubscriptions = document.getElementById("pubsub-websocket-subscriptions");
setInterval(function () {
  var subscriptionsStack = client.subscriptionsStackState;
  pubSubSubscriptions.innerText = "";
  for (var chName of subscriptionsStack.keys()) {
    var li = document.createElement("li");
    li.innerText = chName;
    pubSubSubscriptions.appendChild(li);
  }
}, 500)