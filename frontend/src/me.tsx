import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {ReactQueryDevtools} from '@tanstack/react-query-devtools';
import {MePage} from './pages/MePage';

const queryClient = new QueryClient();

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <MePage/>
    <ReactQueryDevtools initialIsOpen={false}/>
  </QueryClientProvider>,
  document.querySelector('#root'),
);
