//core
import React, { useRef, useEffect, useState, useContext } from "react";
//react-router
import { useSearchParams } from 'react-router-dom';
//date-fns
import { format, getTime } from 'date-fns';
//utils
import {
  logThis
} from '../js/utils';
//context
import { StateContext } from "../context/StateContext";
//components
import ChannelElementDropdown from './ChannelElementDropdown';
import CircularLoader from './loaders/CircularLoader';

const ChannelElement = (props) => {
  if (process.env.DEBUG == "true") logThis("ChannelElement: ON RENDER");
  //exctract env variable
  var projectId = process.env.PROJECT_ID;
  //extract props
  const client = props.client;
  const index = props.index;
  const name = props.name;
  const showEditor = props.showEditor;
  const selectionState = props.selectionState;
  const setSelectedChannel = props.setSelectedChannel;
  const setSelectedChannelIndex = props.setSelectedChannelIndex;
  const setOpenMobMsg = props.setOpenMobMsg;
  const reloadChannelsList = props.reloadChannelsList;
  const reloadChannelsListRef = useRef(null);
  reloadChannelsListRef.current = reloadChannelsList;
  const setReloadChannelsList = props.setReloadChannelsList;
  //ref to the element
  const channelButtonRef = useRef(null);
  //react router hooks used to set query string
  const [searchParams, setSearchParams] = useSearchParams();
  //if present, this is the database whose tables you want to register to
  const queryDBName = searchParams.get("dbName");
  //this state store the listen command result
  const [listenResponse, setListenResponse] = useState(null);
  //timestamp last message
  const [msgTimestamp, setMsgTimestamp] = useState("");
  //read from context state dedicated to save all received messages
  const { chsMapRef, setChsMap } = useContext(StateContext);
  const [prevMsgLenght, setPrevMsgLenght] = useState(0);
  const [alertNewMsg, setAlertNewMsg] = useState(0);
  //we need to create a reference to context state since the listen callback is called inside an event listener
  //see here https://medium.com/geographit/accessing-react-state-in-event-listeners-with-usestate-and-useref-hooks-8cceee73c559
  //whene a new message arrives it is added in the chsMap in correspondence of the element which has the channel name as its key
  const listen = async (message) => {
    console.log(message)
    if (message.channel == name && message.channel !== projectId && message.payload) {
      let newChsMap = new Map(JSON.parse(JSON.stringify(Array.from(chsMapRef.current))));
      let newMessages = JSON.parse(JSON.stringify(newChsMap.get(name)));
      message.time = format(new Date(), "yyyy-MM-dd' | 'HH:mm:ss")
      message.timeMs = getTime(new Date())
      setMsgTimestamp(message.time);
      newMessages.push(message);
      newChsMap.set(name, newMessages)
      chsMapRef.current = newChsMap;
      setChsMap(newChsMap);
    }
    if (
      (message.channel == name && message.channel == projectId && !message.ownMessage)
      // ||
      // (message.channel == name && message.channel !== projectId && message.type == "REMOVE" && !message.ownMessage)
    ) {
      console.log("HERE")
      const payload = message.payload;
      const action = payload.action;
      const channel = payload.channel;
      if (process.env.DEBUG == "true") logThis("ChannelElement: service message received");
      switch (action) {
        case 'remove':
          if (process.env.DEBUG == "true") logThis("ChannelElement: removed channel " + channel);
          //verify if the dropped channel is the one selected
          setIsDroppingChannel(true);
          setIsErrorDropping(null);
          if (selectionState) {
            //if the channel is the one selected
            //remove from the query string the channel name, but veryfing that the db name if present is not removed
            if (queryDBName !== null) {
              setSearchParams({
                dbName: queryDBName
              });
            } else {
              setSearchParams({});
            }
            // setShowMessages(false);
            setSelectedChannelIndex(-1);
            setSelectedChannel("");
          }
          //unlisten the channel
          await client.unlistenChannel(channel);
          //remove from the Maps the dropped channel
          let newChsMap = new Map(JSON.parse(JSON.stringify(Array.from(chsMapRef.current))));
          newChsMap.delete(channel);
          chsMapRef.current = newChsMap;
          setChsMap(newChsMap);
          setReloadChannelsList(!reloadChannelsListRef.current);
          break;
        case 'create':
          if (process.env.DEBUG == "true") logThis("ChannelElement: created channel " + channel);
          setReloadChannelsList(!reloadChannelsListRef.current);
          break;
        default:
          logThis("ChannelElement: action not defined");
      }
    }
  }
  //when loading for the first time start listen for incoming messages 
  //then check if the actual query params is equal to channel name
  //in this case set the selected channel equal to the channel index
  useEffect(() => {
    const registerToCh = async () => {
      let response;
      if (showEditor) {
        response = await client.listenChannel(name, listen);
      } else {
        response = await client.listenTable(name, listen);
      }
      setListenResponse(response);
      const queryChannel = searchParams.get("channel");
      if (queryChannel == name) {
        setSelectedChannelIndex(index)
      };
      let newChsMap = new Map(JSON.parse(JSON.stringify(Array.from(chsMapRef.current))));
      newChsMap.set(name, []);
      chsMapRef.current = newChsMap;
      setChsMap(newChsMap);
    }
    registerToCh();
  }, [])
  //listen again
  useEffect(() => {
    const listenAgain = async () => {
      let response;
      if (showEditor) {
        response = await client.listenChannel(name, listen);
      } else {
        response = await client.listenTable(name, listen);
      }
      setListenResponse(response);
    }
    listenAgain();
  }, [reloadChannelsList])
  //update the counter indicating the unread messages everytime the chsMap changes
  useEffect(() => {
    if (chsMapRef.current.get(name)) {
      if (!selectionState && chsMapRef.current.get(name).length !== prevMsgLenght) {
        setAlertNewMsg(chsMapRef.current.get(name).length - prevMsgLenght);
      } else {
        setAlertNewMsg(0);
        setPrevMsgLenght(chsMapRef.current.get(name).length);
      }
    }
  }, [chsMapRef.current])
  //method called every time a channel is selected
  //whene a channel is selected the counter for unread messages is set to zero
  //the current channel name and current channel index are setted to the current clicked channel
  //query strings are updated based on the clicked channel   
  const updateSelectChannel = (event) => {
    channelButtonRef.current.blur();
    setAlertNewMsg(0);
    setSelectedChannelIndex(index);
    setPrevMsgLenght(chsMapRef.current.get(name).length);
    setSelectedChannel(name);
    if (queryDBName !== null) {
      setSearchParams({
        dbName: queryDBName,
        channel: name,
      });
    } else {
      setSearchParams({
        channel: name
      });
    }
    setOpenMobMsg(false);
  }
  //hadle click that drop a channel
  const [isDroppingChannel, setIsDroppingChannel] = useState(false);
  const [isErrorDropping, setIsErrorDropping] = useState(null);
  const removeChannel = async (event) => {
    event.stopPropagation();
    //verify if the dropped channel is the one selected
    setIsDroppingChannel(true);
    setIsErrorDropping(null);
    if (selectionState) {
      //if the channel is the one selected
      //remove from the query string the channel name, but veryfing that the db name if present is not removed
      if (queryDBName !== null) {
        setSearchParams({
          dbName: queryDBName
        });
      } else {
        setSearchParams({});
      }
      // setShowMessages(false);
      setSelectedChannelIndex(-1);
      setSelectedChannel("");
    }
    const response = await client.removeChannel(name);
    if (process.env.DEBUG == "true") console.log(response);
    if (response.status == "success") {
      await client.unlistenChannel(name);
      //remove from the Maps the dropped channel
      let newChsMap = new Map(JSON.parse(JSON.stringify(Array.from(chsMapRef.current))));
      newChsMap.delete(name);
      chsMapRef.current = newChsMap;
      setChsMap(newChsMap);
      //if createChannel is succesful reload channelsList
      setReloadChannelsList(!reloadChannelsList);
      //notify on service channel that the channel has been removed
      const response = await client.notify(projectId, { action: "remove", channel: name }); //METTI
    } else {
      setIsDroppingChannel(false);
      setIsErrorDropping(response.data.message)
    }
  }
  //render UI
  if (name !== projectId) {
    return (
      <div
        onClick={updateSelectChannel}
        key={name}
        className={`relative flex items-center rounded-lg border border-gray-300 bg-white ${selectionState ? "ring-2 ring-green-500 ring-offset-2" : ""} px-2 md:px-6 py-5 shadow-sm focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 hover:border-gray-400`}>
        {
          // badge for unread messages
          (listenResponse && listenResponse.status == "success" && alertNewMsg !== 0) &&
          <div className="absolute top-0 left-1">
            <span className="inline-flex items-center rounded-full bg-indigo-100 px-2 py-0.5 text-[0.65rem] font-medium text-indigo-800">
              {alertNewMsg}
            </span>
          </div>
        }
        {
          // loader during listen
          !listenResponse &&
          <CircularLoader message={"listening channel " + name} />
        }
        {
          // loader during dropping
          isDroppingChannel &&
          <CircularLoader message={"removing channel " + name} />
        }
        {
          listenResponse && !isDroppingChannel &&
          <div className='relative bg-green-500 flex-shrink-0 flex items-center justify-center w-10 h-10 text-white text-sm font-medium rounded-full'>
            {name.slice(0, 2).toUpperCase()}
            {
              // error indicator on listen
              listenResponse.status == "error" &&
              <span className="box-content flex h-3 w-3 rounded-full absolute -bottom-1 -right-1 border-2 border-white">
                <span className="absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
              </span>
            }
            {
              // success indicator on listen
              listenResponse.status == "success" &&
              <span className="box-content flex h-3 w-3 rounded-full absolute -bottom-1 -right-1 border-2 border-white">
                <span className={`${alertNewMsg !== 0 ? "animate-ping" : ""} absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75`}></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-sky-500"></span>
              </span>
            }
          </div>
        }
        {
          listenResponse && !isDroppingChannel &&
          <div className="ml-3 min-w-0 flex-1">
            <button
              type="button"
              ref={channelButtonRef}
              className="focus:outline-none text-left">
              <span className="absolute inset-0" aria-hidden="true" />
              <p className="text-sm font-medium text-gray-900">{name}</p>
              {
                !listenResponse &&
                <p className="truncate text-xs text-gray-500">connecting...</p>
              }
              {
                listenResponse && listenResponse.status == "error" && !isErrorDropping &&
                <p className="truncate text-xs text-gray-500">{listenResponse.data.message.toLowerCase()}</p>
              }
              {
                listenResponse && listenResponse.status == "success" && !isErrorDropping &&
                <p className="truncate text-xs text-gray-500">{msgTimestamp}</p>
              }
              {
                listenResponse && isErrorDropping &&
                <p className="truncate text-xs text-gray-500">{isErrorDropping}</p>
              }
            </button>
          </div>
        }
        {
          showEditor &&
          <div className="absolute top-1 right-0">
            <ChannelElementDropdown removeChannel={removeChannel} />
          </div>
        }
      </div>
    );
  }
}


export default ChannelElement;