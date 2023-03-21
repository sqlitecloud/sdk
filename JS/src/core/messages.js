export const msg = {
  wsConnectOk: "main websocket connection is active.",
  wsAlreadyConnected: "main websocket connection has already been created.",
  wsCantConnectedWsPubSubExist: "main websocket connection cannot be created because a PubSub websocket connection is already active.",
  wsConnectError: "main websocket connection not established. Check your internet connection, project ID and API key.",
  wsClosingWsPubSubClosingProcess: "closing of the main websocket and PubSub websocket started.",
  wsClosingProcess: "closing of the main websocket started.",
  wsPubSubClosingProcess: "closing of the PubSub websocket started.",
  wsCloseComplete: "main websocket connection has been closed.",
  wsPubSubClosingProcess: "closing of the PubSub websocket started.",
  wsPubSubCloseComplete: "PubSub websocket connection has been closed.",
  wsClosingError: "there is no main websocket that can be closed.",
  wsCloseByClient: "main main websocket connection closed by client.",
  wsPubSubCloseByClient: "PubSub websocket main websocket connection closed by client.",
  wsNotExist: "main websocket connection not exist.",
  wsOnError: "main websocket connection error.",
  wsPubSubOnError: "PubSub websocket connection error.",
  wsConnecting: "CONNECTING",
  wsOpen: "OPEN",
  wsClosing: "CLOSING",
  wsClosed: "CLOSED",
  wsPubSubNotExist: "PubSub websocket connection not exist.",
  wsPubSubConnecting: "CONNECTING",
  wsPubSubOpen: "OPEN",
  wsPubSubClosing: "CLOSING",
  wsPubSubClosed: "CLOSED",
  wsExecErrorNoConnection: "you need to create a main websocket connection. Use the connect method.",
  wsNotifyErrorNoConnection: "you need to create a main websocket connection. Use the connect method.",
  wsListenError: {
    alreadySubscribed: "registration already made to the channel:",
    errorNoConnection: "you need to create a main websocket connection. Use the connect method.",
  },
  wsUnlistenError: {
    missingSubscritption: "it is not possible to unlisten unregistered channel:",
    errorNoConnection: "you have closed the main websocket connection, it is no longer possible to unlisten channels.",
  },
  wsTimeoutError: "the request timed out. ID:",
  createChannelErr:{
    mandatory: "channelName is mandatory",
    string: "channelName has to be a string"
  },
  removeChannelErr:{
    mandatory: "channelName is mandatory",
    string: "channelName has to be a string"
  } 
}