import cookie from 'cookie';
import fetch from 'cross-fetch';

export const sendRequest = async (path: string, body: string) => {
  const headers: { [key: string]: string } = {
    'Content-Type': 'application/json',
  };
  const cookies = cookie.parse(document.cookie);
  if (cookies['apiToken']) {
    headers['Authorization'] = 'Bearer ' + cookies['apiToken'];
  }
  return fetch(path, {
    body: body,
    method: 'POST',
    headers: headers,
  });
};
