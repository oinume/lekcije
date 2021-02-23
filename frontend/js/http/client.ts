import axios from 'axios';
import cookie from 'cookie';

export const createHttpClient = () => {
  const cookies = cookie.parse(document.cookie);
  const headers: { [key: string]: string } = {};
  if (cookies['apiToken']) {
    headers['Grpc-Metadata-Api-Token'] = cookies['apiToken'];
    headers['X-Api-Token'] = cookies['apiToken'];
  }
  return axios.create({
    //baseURL: 'https://some-domain.com/api/',
    timeout: 3000,
    headers: headers,
  });
};
