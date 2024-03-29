import React from 'react';

type Props = {
  readonly children: React.ReactNode;
};

export const PageTitle = ({children}: Props) => <h1 className="page-title">{children}</h1>;

