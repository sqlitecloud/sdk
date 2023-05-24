/*!
 * SQLite Cloud Node.js SDK v1.0.0
 * https://sqlitecloud.io/
 *
 * Copyright 2023, SQLite Cloud
 * Released under the MIT licence.
 */

(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define([], factory);
	else if(typeof exports === 'object')
		exports["SQLiteCloud"] = factory();
	else
		root["SQLiteCloud"] = factory();
})(this, () => {
return /******/ (() => { // webpackBootstrap
/******/ 	var __webpack_modules__ = ({

/***/ 763:
/***/ ((module, __unused_webpack_exports, __webpack_require__) => {

module.exports = __webpack_require__(774)["default"];


/***/ }),

/***/ 774:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

"use strict";
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "default": () => (/* binding */ SQLiteCloud)
/* harmony export */ });
const logThis = (id = "SQLiteCloud", msg) => {
  const dateObject = new Date();
  // adjust 0 before single digit date
  const date = (`0 ${dateObject.getDate()}`).slice(-2);
  // current month
  const month = (`0${dateObject.getMonth() + 1}`).slice(-2);
  // current year
  const year = dateObject.getFullYear();
  // current hours
  const hours = dateObject.getHours();
  // current minutes
  const minutes = dateObject.getMinutes();
  // current seconds
  const seconds = dateObject.getSeconds();
  // current milliseconds
  const milliseconds = dateObject.getMilliseconds()
  // prints date & time in YYYY/MM/DD HH:MM:SS format
  const prefix = `${year}/${month}/${date} ${hours}:${minutes}:${seconds}:${milliseconds}`;
  console.log(`!!!!!!!!! ${id}: ${prefix} - ${msg}`);
}
/*
SQLiteCloud class
*/
class SQLiteCloud {
  /* PRIVATE PROPERTIES */
  /* 
  #debug_sdk 
  */
  #debug_sdk = false;
  /* CONSTRUCTOR */
  /*
  SQLiteCloud class constructor receives:
   - project ID (required)
   - api key (required)
   - debug
  */
  constructor(projectID, apikey, debug_sdk = false) {
    this.#debug_sdk = debug_sdk;
    this.url = `wss://web1.sqlitecloud.io:8443/api/v1/${projectID}/ws?apikey=${apikey}`;
  }
  /*
  method used to send SQL commands using fetch API
  */
  async sendCommands(commands) {
    if (this.#debug_sdk) logThis("", "sending commands: " + commands);
    const data = {
      "statement": commands
    }
    const response = await fetch(this.url, {
      method: "POST",
      mode: "cors",
      cache: "no-cache",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json",
      },
      redirect: "follow",
      referrerPolicy: "no-referrer",
      body: JSON.stringify(data),
    });
    const jsonData = await response.json();
    return jsonData
  }
}

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module used 'module' so it can't be inlined
/******/ 	var __webpack_exports__ = __webpack_require__(763);
/******/ 	
/******/ 	return __webpack_exports__;
/******/ })()
;
});