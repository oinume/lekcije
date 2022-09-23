import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {ReactQueryDevtools} from '@tanstack/react-query-devtools';
import {MePage} from './pages/MePage';
import {defaultQueryClientOptions} from './http/query';

const queryClient = new QueryClient({
  defaultOptions: defaultQueryClientOptions,
});

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <MePage/>
    <ReactQueryDevtools initialIsOpen={false}/>
  </QueryClientProvider>,
  document.querySelector('#root'),
);
