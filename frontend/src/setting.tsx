import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {SettingPage} from './pages/SettingPage';
import {ReactQueryDevtools} from '@tanstack/react-query-devtools';

const queryClient = new QueryClient();

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <SettingPage/>
    <ReactQueryDevtools initialIsOpen={false} />
  </QueryClientProvider>,
  document.querySelector('#root'),
);
