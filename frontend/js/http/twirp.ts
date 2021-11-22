import cookie from 'cookie';
import fetch from 'cross-fetch';

export const twirpRequest = async (path: string, body: string, token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token === undefined) {
    const cookies = cookie.parse(document.cookie);
    if (cookies.apiToken) {
      headers.Authorization = `Bearer ${cookies.apiToken}`;
    }
  } else {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(path, {
    body,
    method: 'POST',
    headers,
  });
  if (response.status >= 400) {
    const body = await response.text();
    // Console.log(j);
    throw TwirpError.fromJson(body);
  }

  return response;
};

export class TwirpError extends Error {
  static fromJson(json: string): TwirpError {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const object = JSON.parse(json);
    return new TwirpError(object.code, object.msg);
  }

  get message(): string {
    return `${this.code}:${this.msg}`;
  }

  public isInternal(): boolean {
    return this.code === 'Internal';
  }

  constructor(code: string, message: string) {
    super(`code:${code}, message:${message}`);
    this.code = code;
    this.msg = message;
  }

  code: string;
  msg: string;
}
