//core
import React, { Fragment, useState } from 'react';
//@heroicons/react
import { XCircleIcon, ExclamationTriangleIcon, CheckCircleIcon, InformationCircleIcon } from '@heroicons/react/20/solid'
//utils
import { logThis } from '../../js/utils';
/*
opt =  {
  description: "message describing errors"
  severity: error | warning | info | success
  messages: ["lorem ipsum", "lorem ipsum"]
}
*/
export default function Alert(props) {
  if (process.env.DEBUG == "true") logThis("Alert: ON RENDER");
  //extract params from opt
  const description = props.opt.description;
  const severity = props.opt.severity;
  const messages = props.opt.messages;
  //define colors based on severity
  let iconColor;
  let bgColor;
  let descriptionColor;
  let messagesColor;
  switch (severity) {
    case "error":
      iconColor = "text-red-400";
      descriptionColor = "text-red-800";
      messagesColor = "text-red-700";
      bgColor = "bg-red-50";
      break;
    case "warning":
      iconColor = "text-yellow-400";
      descriptionColor = "text-yellow-800";
      messagesColor = "text-yellow-700";
      bgColor = "bg-yellow-50";
      break;
    case "info":
      iconColor = "text-blue-400";
      descriptionColor = "text-blue-800";
      messagesColor = "text-blue-700";
      bgColor = "bg-blue-50";
      break;
    case "success":
      iconColor = "text-green-400";
      descriptionColor = "text-green-800";
      messagesColor = "text-green-700";
      bgColor = "bg-green-50";
      break;
    default:
      iconColor = "text-black";
      descriptionColor = "text-red-black";
      messagesColor = "text-black";
      bgColor = "bg-white";
      break;
  }
  //define elements classes based on severity
  const iconClasses = `h-5 w-5 ${iconColor}`
  const bgClasses = `rounded-md p-4 ${bgColor}`
  const descriptionClasses = `text-sm font-medium ${descriptionColor}`
  const messagesClasses = `mt-2 text-sm ${messagesColor}`
  //render UI
  return (
    <div className={bgClasses}>
      <div className="flex">
        <div className="flex-shrink-0">
          {
            severity === "error" &&
            <XCircleIcon className={iconClasses} aria-hidden="true" />
          }
          {
            severity === "warning" &&
            <ExclamationTriangleIcon className={iconClasses} aria-hidden="true" />
          }
          {
            severity === "info" &&
            <InformationCircleIcon className={iconClasses} aria-hidden="true" />
          }
          {
            severity === "success" &&
            <CheckCircleIcon className={iconClasses} aria-hidden="true" />
          }
        </div>
        <div className="ml-3">
          {
            description &&
            <h3 className={descriptionClasses}>{description}</h3>
          }
          {
            messages &&
            <>
              {
                messages.length == 1 &&
                <div className={messagesClasses}>
                  {
                    messages.map((message, i) => <p key={i}>{message}</p>)
                  }
                </div>
              }
              {
                messages.length > 1 &&
                <div className={messagesClasses}>
                  <ul role="list" className="list-disc space-y-1 pl-5">
                    {
                      messages.map((message, i) => <li key={i}>{message}</li>)
                    }
                  </ul>
                </div>
              }
            </>
          }
        </div>
      </div>
    </div>
  )
}
