import React from 'react';
import ReactDOM from 'react-dom';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {MePage} from './pages/MePage';

const queryClient = new QueryClient();

ReactDOM.render(
  <QueryClientProvider client={queryClient}>
    <MePage/>
  </QueryClientProvider>,
  document.querySelector('#root'),
);
