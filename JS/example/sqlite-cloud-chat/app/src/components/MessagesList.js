//core
import React, { useRef, useEffect } from 'react';
//react-json-view
import { JsonViewer } from "@textea/json-viewer"
//utils
import {
  logThis,
} from '../js/utils';


function NewlineText({ text }) {
  let newText = "";
  if (text !== null) {
    newText = text.split('\n').map((str, i) => {
      if (str == "") {
        return (
          <br key={i} />
        )
      } else {
        return (
          <p key={i} className="mt-0.5 text-sm text-gray-500">
            {str}
          </p>
        )
      }
    });
    return newText;
  }
}

export default function MessagesList(props) {
  if (process.env.DEBUG == "true") logThis("MessagesList: ON RENDER");
  //extract props
  const messages = props.messages;
  const showEditor = props.showEditor;
  //init ref to the messages container
  const contRef = useRef(null);
  //on each new messages scroll to bottom the container
  useEffect(() => {
    contRef.current.scrollTo({ left: 0, top: contRef.current.scrollHeight, behavior: 'smooth' });
  }, [messages]);
  //render ui
  return (
    <div ref={contRef} className="absolute bottom-0 w-[95%] max-h-full overflow-auto py-4" >
      {
        messages &&
        <ul role="list">
          {messages.map((message) => {
            const bgColor = message.ownMessage ? "bg-green-500" : "bg-indigo-500";
            return (
              <li key={message.user + message.time} className="flex flex-row p-4 my-4 bg-gray-50 rounded-md transition-all duration-200">
                <div className="relative">
                  <div className={`relative flex items-start space-x-3`}>
                    <div className="relative">
                      <div className={`relative ${bgColor} flex-shrink-0  flex items-center justify-center w-10 h-10 text-white text-sm font-medium rounded-full`}>
                        {message.sender.slice(0, 2).toUpperCase()}
                      </div>
                    </div>
                    <div className="min-w-0 flex-1">
                      <div>
                        <div className="text-sm font-medium text-gray-900">
                          {message.sender}
                        </div>
                        <p className="mt-0.5 text-xs text-gray-500">{message.time}</p>
                        <div className="w-6 h-px my-1 bg-gray-500"></div>
                        <div>
                          {
                            showEditor && message.payload &&
                            <NewlineText text={message.payload.message} />
                          }
                          {
                            !showEditor && message.payload &&
                            <JsonViewer value={message.payload} />
                          }
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            )
          })}
        </ul>
      }
    </div>
  )
}
