import type {UseQueryOptions} from '@tanstack/react-query';

export const queryKeyMe = '/twirp/api.v1.Me/GetMe';
export const queryKeyFollowingTeachers = '/twirp/api.v1.Me/ListFollowingTeachers';
export const defaultUseQueryOptions: Omit<UseQueryOptions, any> = {
  refetchOnMount: 'always',
  retry: 0,
};
