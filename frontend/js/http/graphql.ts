import cookie from 'cookie';
import {GraphQLClient} from 'graphql-request';

export const createGraphQLClient = (path?: string, token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (path === undefined) {
    path = '/graphql';
  }

  if (token === undefined) {
    const cookies = cookie.parse(document.cookie);
    if (cookies.apiToken) {
      headers.Authorization = `Bearer ${cookies.apiToken}`;
    }
  } else {
    headers.Authorization = `Bearer ${token}`;
  }

  return new GraphQLClient(path, {
    headers,
  });
};

// TODO: Define error class correctly
export class GraphQLError extends Error {
  static fromJson(json: string): GraphQLError {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const object = JSON.parse(json);
    return new GraphQLError(object.path, object.message);
  }

  path: string;
  message: string;

  get string(): string {
    return `path=${this.path} message=${this.message}`;
  }

  constructor(path: string, message: string) {
    super(`path:${path}, message:${message}`);
    this.path = path;
    this.message = message;
  }
}
