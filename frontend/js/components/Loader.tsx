import React from 'react';
import {ClipLoader} from 'react-spinners';

type Props = {
  isLoading: boolean;
  message?: string;
  css?: string;
  size?: number;
};

export const Loader = ({isLoading, message, css, size}: Props) => {
  if (message === undefined) {
    message = 'Loading data ...';
  }

  if (css === undefined) {
    css = 'background: rgba(255, 255, 255, 0)';
  }

  if (size === undefined) {
    size = 50;
  }

  return isLoading ? (
    <div className="overlay-content">
      <div className="wrapper" style={{textAlign: 'center'}}>
        <ClipLoader css={css} size={size} color="#123abc" loading={isLoading}/>
        <p>
          <span className="message">{message}</span>
        </p>
      </div>
    </div>
  ) : null;
};
