import cookie from 'cookie';
import {GraphQLClient} from 'graphql-request';

export const createGraphQLClient = (path?: string, token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (path === undefined) {
    path = '/graphql'
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
    headers: headers,
  })
}

// TODO: Define error class correctly
export class GraphQLError extends Error {}
