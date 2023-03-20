export const msg = {
  wsConnectOk: "main WebSocket connection is active.",
  wsAlreadyConnected: "main WebSocket connection has already been created.",
  wsCantConnectedWsPubSubExist: "main WebSocket connection cannot be created because a PubSub WebSocket connection is already active.",
  wsConnectError: "main WebSocket connection not established. Check your internet connection, project ID and API key.",
  wsClosingWsPubSubClosingProcess: "closing of the main WebSocket and PubSub WebSocket started.",
  wsClosingProcess: "closing of the main WebSocket started.",
  wsPubSubClosingProcess: "closing of the PubSub WebSocket started.",
  wsCloseComplete: "main WebSocket connection has been closed.",
  wsPubSubClosingProcess: "closing of the PubSub WebSocket started.",
  wsPubSubCloseComplete: "PubSub WebSocket connection has been closed.",
  wsClosingError: "there is no main WebSocket that can be closed.",
  wsCloseByClient: "main main WebSocket connection closed by client.",
  wsPubSubCloseByClient: "PubSub WebSocket main WebSocket connection closed by client.",
  wsNotExist: "main WebSocket connection not exist.",
  wsOnError: "main WebSocket connection error.",
  wsPubSubOnError: "PubSub WebSocket connection error.",
  wsConnecting: "CONNECTING",
  wsOpen: "OPEN",
  wsClosing: "CLOSING",
  wsClosed: "CLOSED",
  wsPubSubNotExist: "PubSub WebSocket connection not exist.",
  wsPubSubConnecting: "CONNECTING",
  wsPubSubOpen: "OPEN",
  wsPubSubClosing: "CLOSING",
  wsPubSubClosed: "CLOSED",
  wsExecErrorNoConnection: "you need to create a main WebSocket connection. Use the connect method.",
  wsNotifyErrorNoConnection: "you need to create a main WebSocket connection. Use the connect method.",
  wsListenError: {
    alreadySubscribed: "registration already made to the channel:",
    errorNoConnection: "you need to create a main WebSocket connection. Use the connect method.",
  },
  wsUnlistenError: {
    missingSubscritption: "it is not possible to unlisten unregistered channel:",
    errorNoConnection: "you have closed the main WebSocket connection, it is no longer possible to unlisten channels.",
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