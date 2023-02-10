//core
import React, { Fragment, useEffect, useState } from 'react';
//utils
import { logThis } from '../js/utils';
//context
import { StateProvider } from '../context/StateContext';
//components
import Layout from './Layout';

const App = () => {
  if (process.env.DEBUG == 'true') logThis('App: ON RENDER');
  return (
    <Fragment>
      <StateProvider>
        <Layout />
      </StateProvider>
    </Fragment>
  );
}


export default App;