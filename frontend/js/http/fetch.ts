import cookie from 'cookie';
import fetch, {Response} from 'cross-fetch';

export class HttpError extends Error {
  constructor(public message: string, public response: Response) {
    super(message);
  }
}

export const sendRequest = async (path: string, body: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  const cookies = cookie.parse(document.cookie);
  if (cookies.apiToken) {
    headers.Authorization = 'Bearer ' + cookies.apiToken;
  }

  const response = await fetch(path, {
    body,
    method: 'POST',
    headers,
  });
  if (response.status >= 400) {
    throw new HttpError('HTTP request failed on ' + path, response);
  }

  return response;
};
