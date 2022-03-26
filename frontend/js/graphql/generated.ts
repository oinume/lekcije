import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import { useQuery, UseQueryOptions } from 'react-query';
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

export type User = {
  __typename?: 'User';
  email: Scalars['String'];
  followingTeachers: Array<FollowingTeacher>;
  id: Scalars['ID'];
};

export type GetViewerQueryVariables = Exact<{ [key: string]: never; }>;


export type GetViewerQuery = { __typename?: 'Query', viewer: { __typename?: 'User', id: string, email: string } };

export type GetViewerFollowingTeachersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetViewerFollowingTeachersQuery = { __typename?: 'Query', viewer: { __typename?: 'User', followingTeachers: Array<{ __typename?: 'FollowingTeacher', teacher: { __typename?: 'Teacher', id: string, name: string } }> } };


export const GetViewerDocument = `
    query getViewer {
  viewer {
    id
    email
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

export const GetViewerFollowingTeachersDocument = `
    query getViewerFollowingTeachers {
  viewer {
    followingTeachers {
      teacher {
        id
        name
      }
    }
  }
}
    `;
export const useGetViewerFollowingTeachersQuery = <
      TData = GetViewerFollowingTeachersQuery,
      TError = unknown
    >(
      client: GraphQLClient,
      variables?: GetViewerFollowingTeachersQueryVariables,
      options?: UseQueryOptions<GetViewerFollowingTeachersQuery, TError, TData>,
      headers?: RequestInit['headers']
    ) =>
    useQuery<GetViewerFollowingTeachersQuery, TError, TData>(
      variables === undefined ? ['getViewerFollowingTeachers'] : ['getViewerFollowingTeachers', variables],
      fetcher<GetViewerFollowingTeachersQuery, GetViewerFollowingTeachersQueryVariables>(client, GetViewerFollowingTeachersDocument, variables, headers),
      options
    );

useGetViewerFollowingTeachersQuery.getKey = (variables?: GetViewerFollowingTeachersQueryVariables) => variables === undefined ? ['getViewerFollowingTeachers'] : ['getViewerFollowingTeachers', variables];
;
