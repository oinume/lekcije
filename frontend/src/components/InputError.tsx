import React from 'react';

type Props = {
  readonly message: string | undefined;
};

export const InputError: React.FC<Props> = ({message}) => <div className="invalid-feedback">{message}</div>;
