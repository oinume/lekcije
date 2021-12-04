import {useQuery} from 'react-query';
import {NotificationTimeSpan} from '../models/NotificatonTimeSpan';
import {TwirpError, twirpRequest} from '../http/twirp';
import {User} from '../models/User';
import {defaultUseQueryOptions, queryKeyMe} from './common';

type GetMeRequest = Record<string, unknown>;

type GetMeResponse = {
  email: string;
  notificationTimeSpans: NotificationTimeSpan[];
  user: User;
  showTutorial: boolean;
};

export const useGetMe = (
  request: GetMeRequest,
) => useQuery<GetMeResponse, TwirpError>(
  queryKeyMe,
  async () => {
    const response = await twirpRequest('/twirp/api.v1.Me/GetMe', JSON.stringify(request));
    const data = await response.json() as GetMeResponse;
    // Console.log('----- data -----');
    // console.log(data);
    return data;
  },
  defaultUseQueryOptions,
);
