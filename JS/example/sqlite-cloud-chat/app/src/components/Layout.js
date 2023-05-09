//core
import React, { Fragment, useRef, useContext, useState, useEffect } from 'react';
//react-router
import { useSearchParams } from 'react-router-dom';
//SQLiteCloud
import SQLiteCloud from 'sqlitecloud-sdk'
/*
usare client come nome dell'istanza locale 
*/
//@heroicons/react
import {
  Bars3Icon,
  XMarkIcon
} from '@heroicons/react/24/outline'
//utils
import {
  logThis,
  checkChannelExist
} from '../js/utils';
//image
import lightLogo from '../../assets/logo/logo-dark@4x.png';
//context
import { StateContext } from "../context/StateContext"
//components
import Alert from './alert/Alert'
import CircularLoader from './loaders/CircularLoader'
import CreateChannel from './CreateChannel'
import ChannelsList from './ChannelsList'
import MessageEditor from './MessageEditor'
import MessagesList from './MessagesList'

export default function Layout(props) {
  if (process.env.DEBUG == "true") logThis("Layout: ON RENDER");
  //read from context all channels registered  
  const { chsMap } = useContext(StateContext);
  //credentials to establish the websocket connection
  var projectId = process.env.PROJECT_ID;
  var apikey = process.env.API_KEY;
  //state dedicated to SQLiteCloud instance used to handle websocket connection
  const [client, setClient] = useState(null);
  const clientRef = useRef(null);
  clientRef.current = client;
  const [isConnecting, setIsConnecting] = useState(true);
  const [connectionResponse, setConnectionResponse] = useState(null);
  const [isLoadingChannels, setIsLoadingChannels] = useState(false);
  const [channelsListResponse, setChannelsListResponse] = useState(null);
  //callback function passed to websocket and register on error and close events
  let onErrorCallback = function (event, msg) {
    console.log(msg);
  }
  let onCloseCallback = function (msg) {
    console.log(msg);
  }
  //react router hooks used to set & get query string
  const [searchParams, setSearchParams] = useSearchParams();
  //if present, this is the database whose tables you want to register to
  const queryDBName = searchParams.get('dbName');
  //if present, this is the channel you want to listen to 
  const queryChannel = searchParams.get('channel');
  //check if listening to the actual selected channel is completed
  const isListeningActualChannel = chsMap ? chsMap.get(queryChannel) : null;
  //state to handle actual selected channel
  const [selectedChannel, setSelectedChannel] = useState(queryChannel);
  const [selectedChannelIndex, setSelectedChannelIndex] = useState(-1);
  //state used to store the available available channels. In case queryDBName !== null available channels are db tables
  const [channelsList, setChannelsList] = useState(undefined);
  //based on value of query channel show or not Messages component
  const [showMessages, setShowMessages] = useState(false);
  //show message editor and other ui element based on query parameters editor is shown when queryDBName is null
  const [showEditor, setShowEditor] = useState(false);
  //state used to open or close mobile sidebar holding messages
  const [openMobMsg, setOpenMobMsg] = useState(false);
  //method that handle channelsListResponse
  const handlingChannelsListResponse = (channelsListResponse) => {
    setChannelsListResponse(channelsListResponse);
    //check how listChannels completed  
    if (channelsListResponse.status == 'success') {
      //successful listChannels
      if (process.env.DEBUG == 'true') logThis('Received channels list');
      var channels = [];
      if (queryDBName !== null) {
        channelsListResponse.data.rows.forEach(c => {
          channels.push(c.chname);
        })
        setChannelsList(channels);
        setShowEditor(false);
      } else {
        if (channelsListResponse.data.rows == undefined) {
          setChannelsList([]);
        } else {
          channelsListResponse.data.rows.forEach(c => {
            channels.push(c.chname);
          })
          setChannelsList(channels);
        }
        setShowEditor(true);
      }
      //check if the channel in query string exist
      const testChannelExist = checkChannelExist(channels, queryChannel);
      if (testChannelExist !== -1) {
        //if true show message components
        setShowMessages(true);
        setSelectedChannelIndex(testChannelExist);
        setOpenMobMsg(true);
      } else {
        //if false not show message components and remove query string from url
        if (queryDBName !== null) {
          setSearchParams({
            dbName: queryDBName
          });
        } else {
          setSearchParams({});
        }
        setShowMessages(false);
        setSelectedChannelIndex(-1);
        setOpenMobMsg(false);
      }
    } else {
      //error on listChannels
      if (process.env.DEBUG == 'true') logThis(channelsListResponse.data.message);
    }
    setIsLoadingChannels(false);
  }
  //useEffect triggered only onMount
  //create websocket connection
  //retrieve channels list, or tables list if queryDBName is present as query string
  useEffect(() => {
    const onMountWrapper = async () => {
      if (process.env.DEBUG == 'true') logThis('App: ON useEffect');
      //init SQLiteCloud instance using provided credentials
      let localClient = new SQLiteCloud(projectId, apikey, onErrorCallback, onCloseCallback);
      //set websocket request timeout
      localClient.setRequestTimeout(5000);
      //try to enstablish websocket connection
      const connectionResponse = await localClient.connect();
      setConnectionResponse(connectionResponse);
      setIsConnecting(false);
      if (process.env.DEBUG == 'true') console.log(connectionResponse);
      //check how websocket connection completed  
      if (connectionResponse.status == 'success') {
        //successful websocket connection
        setClient(localClient)
        //based on query parameters select if retrieve tables db or channels
        //in case of db, tables will be saved as channels
        let channelsListResponse = null;
        setIsLoadingChannels(true);
        if (queryDBName !== null) {
          const execMessage = `USE DATABASE ${queryDBName}; LIST TABLES PUBSUB`
          channelsListResponse = await localClient.exec(execMessage);
        } else {
          channelsListResponse = await localClient.listChannels();
        }
        handlingChannelsListResponse(channelsListResponse);
        //create if not exists a reserved channel deditaced to system communication
        const response = await localClient.createChannel(projectId, true);
      } else {
        //error on websocket connection
        if (process.env.DEBUG == 'true') logThis(connectionResponse.data.message);
      }
    }
    onMountWrapper();
  }, []);
  //useEffect triggered every time selectedChannel changes
  useEffect(() => {
    var testChannelExist = checkChannelExist(channelsList, selectedChannel);
    if (testChannelExist == -1) {
      setShowMessages(false);
      setSelectedChannelIndex(testChannelExist);
      setOpenMobMsg(false);
    }
    if (testChannelExist != -1) {
      setShowMessages(true);
      setSelectedChannelIndex(testChannelExist);
    }
  }, [selectedChannel]);
  //section dedicated to handling channels list after a channel is created o dropped
  const [reloadChannelsList, setReloadChannelsList] = useState(false);
  const reloadChannelsListRef = useRef(null);
  reloadChannelsListRef.current = reloadChannelsList;
  useEffect(() => {
    const onMountWrapper = async () => {
      if (client) {
        setIsLoadingChannels(true);
        let channelsListResponse = await client.listChannels();
        handlingChannelsListResponse(channelsListResponse);
      }
    }
    onMountWrapper();
  }, [reloadChannelsList]);
  //method called to close channel
  const closeChannel = () => {
    //if false not show message components and remove query string from url
    if (queryDBName !== null) {
      setSearchParams({
        dbName: queryDBName
      });
    } else {
      setSearchParams({});
    }
    setShowMessages(false);
    setSelectedChannelIndex(-1);
    setSelectedChannel(null);
    setOpenMobMsg(false);
  }
  //useEffect triggered when chsMap change
  useEffect(() => {
    let lastMessagesTimestamp = [];
    //build an array with the last received message for each channel
    chsMap.forEach((ch, key) => {
      const lastChEl = ch[ch.length - 1];
      if (lastChEl) {
        lastMessagesTimestamp.push({
          name: key,
          time: lastChEl.timeMs
        });
      } else {
        lastMessagesTimestamp.push({
          name: key,
          time: 0
        });
      }
    });
    //reordered the array from the old messages to the newest
    lastMessagesTimestamp.sort(function (a, b) {
      return ((a.time < b.time) ? -1 : ((a.time == b.time) ? 0 : 1));
    });
    //reverse the array to have as first the newest channel
    lastMessagesTimestamp = lastMessagesTimestamp.reverse();
    //reorder channelsList following the new order based on last messages
    if (channelsList) {
      let newChannelsList = [];
      let newSelectedIndex = -1;
      lastMessagesTimestamp.forEach((lastMsg, i) => {
        if (lastMsg.name === selectedChannel) newSelectedIndex = i;
        newChannelsList.push(lastMsg.name);
      })
      //check if all the original channels are present in the new array
      //if a channel is not present is added to the array
      //this can happen when there is a problem in listening a specific channel
      channelsList.forEach((channel) => {
        if (newChannelsList.indexOf(channel) === -1) {
          newChannelsList.push(channel);
        }
      })
      setChannelsList(newChannelsList);
      setSelectedChannelIndex(newSelectedIndex);
    }
  }, [chsMap])
  //handle reconnection of windows focus
  useEffect(() => {
    const handleWindowFocusEvent = async () => {
      if (clientRef.current) {
        console.log(clientRef.current.connectionState);
        const connectionState = clientRef.current.connectionState;
        if (connectionState.state !== 1) {
          setIsConnecting(true);
          //try to enstablish websocket connection
          const connectionResponse = await clientRef.current.connect();
          setConnectionResponse(connectionResponse);
          setIsConnecting(false);
          //check how websocket connection completed  
          if (connectionResponse.status == 'success') {
            setReloadChannelsList(!reloadChannelsListRef.current);
          } else {
            //error on websocket connection
            if (process.env.DEBUG == 'true') logThis(connectionResponse.data.message);
          }
        }
      }
    }
    window.addEventListener("focus", handleWindowFocusEvent);
    return () => {
      window.removeEventListener("focus", handleWindowFocusEvent);
    };
  }, []);
  //show or not sidebar on mobile
  const showSideBarOnMob = openMobMsg || selectedChannelIndex == -1;
  //render UI
  return (
    <div className='h-full'>
      {/* Static sidebar for desktop */}
      <div className={`z-[5000] fixed inset-y-0 transform transition-all duration-300 ${showSideBarOnMob ? "translate-x-0" : "-translate-x-full"} sm:translate-x-0 sm:flex w-full sm:w-96 flex-col`}>
        <div className='flex flex-col h-full bg-gray-100'>
          {/* Sidebar component, swap this element with another sidebar if you like */}
          <div className='flex-grow-0 mb-2'>
            <img
              className='box-content h-8 px-4 py-2 w-auto'
              src={lightLogo}
              alt=''
            />
            {/* TEST CLOSE CONNECTION TOGLI*/}
            <button onClick={()=>{client.close()}}>CLOSE</button>
            <button onClick={()=>{console.log(client.subscriptionsStackState);}}>List active pub sub</button>
          </div>
          {
            showEditor &&
            <div className='flex-grow-0'>
              <CreateChannel client={client} reloadChannelsList={reloadChannelsList} setReloadChannelsList={setReloadChannelsList} />
            </div>
          }
          {
            isConnecting &&
            <div className='mx-4'>
              <CircularLoader message={"connecting"} />
            </div>
          }
          <div className='flex min-h-0 flex-1 flex-col bg-gray-100'>
            <div className='flex flex-1 flex-col overflow-y-hidden px-2 sm:px-0 pb-4'>
              <ChannelsList client={client} showEditor={showEditor} isLoadingChannels={isLoadingChannels} channelsList={channelsList} reloadChannelsList={reloadChannelsList} setReloadChannelsList={setReloadChannelsList} selectedChannelIndex={selectedChannelIndex} setSelectedChannelIndex={setSelectedChannelIndex} setSelectedChannel={setSelectedChannel} setOpenMobMsg={setOpenMobMsg} />
            </div>
            {/* Sidebar footer for desktop */}
            <div className='flex flex-shrink-0 bg-gray-800 p-4'>
              <div className='group block w-full flex-shrink-0'>
                <div className='flex items-center'>
                  <div className='ml-3'>
                    <p className='text-xs font-medium text-white'>project ID:</p>
                    <p className='mt-1 text-xs font-normal text-gray-400'>{projectId}</p>
                  </div>
                </div>
                {
                  queryDBName &&
                  <div className='mt-2 ml-3'>
                    <p className='text-xs font-medium text-white'>database name:</p>
                    <p className='mt-1 text-xs font-normal text-gray-400'>{queryDBName}</p>
                  </div>
                }
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Layout right section */}
      <div className='h-full max-h-full flex flex-1 flex-col justify-end sm:pl-96'>
        {
          connectionResponse && connectionResponse.status == "error" &&
          <div className='w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 overflow-y-auto'>
            <Alert opt={{
              description: "websocket connection errors",
              severity: "error",
              messages: [connectionResponse.data.message]
            }} />
          </div>
        }
        {
          connectionResponse && connectionResponse.status == "success" && !showMessages &&
          <div className='w-full h-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 overflow-y-auto'>
            <div className='max-w-xl mx-auto h-full flex flex-col items-center justify-start'>
              <div className='mt-[40%] text-lg font-semibold text-center text-gray-600'>
                Welcome to SQLite Cloud Chat!
              </div>
              <div className='mt-12 text-base text-gray-700 flex flex-col space-y-6'>
                <div>
                  This sample Web App demonstrates the power of the <a href="https://docs.sqlitecloud.io/docs/introduction/pubsub_implementation" className="text-blue-600" target="_blank">Pub/Sub</a> capabilities built into SQLite Cloud.
                </div>
                <div>
                  Open this page on two or more separate devices and try to send messages to different channels.
                </div>
                <div>
                  The underline <a href="https://docs.sqlitecloud.io/docs/sdk" className="text-blue-600" target="_blank">Javascript SDK</a> communicates with an SQLite Cloud cluster, and it uses the standard PUB/SUB commands to efficiently broadcast messages between the channels.
                </div>
                <div>
                  Pub/Sub is also implemented on the database level so you can LISTEN to a database table and start receiving JSON payloads each time that table changes.
                </div>
              </div>
            </div>
          </div>
        }
        {
          channelsListResponse && channelsListResponse.status == "error" &&
          <div className='w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 overflow-y-auto'>
            <Alert opt={{
              description: "channels list errors",
              severity: "error",
              messages: [channelsListResponse.data.message]
            }} />
          </div>
        }
        {
          showMessages && isListeningActualChannel &&
          <div className='sticky top-0 w-full pt-4 '>
            <div className='flex flex-row justify-between max-w-7xl mx-auto border-b'>
              <h1 className='text-2xl font-semibold text-gray-900 pb-2 px-4 sm:px-6 lg:px-8'>{selectedChannel}</h1>
              <button
                type='button'
                className='ml-1 flex h-10 w-10 items-center justify-center rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 group'
                onClick={closeChannel}
              >
                <span className='sr-only'>Close channel</span>
                <XMarkIcon className='h-6 w-6 text-gray-900 group-hover:text-gray-500' aria-hidden='true' />
              </button>
            </div>
          </div>
        }
        <div className='relative flex-grow w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 my-4'>
          {
            showMessages && isListeningActualChannel &&
            <>
              {
                chsMap && chsMap.size > 0 && <MessagesList messages={chsMap.get(selectedChannel)} showEditor={showEditor} />
              }
            </>
          }
        </div>
        {
          showEditor && showMessages && isListeningActualChannel &&
          <div className='w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8'>
            <MessageEditor client={client} chsMap={chsMap} />
          </div>
        }
      </div >
    </div >
  )
}
