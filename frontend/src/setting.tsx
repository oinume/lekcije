import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {ReactQueryDevtools} from '@tanstack/react-query-devtools';
import {SettingPage} from './pages/SettingPage';

const queryClient = new QueryClient();

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <SettingPage/>
    <ReactQueryDevtools initialIsOpen={false}/>
  </QueryClientProvider>,
  document.querySelector('#root'),
);
