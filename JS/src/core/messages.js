export const msg = {
  wsConnectOk: "webSocket connection is active.",
  wsAlreadyConnected: "webSocket connection has already been created.",
  wsCantConnectedWsPubSubExist: "webSocket connection cannot be created because a pubSub connection is already active.",
  wsConnectError: "websocket connection not established. Check your internet connection, project ID and API key.",
  wsClosingWsPubSubClosingProcess: "closing of the WebSocket and pubSub started.",
  wsClosingProcess: "closing of the WebSocket started.",
  wsPubSubClosingProcess: "closing of the pubSub started.",
  wsCloseComplete: "webSocket connection has been closed.",
  wsPubSubClosingProcess: "closing of the pubSub started.",
  wsPubSubCloseComplete: "pubSub connection has been closed.",
  wsClosingError: "there is no WebSocket that can be closed.",
  wsCloseByClient: "main WebSocket connection closed by client.",
  wsPubSubCloseByClient: "pubSub WebSocket connection closed by client.",
  wsNotExist: "websocket connection not exist.",
  wsOnError: "websocket connection error.",
  wsPubSubOnError: "pubSub connection error.",
  wsConnecting: "CONNECTING",
  wsOpen: "OPEN",
  wsClosing: "CLOSING",
  wsClosed: "CLOSED",
  wsPubSubNotExist: "pubSub connection not exist.",
  wsPubSubConnecting: "CONNECTING",
  wsPubSubOpen: "OPEN",
  wsPubSubClosing: "CLOSING",
  wsPubSubClosed: "CLOSED",
  wsExecErrorNoConnection: "you need to create a WebSocket connection. Use the connect method.",
  wsNotifyErrorNoConnection: "you need to create a WebSocket connection. Use the connect method.",
  wsListenError: {
    alreadySubscribed: "registration already made to the channel:",
    errorNoConnection: "you need to create a WebSocket connection. Use the connect method.",
  },
  wsUnlistenError: {
    missingSubscritption: "it is not possible to unlisten unregistered channel:",
    errorNoConnection: "you have closed the WebSocket connection, it is no longer possible to unlisten channels.",
  },
  wsTimeoutError: "the request timed out. ID:",
  createChannelErr:{
    mandatory: "channelName is mandatory",
    string: "channelName has to be a string"
  },
  dropChannelErr:{
    mandatory: "channelName is mandatory",
    string: "channelName has to be a string"
  } 
}