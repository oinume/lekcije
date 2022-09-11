import type {UseMutationResult} from '@tanstack/react-query';
import React from 'react';
import type {TwirpError} from '../http/twirp';
import {Alert} from './Alert';

type Props = {
  result: UseMutationResult<any, any, any, any>;
  name: string;
};

export const UseMutationResultAlert = ({
  result,
  name,
}: Props) => {
  const error: TwirpError = result.error as TwirpError;
  // TODO: error handlingをまともにする
  switch (result.status) {
    case 'success':
      return <Alert kind="success" message={name + 'を更新しました！'}/>;
    case 'error':
      return <Alert kind="danger" message={name + `の更新に失敗しました。(${error.message})`}/>;
    default:
      return <div/>;
  }
};
