import {useQuery} from 'react-query';
import {NotificationTimeSpan} from '../models/NotificatonTimeSpan';
import {twirpRequest} from '../http/twirp';
import {queryKeyMe} from './common';

type GetMeRequest = Record<string, unknown>;

type GetMeResponse = {
  email: string;
  notificationTimeSpans: NotificationTimeSpan[];
};

export const useGetMe = (
  request: GetMeRequest,
) => useQuery<GetMeResponse, Error>(
  queryKeyMe,
  async () => {
    const response = await twirpRequest('/twirp/api.v1.Me/GetMe', JSON.stringify(request));
    const data = await response.json() as GetMeResponse;
    // Console.log('----- data -----');
    // console.log(data);
    return data;
  },
  {
    retry: 0, // TODO: commonize option
  },
);
