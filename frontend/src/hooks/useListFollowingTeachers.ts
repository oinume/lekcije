import {useQuery} from '@tanstack/react-query';
import type {TwirpError} from '../http/twirp';
import {twirpRequest} from '../http/twirp';
import type {Teacher} from '../models/Teacher';
import {defaultUseQueryOptions, queryKeyFollowingTeachers} from './common';

type ListFollowingTeachersRequest = Record<string, unknown>;

type ListFollowingTeachersResponse = {
  teachers: Teacher[];
};

export const useListFollowingTeachers = (
  request: ListFollowingTeachersRequest,
) => useQuery<ListFollowingTeachersResponse, TwirpError>(
  [queryKeyFollowingTeachers],
  async () => {
    const response = await twirpRequest('/twirp/api.v1.Me/ListFollowingTeachers', JSON.stringify(request));
    const data = await response.json() as ListFollowingTeachersResponse;
    /*
    console.log('----- data -----');
    console.log(data);
     */
    return data;
  },
  defaultUseQueryOptions,
);
