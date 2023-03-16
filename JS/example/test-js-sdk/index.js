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
  //check how channel request creation completed  
  if (createChannelsResponse.status == 'success') {
    //successful channel request creation
    createChannelResult.innerHTML = "";
    console.log(createChannelsResponse)
    logThis("creation channel " + ' ' + createChannelsResponse.data);
    createChannelResult.innerHTML = createChannelsResponse.data;
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
  //check how channel request creation completed  
  if (removeChannelsResponse.status == 'success') {
    //successful channel request creation
    removeChannelResult.innerHTML = "";
    console.log(removeChannelsResponse)
    logThis("creation channel " + ' ' + removeChannelsResponse.data);
    removeChannelResult.innerHTML = removeChannelsResponse.data;
  } else {
    //error on channel request creation
    removeChannelResult.innerHTML = removeChannelsResponse.data.message;
    logThis(removeChannelsResponse.data.message);
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
  if (execCommandResponse.status == 'success') {
    //successful channel request creation
    execCommandResult.innerHTML = "OK. Read console to see payload";
    execCommandResult.innerHTML = JSON.stringify(execCommandResponse.data);
    logThis("response to " + commandInput.value);
    console.log(execCommandResponse)
  } else {
    //error on channel request creation
    execCommandResult.innerHTML = execCommandResponse.data.message;
    logThis(execCommandResponse.data.message);
  }
}
execCommandButton.addEventListener("click", execCommand);