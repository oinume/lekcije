import cookie from "cookie";
import fetch from "cross-fetch";

export const twirpRequest = async (path: string, body: string, token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token === undefined) {
    const cookies = cookie.parse(document.cookie);
    if (cookies.apiToken) {
      token = cookies.apiToken;
    }
  }
  headers.Authorization = 'Bearer ' + token;

  const response = await fetch(path, {
    body,
    method: 'POST',
    headers,
  });
  if (response.status >= 400) {
    const j = await response.text();
    //console.log(j);
    const error = TwirpError.fromJson(j);
    throw error;
  }

  return response;
};

export class TwirpError extends Error {
  static fromJson(json: string): TwirpError {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const obj = JSON.parse(json);
    return new TwirpError(obj.code, obj.msg);
  }

  get message(): string {
    return `${this.code}:${this.msg}`;
  }

  constructor(code: string, msg: string) {
    super(`code:${code}, message:${msg}`);
    this.code = code;
    this.msg = msg;
  }

  code: string;
  msg: string;
}
