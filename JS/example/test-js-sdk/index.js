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
  if (connectionResponse.status == 'success' || connectionResponse.status == 'warning') {
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
  if (closeResponse.status == 'success' || closeResponse.status == 'warning') {
    //successful websocket close
    logThis(closeResponse.data.message);
  } else {
    //error on websocket close
    logThis(closeResponse.data.message);
  }
}
closeAllButton.addEventListener("click", function () { close(true) });
closeOnlyMainButton.addEventListener("click", function () { close(false) });
//MAIN WEBSOCKET AND PUBSUB WEBSOCKET STATE AND MAIN WEBSOCKET PENDING REQUEST
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
  var listChannelsResponse = await client.listChannels();
  //check how channels request completed  
  if (listChannelsResponse.status == 'success' || listChannelsResponse.status == 'warning') {
    //successful channels request connection
    listChannelsResult.innerHTML = "";
    logThis("received channels list");
    console.log(listChannelsResponse);
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
//CREATE CHANNEL
var createChannelButton = document.getElementById("create-channel");
var createChannelNameInput = document.getElementById("create-channel-name");
var createChannelResult = document.getElementById("create-channel-result");
var createChannel = async function () {
  //try to create a channel
  var createChannelsResponse = await client.createChannel(createChannelNameInput.value);
  console.log(createChannelsResponse) //TOGLI
  //check how channel request creation completed  
  if (createChannelsResponse.status == 'success' || createChannelsResponse.status == 'warning') {
    //successful channel request creation
    createChannelResult.innerHTML = "";
    console.log(createChannelsResponse)
    logThis("creation channel " + ' ' + createChannelsResponse.data.message);
    createChannelResult.innerHTML = createChannelsResponse.data.message;
  } else {
    //error on channel request creation
    createChannelResult.innerHTML = createChannelsResponse.data.message;
    logThis(createChannelsResponse.data.message);
  }
}
createChannelButton.addEventListener("click", createChannel);
//REMOVE CHANNEL
var removeChannelButton = document.getElementById("remove-channel");
var removeChannelNameInput = document.getElementById("remove-channel-name");
var removeChannelResult = document.getElementById("remove-channel-result");
var removeChannel = async function () {
  //try to remove a channel
  var removeChannelsResponse = await client.removeChannel(removeChannelNameInput.value);
  console.log(removeChannelsResponse); //TOGLI
  //check how channel request creation completed  
  if (removeChannelsResponse.status == 'success' || removeChannelsResponse.status == 'warning') {
    //successful channel request creation
    removeChannelResult.innerHTML = "";
    console.log(removeChannelsResponse)
    logThis("creation channel " + ' ' + removeChannelsResponse.data.message);
    removeChannelResult.innerHTML = removeChannelsResponse.data.message;
  } else {
    //error on channel request creation
    removeChannelResult.innerHTML = removeChannelsResponse.data.message;
    logThis(removeChannelsResponse);
  }
}
removeChannelButton.addEventListener("click", removeChannel);
//EXEC COMMAND
var execCommandButton = document.getElementById("exec-command");
var commandInput = document.getElementById("command");
var execCommandResult = document.getElementById("exec-command-result");
var execCommand = async function () {
  //try to exec command
  var execCommandResponse = await client.exec(commandInput.value);
  console.log(execCommandResponse);
  //check how command execution request completed  
  if (execCommandResponse.status == 'success' || execCommandResponse.status == 'warning') {
    //successful channel request creation
    execCommandResult.innerHTML = "OK. Read console to see payload";
    execCommandResult.innerHTML = JSON.stringify(execCommandResponse.data);
    logThis("response to " + commandInput.value);
  } else {
    //error on channel request creation
    execCommandResult.innerHTML = execCommandResponse.data.message;
    logThis(execCommandResponse.data.message);
  }
}
execCommandButton.addEventListener("click", execCommand);
//CREATE CHANNEL
var notifyButton = document.getElementById("notify");
var notifyChannelNameInput = document.getElementById("notify-channel-name");
var notifyMessageInput = document.getElementById("notify-message");
var notifyResult = document.getElementById("notify-result");
var notify = async function () {
  var payload = { message: notifyMessageInput.value };
  //try to send notification to a channel
  var notificationResponse = await client.notify(notifyChannelNameInput.value, payload);
  console.log(notificationResponse)
  //check how notifcation completed  
  if (notificationResponse.status == 'success' || notificationResponse.status == 'warning') {
    //successful notification
    notifyResult.innerHTML = "";
    logThis("success on notify to channel " + notifyChannelNameInput.value);
    notifyResult.innerHTML = notificationResponse.data.message;
  } else {
    //error on notification
    if (notificationResponse.data.message) {
      notifyResult.innerHTML = notificationResponse.data.message;
    } else {
      notifyResult.innerHTML = "ERROR";
    }
    logThis(notificationResponse.data);
  }
}
notifyButton.addEventListener("click", notify);
//LISTEN CHANNEL
var listenChannelButton = document.getElementById("listen-channel");
var listenChannelNameInput = document.getElementById("listen-channel-name");
var listenChannelResult = document.getElementById("listen-channel-result");
var listenChannelIncomingMessage = document.getElementById("listen-channel-incoming-message");
var listenChannelCallback = function (incomingMessage) {
  console.log(incomingMessage) //TOGLI
  logThis("incoming message on channel " + listenChannelNameInput.value);
  console.log(incomingMessage);
  listenChannelResult.innerHTML = "received message on " + incomingMessage.channel;
  listenChannelIncomingMessage.innerHTML = incomingMessage.payload.message;
};
var listenChannel = async function () {
  //try to listen channel
  var listenChannelResponse = await client.listenChannel(listenChannelNameInput.value, listenChannelCallback);
  console.log(listenChannelResponse)
  //check how listen completed  
  if (listenChannelResponse.status == 'success' || listenChannelResponse.status == 'warning') {
    //successful listening
    listenChannelResult.innerHTML = "";
    logThis("success on listening to " + listenChannelNameInput.value,);
    listenChannelResult.innerHTML = "listening on channel " + listenChannelNameInput.value;
    listenChannelIncomingMessage.innerHTML = "";
  } else {
    //error on listetning
    if (listenChannelResponse.data.message) {
      listenChannelResult.innerHTML = listenChannelResponse.data.message;
      logThis(listenChannelResponse.data.message);
    } else {
      listenChannelResult.innerHTML = "ERROR";
    }
  }
}
listenChannelButton.addEventListener("click", listenChannel);
//UNLISTEN CHANNEL
var unlistenChannelButton = document.getElementById("unlisten-channel");
var unlistenChannelNameInput = document.getElementById("unlisten-channel-name");
var unlistenChannelResult = document.getElementById("unlisten-channel-result");
var unlistenChannel = async function () {
  //try to unlisten channel
  var unlistenChannelResponse = await client.unlistenChannel(unlistenChannelNameInput.value);
  //check how unlisten completed  
  if (unlistenChannelResponse.status == 'success' || unlistenChannelResponse.status == 'warning') {
    //unsuccessful listening
    logThis("success on unlistening to channel " + unlistenChannelNameInput.value,);
    console.log(unlistenChannelResponse)
    listenChannelResult.innerHTML = "unlistening on channel " + listenChannelNameInput.value;
  } else {
    //error on listetning
    if (unlistenChannelResponse.data.message) {
      unlistenChannelResult.innerHTML = unlistenChannelResponse.data.message;
    } else {
      unlistenChannelResult.innerHTML = "ERROR";
    }
    logThis(unlistenChannelResponse.data.message);
  }
}
unlistenChannelButton.addEventListener("click", unlistenChannel);
//LISTEN TABLE
var listenTableButton = document.getElementById("listen-table");
var listenTableNameInput = document.getElementById("listen-table-name");
var listenTableResult = document.getElementById("listen-table-result");
var listenTableIncomingMessage = document.getElementById("listen-table-incoming-message");
var listenTableCallback = function (incomingMessage) {
  logThis("incoming message on table " + listenTableNameInput.value);
  console.log(incomingMessage);
  listenTableResult.innerHTML = "received message on " + incomingMessage.table;
  listenTableIncomingMessage.innerHTML = incomingMessage.payload.message;
};
var listenTable = async function () {
  //try to listen table
  var listenTableResponse = await client.listenTable(listenTableNameInput.value, listenTableCallback);
  console.log(listenTableResponse)
  //check how listen completed  
  if (listenTableResponse.status == 'success' || listenTableResponse.status == 'warning') {
    //successful listening
    listenTableResult.innerHTML = "";
    logThis("success on listening to " + listenTableNameInput.value,);
    listenTableResult.innerHTML = "listening on table " + listenTableNameInput.value;
    listenTableIncomingMessage.innerHTML = "";
  } else {
    //error on listetning
    if (listenTableResponse.data.message) {
      listenTableResult.innerHTML = listenTableResponse.data.message;
    } else {
      listenTableResult.innerHTML = "ERROR";
    }
    logThis(listenTableResponse.data.message);
  }
}
listenTableButton.addEventListener("click", listenTable);
//UNLISTEN TABLE
var unlistenTableButton = document.getElementById("unlisten-table");
var unlistenTableNameInput = document.getElementById("unlisten-table-name");
var unlistenTableResult = document.getElementById("unlisten-table-result");
var unlistenTable = async function () {
  //try to unlisten table
  var unlistenTableResponse = await client.unlistenTable(unlistenTableNameInput.value);
  //check how unlisten completed  
  if (unlistenTableResponse.status == 'success' || unlistenTableResponse.status == 'warning') {
    //unsuccessful listening
    logThis("success on unlistening to table " + unlistenTableNameInput.value,);
    console.log(unlistenTableResponse)
    listenTableResult.innerHTML = "unlistening on table " + listenTableNameInput.value;
  } else {
    //error on listetning
    if (unlistenTableResponse.data.message) {
      unlistenTableResult.innerHTML = unlistenTableResponse.data.message;
    } else {
      unlistenTableResult.innerHTML = "ERROR";
    }
    logThis(unlistenTableResponse.data.message);
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