import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from 'react-query';
import {SettingPage} from './components/setting/SettingPage';

const queryClient = new QueryClient();

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <SettingPage/>
  </QueryClientProvider>,
  document.querySelector('#root'),
);
