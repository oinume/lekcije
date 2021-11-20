import {useQuery} from 'react-query';
import {NotificationTimeSpan} from '../components/setting/NotificationTimeSpanForm';
import {twirpRequest} from "../http/twirp";

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
      const response = await twirpRequest(path, JSON.stringify(request));
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
