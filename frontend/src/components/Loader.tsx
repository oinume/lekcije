import type {CSSProperties} from 'react';
import React from 'react';
import {ClipLoader} from 'react-spinners';

type Props = {
  readonly isLoading: boolean;
  readonly message?: string;
  readonly css?: CSSProperties;
  readonly size?: number;
};

// TODO: Use bootstrap loader: https://getbootstrap.com/docs/4.4/components/spinners/
export const Loader = ({isLoading, message, css, size}: Props) => {
  if (message === undefined) {
    message = 'Loading data ...';
  }

  if (css === undefined) {
    css = {
      background: 'rgba(255, 255, 255, 0)',
    };
  }

  if (size === undefined) {
    size = 50;
  }

  return isLoading ? (
    <div className="overlay-content">
      <div className="wrapper" style={{textAlign: 'center'}}>
        <ClipLoader cssOverride={css} size={size} color="#123abc" loading={isLoading}/>
        <p>
          <span className="message">{message}</span>
        </p>
      </div>
    </div>
  ) : null;
};
