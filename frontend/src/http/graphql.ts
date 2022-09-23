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

type ErrorExtensions = {
  code?: string;
  localizedMessage?: string;
};

type ErrorInfo = {
  path: string;
  message: string;
  extensions: ErrorExtensions;
};

export type GraphQLError = {
  response: {errors: ErrorInfo[]};
};

export const toMessage = (error: GraphQLError, defaultMessage?: string): string => {
  const errs = error.response.errors;
  if (errs.length === 0) {
    return '';
  }

  return `${errs[0].extensions.localizedMessage ?? defaultMessage ?? ''} (${errs[0].extensions.code ?? ''})`;
};
