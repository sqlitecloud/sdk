//core
import React, { useState } from 'react';
//utils
import {
  logThis,
} from '../js/utils';
//components
import Alert from './alert/Alert';
import CircularLoader from './loaders/CircularLoader';
/*
opt =  {
}
*/
export default function CreateChannel(props) {
  if (process.env.DEBUG == "true") logThis("CreateChannel: ON RENDER");
  //exctract env variable
  var projectId = process.env.PROJECT_ID;
  //extract params from opt
  const client = props.client;
  const setReloadChannelsList = props.setReloadChannelsList;
  const reloadChannelsList = props.reloadChannelsList;
  //hadle input channel name
  const [isError, setIsError] = useState(null);
  const [channelName, setChannelName] = useState("");
  const handleChange = (event) => {
    setChannelName(event.target.value);
  }
  //hadle click that submit channel creation
  const [isCreatingChannel, setIsCreatingChannel] = useState(false);
  const createChannel = async () => {
    setIsCreatingChannel(true);
    setIsError(null);
    const response = await client.createChannel(channelName, true);
    if (process.env.DEBUG == "true") console.log(response);
    if (response.status == "success") {
      //if createChannel is succesful reload channelsList
      setReloadChannelsList(!reloadChannelsList);
      const response = await client.notify(projectId, { action: "create", channel: channelName });
    } else {
      setIsError(response.data.message);
    }
    setChannelName("");
    setIsCreatingChannel(false);
  }
  const handleOnClick = async () => {
    await createChannel();
  }
  //method to handle key down
  //if pressed enter key the message is sent
  //if pressed enter+shift keys, a new line is created
  const handleKey = async (event) => {
    if (event.keyCode == 13) {
      await createChannel();
    } else {
      return false;
    }
  };
  //render UI
  return (
    <div className="mx-4 my-4 bg-white shadow sm:rounded-lg">
      <div className="px-4 py-2">
        <h4 className="text-base font-medium text-gray-900">Create channel</h4>
        <div className="mt-3 sm:flex">
          {
            isCreatingChannel &&
            <CircularLoader message={"creating channel " + channelName} />
          }
          {
            !isCreatingChannel &&
            <>
              <div className="w-full">
                <div className="rounded-md border border-gray-300 px-3 py-1 shadow-sm bg-white focus-within:border-indigo-600 focus-within:ring-1 focus-within:ring-indigo-600">
                  <label htmlFor="name" className="block text-xs font-medium text-gray-900">
                    Channel name
                  </label>
                  <input
                    onChange={handleChange}
                    type="text"
                    name="name"
                    id="name"
                    value={channelName}
                    className="block w-full border-0 p-0 text-gray-900 placeholder-gray-500 focus:ring-0 sm:text-sm"
                    placeholder="channel_123"
                    onKeyDown={handleKey}
                  />
                </div>
              </div>
              {
                channelName &&
                <button
                  onClick={handleOnClick}
                  type="button"
                  className="mt-3 inline-flex w-full items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                >
                  Create
                </button>
              }
            </>
          }
        </div>
        {
          isError &&
          <div className='w-full flex-1 mt-2'>
            <Alert opt={{
              description: "error creating channel " + channelName,
              severity: "error",
              messages: [isError]
            }} />
          </div>
        }
      </div>
    </div>
  )
}
