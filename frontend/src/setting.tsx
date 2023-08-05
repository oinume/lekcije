import React from 'react';
import {createRoot} from 'react-dom/client.js';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {ReactQueryDevtools} from '@tanstack/react-query-devtools';
import {SettingPage} from './pages/SettingPage';
import {defaultQueryClientOptions} from './http/query';

const queryClient = new QueryClient({
  defaultOptions: defaultQueryClientOptions,
});
const container = document.querySelector('#root');
const root = createRoot(container!);
root.render(
  <QueryClientProvider client={queryClient}>
    <SettingPage/>
    <ReactQueryDevtools initialIsOpen={false}/>
  </QueryClientProvider>,
);
