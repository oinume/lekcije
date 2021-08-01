import React, { useState } from 'react';

type Props = {
  kind: string;
  message: string;
};

export const Alert: React.FC<Props> = ({ kind, message }) => {
  const [visible, setVisible] = useState<boolean>(true);
  if (visible) {
    return (
      <div className={'alert alert-' + kind} role="alert">
        <button
          type="button"
          className="close"
          onClick={() => setVisible(false)}
        >
          &times;
        </button>
        {message}
      </div>
    );
  } else {
    return <div />;
  }
};
