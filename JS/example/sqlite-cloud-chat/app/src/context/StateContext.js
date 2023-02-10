import React, { useState, useRef, createContext } from "react";

const StateContext = createContext({});

const StateProvider = ({ children }) => {

  //queue used to store all received mgs from PubSub
  const [chsMap, setChsMap] = useState(new Map());
  const chsMapRef = useRef(chsMap);

  return (
    <StateContext.Provider
      value={{
        chsMap, chsMapRef, setChsMap
      }}
    >
      {children}
    </StateContext.Provider>
  )
}

export { StateContext, StateProvider }
