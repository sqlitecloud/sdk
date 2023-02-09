//core
import React, { useContext } from "react";
//utils
import { logThis } from '../js/utils';
//components
import ChannelElement from './ChannelElement'
import CircularLoader from './loaders/CircularLoader'

const ChannelsList = (props) => {
  if (process.env.DEBUG == "true") logThis("ChannelsList: ON RENDER");
  //extract props
  const liter = props.liter;
  let channelsList = props.channelsList;
  const showEditor = props.showEditor;
  const isLoadingChannels = props.isLoadingChannels;
  const setSelectedChannel = props.setSelectedChannel;
  const setOpenMobMsg = props.setOpenMobMsg;
  const selectedChannelIndex = props.selectedChannelIndex;
  const setSelectedChannelIndex = props.setSelectedChannelIndex;
  const reloadChannelsList = props.reloadChannelsList;
  const setReloadChannelsList = props.setReloadChannelsList;
  //render UI
  return (
    <nav className='py-2 space-y-3 px-2 md:px-4 min-h-0 transition-all duration-500 overflow-y-auto'>
      {
        isLoadingChannels &&
        <CircularLoader message={"loading channels list"} />
      }
      {
        channelsList &&
        <>
          {
            channelsList.length == 0 &&
            <span className="inline-flex items-center rounded-md bg-indigo-100 px-2.5 py-0.5 text-sm font-medium text-indigo-800">
              Doesn't exist any channels for this project. Create one!
            </span>
          }
          {
            channelsList.length > 0 &&
            <>
              {
                channelsList.map((channel, i) =>
                  <ChannelElement key={channel} index={i} liter={liter} name={channel} showEditor={showEditor} reloadChannelsList={reloadChannelsList} setReloadChannelsList={setReloadChannelsList} selectionState={selectedChannelIndex == i} setSelectedChannelIndex={setSelectedChannelIndex} setSelectedChannel={setSelectedChannel} setOpenMobMsg={setOpenMobMsg} />
                )
              }
            </>
          }
        </>
      }
    </nav>
  );
}


export default ChannelsList;