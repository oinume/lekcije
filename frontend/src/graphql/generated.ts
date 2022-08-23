import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, useMutation, UseQueryOptions, UseMutationOptions } from '@tanstack/react-query';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };

function fetcher<TData, TVariables>(client: GraphQLClient, query: string, variables?: TVariables, headers?: RequestInit['headers']) {
  return async (): Promise<TData> => client.request<TData, TVariables>(query, variables, headers);
}
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Empty = {
  __typename?: 'Empty';
  id: Scalars['ID'];
};

export type FollowingTeacher = {
  __typename?: 'FollowingTeacher';
  createdAt: Scalars['String'];
  id: Scalars['ID'];
  teacher: Teacher;
};

export type Mutation = {
  __typename?: 'Mutation';
  createEmpty?: Maybe<Empty>;
  updateViewer: User;
};


export type MutationUpdateViewerArgs = {
  input: UpdateViewerInput;
};

export type NotificationTimeSpan = {
  __typename?: 'NotificationTimeSpan';
  fromHour: Scalars['Int'];
  fromMinute: Scalars['Int'];
  toHour: Scalars['Int'];
  toMinute: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  empty?: Maybe<Empty>;
  followingTeachers: Array<FollowingTeacher>;
  viewer: User;
};

export type Teacher = {
  __typename?: 'Teacher';
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type UpdateViewerInput = {
  email?: InputMaybe<Scalars['String']>;
};

export type User = {
  __typename?: 'User';
  email: Scalars['String'];
  followingTeachers: Array<FollowingTeacher>;
  id: Scalars['ID'];
  notificationTimeSpans: Array<NotificationTimeSpan>;
  showTutorial: Scalars['Boolean'];
};

export type GetViewerQueryVariables = Exact<{ [key: string]: never; }>;


export type GetViewerQuery = { __typename?: 'Query', viewer: { __typename?: 'User', id: string, email: string, showTutorial: boolean } };

export type GetViewerWithNotificationTimeSpansQueryVariables = Exact<{ [key: string]: never; }>;


export type GetViewerWithNotificationTimeSpansQuery = { __typename?: 'Query', viewer: { __typename?: 'User', id: string, email: string, showTutorial: boolean, notificationTimeSpans: Array<{ __typename?: 'NotificationTimeSpan', fromHour: number, fromMinute: number, toHour: number, toMinute: number }> } };

export type GetViewerWithFollowingTeachersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetViewerWithFollowingTeachersQuery = { __typename?: 'Query', viewer: { __typename?: 'User', id: string, email: string, showTutorial: boolean, followingTeachers: Array<{ __typename?: 'FollowingTeacher', teacher: { __typename?: 'Teacher', id: string, name: string } }> } };

export type UpdateViewerMutationVariables = Exact<{
  input: UpdateViewerInput;
}>;


export type UpdateViewerMutation = { __typename?: 'Mutation', updateViewer: { __typename?: 'User', id: string, email: string } };


export const GetViewerDocument = `
    query getViewer {
  viewer {
    id
    email
    showTutorial
  }
}
    `;
export const useGetViewerQuery = <
      TData = GetViewerQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetViewerQueryVariables,
      options?: UseQueryOptions<GetViewerQuery, TError, TData>,
      headers?: RequestInit['headers']
    ) =>
    useQuery<GetViewerQuery, TError, TData>(
      variables === undefined ? ['getViewer'] : ['getViewer', variables],
      fetcher<GetViewerQuery, GetViewerQueryVariables>(client, GetViewerDocument, variables, headers),
      options
    );

useGetViewerQuery.getKey = (variables?: GetViewerQueryVariables) => variables === undefined ? ['getViewer'] : ['getViewer', variables];
;

export const GetViewerWithNotificationTimeSpansDocument = `
    query getViewerWithNotificationTimeSpans {
  viewer {
    id
    email
    notificationTimeSpans {
      fromHour
      fromMinute
      toHour
      toMinute
    }
    showTutorial
  }
}
    `;
export const useGetViewerWithNotificationTimeSpansQuery = <
      TData = GetViewerWithNotificationTimeSpansQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetViewerWithNotificationTimeSpansQueryVariables,
      options?: UseQueryOptions<GetViewerWithNotificationTimeSpansQuery, TError, TData>,
      headers?: RequestInit['headers']
    ) =>
    useQuery<GetViewerWithNotificationTimeSpansQuery, TError, TData>(
      variables === undefined ? ['getViewerWithNotificationTimeSpans'] : ['getViewerWithNotificationTimeSpans', variables],
      fetcher<GetViewerWithNotificationTimeSpansQuery, GetViewerWithNotificationTimeSpansQueryVariables>(client, GetViewerWithNotificationTimeSpansDocument, variables, headers),
      options
    );

useGetViewerWithNotificationTimeSpansQuery.getKey = (variables?: GetViewerWithNotificationTimeSpansQueryVariables) => variables === undefined ? ['getViewerWithNotificationTimeSpans'] : ['getViewerWithNotificationTimeSpans', variables];
;

export const GetViewerWithFollowingTeachersDocument = `
    query getViewerWithFollowingTeachers {
  viewer {
    id
    email
    followingTeachers {
      teacher {
        id
        name
      }
    }
    showTutorial
  }
}
    `;
export const useGetViewerWithFollowingTeachersQuery = <
      TData = GetViewerWithFollowingTeachersQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetViewerWithFollowingTeachersQueryVariables,
      options?: UseQueryOptions<GetViewerWithFollowingTeachersQuery, TError, TData>,
      headers?: RequestInit['headers']
    ) =>
    useQuery<GetViewerWithFollowingTeachersQuery, TError, TData>(
      variables === undefined ? ['getViewerWithFollowingTeachers'] : ['getViewerWithFollowingTeachers', variables],
      fetcher<GetViewerWithFollowingTeachersQuery, GetViewerWithFollowingTeachersQueryVariables>(client, GetViewerWithFollowingTeachersDocument, variables, headers),
      options
    );

useGetViewerWithFollowingTeachersQuery.getKey = (variables?: GetViewerWithFollowingTeachersQueryVariables) => variables === undefined ? ['getViewerWithFollowingTeachers'] : ['getViewerWithFollowingTeachers', variables];
;

export const UpdateViewerDocument = `
    mutation updateViewer($input: UpdateViewerInput!) {
  updateViewer(input: $input) {
    id
    email
  }
}
    `;
export const useUpdateViewerMutation = <
      TError = unknown,
      TContext = unknown
    >(
      client: GraphQLClient,
      options?: UseMutationOptions<UpdateViewerMutation, TError, UpdateViewerMutationVariables, TContext>,
      headers?: RequestInit['headers']
    ) =>
    useMutation<UpdateViewerMutation, TError, UpdateViewerMutationVariables, TContext>(
      ['updateViewer'],
      (variables?: UpdateViewerMutationVariables) => fetcher<UpdateViewerMutation, UpdateViewerMutationVariables>(client, UpdateViewerDocument, variables, headers)(),
      options
    );