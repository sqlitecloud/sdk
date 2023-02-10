//core
import React, { Fragment, useRef, useState, useEffect } from "react";
//react-router
import { useSearchParams } from 'react-router-dom';
//headlessui/react
// import { Listbox, Transition } from '@headlessui/react'
//heroicons/
// import { CalendarIcon, PaperClipIcon, TagIcon, UserCircleIcon } from '@heroicons/react/20/solid'
//utils
import {
  logThis
} from '../js/utils';
//components
import Alert from './alert/Alert'
import CircularLoader from './loaders/CircularLoader';
/* dummy data from demo component. can be usefull in future 
const assignees = [
  { name: 'Unassigned', value: null },
  {
    name: 'Wade Cooper',
    value: 'wade-cooper',
    avatar:
      'https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
  },
  // More items...
]
const labels = [
  { name: 'Unlabelled', value: null },
  { name: 'Engineering', value: 'engineering' },
  // More items...
]
const dueDates = [
  { name: 'No due date', value: null },
  { name: 'Today', value: 'today' },
  // More items...
]
*/
const MessageEditor = ({ client }) => {
  if (process.env.DEBUG == "true") logThis("MessageEditor: ON RENDER");
  //react router hooks used to set query string
  const [searchParams, setSearchParams] = useSearchParams();
  //store and set actual message written by the user
  const [value, setValue] = useState("");
  //handling change in textarea
  const handleChange = (event) => {
    if (!sendingMessage) {
      setValue(event.target.value);
    }
  };
  //init reference to textarea 
  const editorRef = useRef(null);
  //when query parameters change, set automaticcally focus on textarea 
  useEffect(() => {
    setValue("");
    setErrorSending(null);
    editorRef.current.focus();
  }, [searchParams])
  /* state not used at the moment. taken from demo component
  const [assigned, setAssigned] = useState(assignees[0])
  const [labelled, setLabelled] = useState(labels[0])
  const [dated, setDated] = useState(dueDates[0])
  */
  //method to handle key down
  //if pressed enter key the message is sent
  //if pressed enter+shift keys, a new line is created
  const handleKey = (event) => {
    if (event.keyCode == 13 && event.shiftKey) {
    } else if (event.keyCode == 13) {
      event.preventDefault();
      if (value !== "") sendMsg();
      return false;
    }
  };
  //state used to know if we are sending a message
  const [sendingMessage, setSendingMessage] = useState(false);
  //method used to send message
  const [errorSending, setErrorSending] = useState(null);
  const sendMsg = async (event) => {
    if (event) event.preventDefault();
    const queryChannel = searchParams.get("channel");
    if (queryChannel && !sendingMessage && value) {
      setSendingMessage(true);
      setErrorSending(null);
      const response = await client.notify(queryChannel, { message: value });
      console.log(response) //TOGLI
      if (response.status == "success") {
        setValue("");
      }
      if (response.status == "error") {
        setErrorSending(response.data.message);
      }
      setSendingMessage(false);
    }
    editorRef.current.focus();
  }
  //render ui
  return (
    <div className="flex flex-col">
      {
        errorSending &&
        <div className="mb-2">
          <Alert opt={{
            description: "message can't be sent",
            severity: "error",
            messages: [errorSending]
          }} />
        </div>
      }
      <div className="relative flex flex-row">
        <div className="flex-grow overflow-hidden rounded-lg border border-gray-300 shadow-sm focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500">
          <textarea
            ref={editorRef}
            rows={4}
            name="description"
            id="description"
            className="block w-full resize-none border-0 py-0 pt-2.5 placeholder-gray-500 focus:ring-0 sm:text-sm"
            placeholder="Write your message..."
            value={value}
            onChange={handleChange}
            onKeyDown={handleKey}
          />
        </div>
        <div className={`flex items-center ${sendingMessage ? "justify-between" : "justify-end"} space-x-3 px-2 py-2 sm:px-3`}>
          {
            // loader during sending
            sendingMessage &&
            <CircularLoader message={"sending message"} />
          }
          {
            !sendingMessage && 
            <button
              disabled={sendingMessage || !value}
              type="button"
              className={`inline-flex items-center rounded-md border border-transparent bg-indigo-600 ${sendingMessage || !value ? "opacity-50" : ""} px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2`}
              onClick={sendMsg}
            >
              Send
            </button>
          }
        </div>
      </div>
    </div>
  );
}


export default MessageEditor;