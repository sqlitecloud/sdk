//core
import React, { Fragment, useState } from 'react';
//utils
import {
  logThis,
} from '../../js/utils';
/*
opt =  {
}
*/
export default function CircularLoader(props) {
  if (process.env.DEBUG == "true") logThis("CircularLoader: ON RENDER");
  //extract props
  const message = props.message;
  //render UI
  return (
    <div className='flex flex-row space-x-1'>
      <svg className="animate-spin -ml-1 mr-3 h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      {
        message &&
        <span className="inline-flex items-center rounded-md bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800">
          {message}
        </span>
      }
    </div>
  )
}
