import {DefaultOptions} from '@tanstack/react-query';

export const defaultQueryClientOptions: Omit<DefaultOptions, any> = {
  mutations: {
    retry: 0,
  },
  queries: {
    refetchOnMount: 'always',
    retry: 0,
  },
};
