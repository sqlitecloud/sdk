import { msg } from "./messages";



export default class SQLiteCloud {
  /* PRIVATE PROPERTIES */

  /*
  #ws private property stores the websocket used:
    - to send "exec" type request
    - to send "pubsub" subscription request
    - to receive response to "exec" type request
  
  User receives the responses to his requests reading the result of a Promise.
  */
  #ws = null;

  /*
  #wsPubSub private property stores the websocket used to receive pubSub messages.
  #wsPubSubUrl private property stores the websocket url
  #uuid private property stores the user identifier. This can be used to not received messages sent by the current user
  When a new message is received, based on the channel, is selected the callbacks to be called cycling through subscriptionsStack property
  */
  #wsPubSub = null;
  #wsPubSubUrl = null;
  #uuid = null;

  /*
  #requestsStack private property stores the list of pending requests, in this way the user can send multiple parallel requests.
  For each request an object that contains the following is stored:
  {
    id: //unique id associated with the request 
    onRequestTimeout: //function called when the request times out
    resolve: resolve, //function called when the Promise resolve
    reject: reject //function called when the Promise reject
  }
  */
  #requestsStack = new Map();

  /*
  #subscriptionsStack private property stores the list of pubSub subscriptions.
  For each subscription an object that contains the following is stored:
  {
    channel: //the name of the channel you are subscribed to 
    callback: //the function called when a new message arrives
  }
  */
  #subscriptionsStack = new Map();

  /* PUBLIC PROPERTIES */
  /*
  requestTimeout property stores the time available to receive a response before the request times out.
  filterSentMessages the library not return messages sent by the user
  */
  requestTimeout = 3000;
  filterSentMessages = false;


  /* CONSTRUCTOR */
  /*
  SQLiteCloud class constructor receives:
   - project ID (required)
   - api key (required)
   - webSocket callback event (optional)
  */
  constructor(projectID, apikey, onError = null, onClose = null) {
    this.url = `wss://web1.sqlitecloud.io:8443/api/v1/${projectID}/ws?apikey=${apikey}`;
    this.onError = onError;
    this.onClose = onClose;
  }

  /* PUBLIC METHODS */
  /*
  setRequestTimeout method allows the user to change the request timeout value
  */
  setRequestTimeout(value) {
    this.requestTimeout = value;
  }

  /*
  setFilterSentMessages method allows the user to filter or not sent messagess
  */
  setFilterSentMessages(value) {
    this.filterSentMessages = value;
  }


  /*
  connect method opens websocket connection
  */
  async connect() {
    if (this.#ws == null) {
      if (this.#wsPubSub == null) {
        try {
          this.#ws = await this.#connectWs(this.url, msg.wsConnectError);
          //register the error event on websocket
          this.#ws.addEventListener('error', this.#onErrorWs);
          //register the close event on websocket
          this.#ws.addEventListener('close', this.#onCloseWs);
          return {
            status: "success",
            data: {
              message: msg.wsConnectOk
            }
          }
        } catch (error) {
          return {
            status: "error",
            data: {
              message: error.toString(),
              error: error
            }
          }
        }
      } else {
        return {
          status: "error",
          data: {
            message: msg.wsCantConnectedWsPubSubExist
          }
        };
      }
    } else {
      return {
        status: "error",
        data: {
          message: msg.wsAlreadyConnected
        }
      }
    }
  }

  /*
  close method closes websocket connection and if exist pubSub websocket connection
  - closePubSub = true (default true), this method close both WebSocket
  - closePubSub = false this method close only main WebSocket
  */
  close(closePubSub = true) {
    if (closePubSub) {
      if (this.#wsPubSub !== null && this.#ws !== null) {
        this.#wsPubSub.close(1000, msg.wsPubSubCloseByClient);
        this.#subscriptionsStack = new Map();
        this.#ws.close(1000, msg.wsCloseByClient);
        return (
          {
            status: "success",
            data: {
              message: msg.wsClosingWsPubSubClosingProcess
            }
          }
        )
      } else if (this.#wsPubSub == null && this.#ws !== null) {
        this.#ws.close(1000, msg.wsCloseByClient);
        return (
          {
            status: "success",
            data: {
              message: msg.wsClosingProcess
            }
          }
        )
      } else if (this.#wsPubSub != null && this.#ws == null) {
        this.#subscriptionsStack = new Map();
        this.#wsPubSub.close(1000, msg.wsCloseByClient);
        return (
          {
            status: "success",
            data: {
              message: msg.wsPubSubClosingProcess
            }
          }
        )
      } else {
        return (
          {
            status: "error",
            data: {
              message: msg.wsClosingError
            }
          }
        )
      }
    }
    if (!closePubSub) {
      if (this.#ws !== null) {
        this.#ws.close(1000, msg.wsCloseByClient);
        return (
          {
            status: "success",
            data: {
              message: msg.wsClosingProcess
            }
          }
        )
      } else {
        return (
          {
            status: "error",
            data: {
              message: msg.wsClosingError
            }
          }
        )
      }
    }
  }

  /*
  connectionState method returns the actual state of websocket connection
  */
  get connectionState() {
    let connectionStateString = msg.wsNotExist;
    let connectionState = -1;
    if (this.#ws !== null) {
      connectionState = this.#ws.readyState;
      switch (connectionState) {
        case 0:
          connectionStateString = msg.wsConnecting;
          break;
        case 1:
          connectionStateString = msg.wsOpen;
          break;
        case 2:
          connectionStateString = msg.wsClosing;
          break;
        case 3:
          connectionStateString = msg.wsClosed;
          break;
        default:
          connectionStateString = msg.wsNotExist;
      }
    }
    return {
      state: connectionState,
      description: connectionStateString
    };
  }

  /*
  pubSubState method returns the actual state of pubSubState websocket connection
  */
  get pubSubState() {
    let connectionStateString = msg.wsPubSubNotExist;
    let connectionState = -1;
    if (this.#wsPubSub !== null) {
      connectionState = this.#wsPubSub.readyState;
      switch (connectionState) {
        case 0:
          connectionStateString = msg.wsPubSubConnecting;
          break;
        case 1:
          connectionStateString = msg.wsPubSubOpen;
          break;
        case 2:
          connectionStateString = msg.wsPubSubClosing;
          break;
        case 3:
          connectionStateString = msg.wsPubSubClosed;
          break;
        default:
          connectionStateString = msg.wsPubSubNotExist;
      }
    }
    return {
      state: connectionState,
      description: connectionStateString
    };
  }

  /*
  requestsStackState method return the lits of pending requests
  */
  get requestsStackState() {
    return this.#requestsStack;
  }

  /*
  subscriptionsStackState method return the lits of active subscriptions
  */
  get subscriptionsStackState() {
    return this.#subscriptionsStack;
  }


  /*
  exec method calls the private method promiseSend building the object 
  */
  async exec(command) {
    if (this.#ws !== null) {
      try {
        const response = await this.#promiseSend(
          {
            id: this.#makeid(5),
            type: "exec",
            command: command
          }
        );
        return response;
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return {
        status: "error",
        data: {
          message: msg.wsExecErrorNoConnection
        }
      }
    }
  }


  /*
  notify method calls the private method promiseSend building the object 
  */
  async notify(channel, payload) {
    if (this.#ws !== null) {
      try {
        const response = await this.#promiseSend(
          {
            id: this.#makeid(5),
            type: "notify",
            channel: channel,
            payload: JSON.stringify(payload)
          }
        );
        if (response.status == "success") {
          return (
            {
              status: response.status,
              data: {
                message: "OK"
              }
            }
          )
        }
        if (response.status == "error") {
          return response
        }
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return (
        {
          status: "error",
          data: {
            message: msg.wsNotifyErrorNoConnection
          }
        }
      );
    }
  }


  /*
  listenChannel method calls the private method #pubsub to register to a new channel 
  */
  async listenChannel(channel, callback) {
    if (this.#ws !== null) {
      try {
        const response = await this.#pubsub("listenChannel", channel, callback);
        return response;
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return ({
        status: "error",
        data: {
          message: msg.wsListenError.errorNoConnection
        }
      })
    }
  }


  /*
  listenTable method calls the private method #pubsub to register to a new table 
  */
  async listenTable(table, callback) {
    if (this.#ws !== null) {
      try {
        const response = await this.#pubsub("listenTable", table, callback);
        return response;
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return ({
        status: "error",
        data: {
          message: msg.wsListenError.errorNoConnection
        }
      })
    }
  }


  /*
  unlistenChannel method calls the private method #pubsub to unregister to a channel 
  */
  async unlistenChannel(channel) {
    if (this.#ws !== null) {
      try {
        if (!this.#subscriptionsStack.has(channel)) {
          return (
            {
              status: "error",
              data: {
                message: msg.wsUnlistenError.missingSubscritption + " " + channel
              }
            }
          )
        }
        console.log("STO PER unlinset il canale " + channel)
        const response = await this.#pubsub("unlistenChannel", channel, null);
        return response;
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return (
        {
          status: "error",
          data: {
            message: msg.wsUnlistenError.errorNoConnection
          }
        }
      )
    }
  }


  /*
  unlistenTable method calls the private method #pubsub to unregister to a table 
  */
  async unlistenTable(table) {
    if (this.#ws !== null) {
      try {
        if (!this.#subscriptionsStack.has(table)) {
          return (
            {
              status: "error",
              data: {
                message: msg.wsUnlistenError.missingSubscritption + " " + table
              }
            }
          )
        }
        const response = await this.#pubsub("unlistenTable", table, null);
        return response;
      } catch (error) {
        return {
          status: "error",
          data: {
            message: error.toString(),
            error: error
          }
        }
      }
    } else {
      return (
        {
          status: "error",
          data: {
            message: msg.wsUnlistenError.errorNoConnection
          }
        }
      )
    }
  }

  /*
  listChannels method calls exec method to receive the list of all the existing channels 
  */
  async listChannels() {
    try {
      const response = await this.exec("LIST CHANNELS");
      return (response);
    } catch (error) {
      return {
        status: "error",
        data: error
      };
    }
  }

  /*
  createChannel method calls exec method to create a new channel
  channelName: mandatory, the name of the channel to be created
  ifNotExist: optional, if true set in the command the string [IF NOT EXISTS]
  */
  async createChannel(channelName, ifNotExist = true) {
    try {
      //params validation
      //check channelName has been provided
      if (!channelName) {
        return (
          {
            status: "error",
            data: {
              message: msg.createChannelErr.mandatory
            }
          }
        )
      }
      if (typeof channelName !== "string") {
        return (
          {
            status: "error",
            data: {
              message: msg.createChannelErr.string
            }
          }
        )
      }
      //params are ok
      let command = `CREATE CHANNEL '${channelName}'`;
      command = ifNotExist ? command + " IF NOT EXISTS" : command;
      const response = await this.exec(command);
      if (response.status == "success") {
        return (
          {
            status: response.status,
            data: {
              message: response.data
            }
          }
        )
      }
      if (response.status == "error") {
        return response
      }
    } catch (error) {
      return {
        status: "error",
        data: {
          message: error.toString(),
          error: error
        }
      }
    }
  }

  /*
  removeChannel method calls exec method to remove an existing channel 
  */
  async removeChannel(channelName) {
    try {
      //params validation
      //check channelName has been provided
      if (!channelName) {
        return (
          {
            status: "error",
            data: {
              message: msg.removeChannelErr.mandatory
            }
          }
        )
      }
      if (typeof channelName !== "string") {
        return (
          {
            status: "error",
            data: {
              message: msg.removeChannelErr.string
            }
          }
        )
      }
      //params are ok
      let command = `REMOVE CHANNEL '${channelName}'`;
      const response = await this.exec(command);
      if (response.status == "success") {
        return (
          {
            status: response.status,
            data: {
              message: response.data
            }
          }
        )
      }
      if (response.status == "error") {
        return response
      }
    } catch (error) {
      return {
        status: "error",
        data: {
          message: error.toString(),
          error: error
        }
      }
    }
  }


  /* PRIVATE METHODS */

  /*
  makeid method generates unique IDs to use for requests
  */
  #makeid(length = 5) {
    var result = '';
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < length; i++) {
      result += characters.charAt(Math.floor(Math.random() *
        charactersLength));
    }
    return result;
  }

  /*
  connectWs private method is called to create the websocket conenction used to send command request, receive command request response, pubSub subscription request, 
  */
  #connectWs(url, errorMessage) {
    return new Promise((resolve, reject) => {
      var ws = new WebSocket(url);
      ws.onopen = function () {
        resolve(ws);
      };
      ws.onerror = function (err) {
        reject({
          err: err,
          message: errorMessage
        });
      };
    });
  }

  /*
  onCloseWs private method is called when close event is fired by websocket.
  if user provided a callback, this is invoked
  */
  #onCloseWs = (event) => {
    if (this.onClose !== null) {
      this.onClose(msg.wsCloseComplete);
      this.#ws.removeEventListener('close', this.#onCloseWs);
      this.#ws.removeEventListener('error', this.#onErrorWs);
      this.#ws = null;
    }
  }

  /*
  onErrorWs private method is called when error event is fired by websocket.
  if user provided a callback, this is invoked
  */
  #onErrorWs = (event) => {
    if (this.onError !== null) {
      this.onError(event, msg.wsOnError);
    }
  }

  /*
  onCloseWsPubSub private method is called when close event is fired by websocket.
  if user provided a callback, this is invoked
  */
  #onCloseWsPubSub = (event) => {
    if (this.onClose !== null) {
      this.onClose(msg.wsPubSubCloseComplete);
      this.#wsPubSub.removeEventListener('close', this.#onCloseWsPubSub);
      this.#wsPubSub.removeEventListener('error', this.#onErrorWsPubSub);
      this.#uuid = null;
      this.#wsPubSub = null;
    }
  }

  /*
  onError private method is called when error event is fired by websocket.
  if user provided a callback, this is invoked
  */
  #onErrorWsPubSub = (event) => {
    if (this.onError !== null) {
      this.onError(event, msg.wsPubSubOnError);
    }
  }

  /*
  pubsub method calls 
  */
  async #pubsub(type, channel, callback) {
    //based on the value of of callback, create a new subscription or remove the subscription
    if (callback !== null) {
      //check if the channel subscription is already active
      if (!this.#subscriptionsStack.has(channel)) {
        //if the subscription does not exist 
        try {
          let body;
          //based on type value build the right body
          //important in case of channel, the channel key is present in the body
          //important in case of table, the table key is present in the body
          if (type == "listenChannel") {
            body = {
              id: this.#makeid(5),
              type: "listen",
              channel: channel.toLowerCase(),
            }
          }
          if (type == "unlistenChannel") {
            body = {
              id: this.#makeid(5),
              type: "unlisten",
              channel: channel.toLowerCase(),
            }
          }
          if (type == "listenTable") {
            body = {
              id: this.#makeid(5),
              type: "listen",
              table: channel.toLowerCase(),
            }
          }
          if (type == "unlistenTable") {
            body = {
              id: this.#makeid(5),
              type: "unlisten",
              table: channel.toLowerCase(),
            }
          }
          const response = await this.#promiseSend(body);
          //if this is the first subscription, create the websocket connection dedicated to receive pubSub messages
          if (this.#subscriptionsStack.size == 0 && response.status == "success") {
            //response here we have authentication information
            const uuid = response.data.uuid;
            const secret = response.data.secret;
            try {
              this.#wsPubSubUrl = `wss://web1.sqlitecloud.io:8443/api/v1/wspsub?uuid=${uuid}&secret=${secret}`;
              this.#wsPubSub = await this.#connectWs(this.#wsPubSubUrl, "PubSub connection not established");
              //when the PubSub WebSocket is created the channel is added to the stack
              this.#subscriptionsStack.set(channel.toLowerCase(),
                {
                  channel: channel.toLowerCase(),
                  callback: callback
                }
              );
              this.#uuid = uuid;
              //register the onmessage event on pubSub websocket
              this.#wsPubSub.addEventListener('message', this.#wsPubSubonMessage);
              //register the close event on websocket
              this.#wsPubSub.addEventListener('error', this.#onErrorWsPubSub);
              //register the close event on websocket
              this.#wsPubSub.addEventListener('close', this.#onCloseWsPubSub);
            } catch (error) {
              return {
                status: "error",
                data: error
              };
            }
          }
          //if this isn't the first subscription, just add the supscription to the stack
          if (this.#subscriptionsStack.size > 0 && response.status == "success") {
            this.#subscriptionsStack.set(channel.toLowerCase(),
              {
                channel: channel.toLowerCase(),
                callback: callback
              }
            );
          }
          //build the object returned to client
          let userResponse = {};
          userResponse.status = response.status;
          if (response.status == "success") {
            userResponse.data = {};
            userResponse.data.channel = response.data.channel;
          }
          if (response.status == "error") {
            userResponse.data = response.data;
          }
          return userResponse;
        } catch (error) {
          return {
            status: "error",
            data: error
          };
        }
      } else {
        //if the subscription exists
        return (
          {
            status: "warning",
            data: {
              message: msg.wsListenError.alreadySubscribed + " " + channel
            }
          }
        )
      }
    } else {
      try {
        console.log("sto per chiamare this.#promiseSend") //TOGLI
        const response = await this.#promiseSend(
          {
            id: this.#makeid(5),
            type: type,
            channel: channel.toLowerCase(),
          }
        );
        console.log(response) //TOGLI
        this.#subscriptionsStack.delete(channel)
        //check the remaing active subscription. If zero the websocket connection used for pubSub can be closed
        if (this.#subscriptionsStack.size == 0) {
          this.#wsPubSub.removeEventListener('message', this.#wsPubSubonMessage);
          this.#wsPubSub.close(1000, msg.wsPubSubCloseByClient);
          this.#wsPubSub = null;
        }
        return (response);
      } catch (error) {
        return {
          status: "error",
          data: error
        };
      }
    }
  }

  /*
  wsPubSubonMessage private method is called when ad onmessage event is fired on pubSub websocket.
  based on the channel indicated in the message the right callback is called
  */
  #wsPubSubonMessage = (event) => {
    const pubSubMessage = JSON.parse(event.data);

    //since payload can be both a string or JSON, this function based on check of it is or not a valid JSON return the correct parsed paylod
    function payloadParser(str) {
      try {
        JSON.parse(str);
      } catch (e) {
        return str;
      }
      return JSON.parse(str);
    }

    //build the obj returned to the user removing fields not usefull
    const userPubSubMessage = {
      channel: pubSubMessage.channel,
      sender: pubSubMessage.sender,
      payload: payloadParser(pubSubMessage.payload),
      ownMessage: this.#uuid == pubSubMessage.sender
    }
    //this is the case in which the user decide to filter the message sent by himself
    if (this.filterSentMessages && this.#uuid == pubSubMessage.sender) {

    } else {
      this.#subscriptionsStack.get(pubSubMessage.channel).callback(userPubSubMessage);
    }
  }

  /*
  promiseSend private method send request to the server creating a Promise that resolve when a websocket onmessage event is fired.
  */
  #promiseSend(request) {
    console.log("STO PER CHIAMARE this.#ws.send") //TOGLI
    //request is sent to the server
    this.#ws.send(JSON.stringify(request));
    //extract request id
    const requestId = request.id;
    //define the Promise that wait for the server response 
    return new Promise((resolve, reject) => {
      //define what to do if an answer does not arrive within the set time
      const onRequestTimeout = setTimeout(() => { this.#handlePromiseRejectTimeout(requestId) }, this.requestTimeout);
      //add the new request to the request stack 
      this.#requestsStack.set(
        requestId,
        {
          id: requestId,
          onRequestTimeout: onRequestTimeout,
          resolve: resolve,
          reject: reject
        }
      )
      //if this is the only one request in the stack, register the function to be executed at the onmessage event
      if (this.#requestsStack.size == 1) this.#ws.addEventListener('message', this.#handlePromiseResolve);
    })
  }

  /*
  private handlePromiseResolve method is called when onmessage event is fired.
  */
  #handlePromiseResolve = (event) => {
    //parse the response sent by the server
    const response = JSON.parse(event.data);
    //search in the requests stack the request with the same id of the response received by the server.
    //it is possible that the request no longer exists as it may have timed out and therefore deleted from the stack.
    //if the request was found:
    // - the Promise corresponding to the request is resolved returning the response received by the server
    // - the timeout related to the request is cleared
    // - the new requests stack is stored
    // - if there are no pending requests remove the websocket onmessage event
    if (this.#requestsStack.has(response.id)) {
      //build the obj returned to the user removing based on type
      let userResponse = {};
      switch (response.type) {
        case "exec":
          userResponse = {
            status: response.status,
            data: response.data
          };
          break;
        case "notify":
          userResponse = {
            status: response.status,
            data: response.data
          };
          break;
        case "listen":
          //in this case the message is passed as is because the uuid and secret field, if present, will be used to create wsPubSub connection
          //the messega will be cleaned from this fields in #pubsub method
          userResponse = response;
          break;
        case "unlisten":
          userResponse = {
            status: response.status,
            data: response.data
          };
          break;
        default:
          userResponse = response;
      }

      this.#requestsStack.get(response.id).resolve(userResponse);
      clearTimeout(this.#requestsStack.get(response.id).onRequestTimeout);
      this.#requestsStack.delete(response.id);
      if (this.#requestsStack.size == 0) this.#ws.removeEventListener('message', this.#handlePromiseResolve);
    }
  }

  /*
  private handlePromiseRejectTimeout method is called when a request times out.
  */
  #handlePromiseRejectTimeout = (requestID) => {
    //search in the requests stack the request with the same id of the request that times out.
    //a new requests stack is created by removing the request that times out.
    //once the request is found:
    // - the Promise corresponding to the reject returning an error message
    // - the timeout related to the request is cleared
    // - the new requests stack is stored
    // - if there are no pending requests remove the websocket onmessage event
    if (this.#requestsStack.has(requestID)) {
      clearTimeout(this.#requestsStack.get(requestID).onRequestTimeout);
      this.#requestsStack.get(requestID).reject({
        message: msg.wsTimeoutError + " " + requestID
      });
      this.#requestsStack.delete(requestID);
      if (this.#requestsStack.size == 0) this.#ws.removeEventListener('message', this.#handlePromiseResolve);
    }
  }
}


