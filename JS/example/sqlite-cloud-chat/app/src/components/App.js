//core
import React, { Fragment, useEffect, useState } from 'react';
//utils
import { logThis } from '../js/utils';
//context
import { StateProvider } from '../context/StateContext';
//components
import Layout from './Layout';

const App = () => {
  if (process.env.DEBUG == 'true') logThis('App: ON RENDER TEST');
  console.log(process.env);
  console.log(process.env.PROJECT_ID);
  console.log(process.env.API_KEY);
  return (
    <Fragment>
      <StateProvider>
        <Layout />
      </StateProvider>
    </Fragment>
  );
}


export default App;