//core
import React, { Fragment, useState } from 'react';
//@headlessui
import { Menu, Transition } from '@headlessui/react'
//@heroicons
import { EllipsisVerticalIcon } from '@heroicons/react/20/solid'
//utils
import {
  logThis,
  classNames,
} from '../js/utils';
export default function ChannelElementDropdown(props) {
  if (process.env.DEBUG == "true") logThis("ChannelElementDropdown: ON RENDER");
  //extract params from opt
  const dropChannel = props.dropChannel;
  //hadle click that opens dropdown menu
  const handleOpenDropdown = (event) => {
    event.stopPropagation();
  }
  //render UI
  return (
    <Menu as="div" className="relative inline-block text-left">
      <div>
        <Menu.Button
          onClick={handleOpenDropdown}
          className="flex items-center rounded-full  text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-gray-100">
          <span className="sr-only">Open options</span>
          <EllipsisVerticalIcon className="h-5 w-5" aria-hidden="true" />
        </Menu.Button>
      </div>

      <Transition
        as={Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Menu.Items className="absolute right-0 z-50 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="py-1">
            <Menu.Item
            >
              {({ active }) => (
                <button
                  type="button"
                  onClick={dropChannel}
                  className={classNames(
                    active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                    'block w-full px-4 py-2 text-sm text-left'
                  )}
                >
                  Delete
                </button>
              )}
            </Menu.Item>
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  )
}
