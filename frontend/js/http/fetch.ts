import cookie from 'cookie';
import fetch, { Response } from 'cross-fetch';

class HttpError extends Error {
  constructor(message: string, public response: Response) {
    super(message);
  }
}

export const sendRequest = async (path: string, body: string) => {
  const headers: { [key: string]: string } = {
    'Content-Type': 'application/json',
  };
  const cookies = cookie.parse(document.cookie);
  if (cookies['apiToken']) {
    headers['Authorization'] = 'Bearer ' + cookies['apiToken'];
  }
  const res = await fetch(path, {
    body: body,
    method: 'POST',
    headers: headers,
  });
  if (res.status >= 400) {
    throw new HttpError('HTTP request failed on ' + path, res);
  }
  return res;
};
