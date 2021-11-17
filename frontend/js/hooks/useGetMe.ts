import {useQuery} from 'react-query';
import {sendRequest} from '../http/fetch';
import {NotificationTimeSpan} from '../components/setting/NotificationTimeSpanForm';
import {TwirpError} from './TwirpError';

type GetMeRequest = Record<string, unknown>;

type GetMeResponse = {
  email: string;
  notificationTimeSpans: NotificationTimeSpan[];
};

const path = '/twirp/api.v1.Me/GetMe';

export const useGetMe = (
  request: GetMeRequest,
) => {
  return useQuery<GetMeResponse, Error>(
    path,
    async () => {
      const response = await sendRequest(path, JSON.stringify(request));
      if (!response.ok) {
        // TODO: error
        const error = TwirpError.fromJson(await response.json());
        throw error;
      }

      const data = await response.json() as GetMeResponse;
      // Console.log('----- data -----');
      // console.log(data);
      return data;
    },
    {
      retry: 0, // TODO: commonize option
    },
  );
};
