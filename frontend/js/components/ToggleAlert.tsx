import React from 'react';

type Props = {
  kind: string;
  message: string;
  visible: boolean;
  handleCloseAlert: () => void;
};

export const ToggleAlert: React.FC<Props> = ({
  kind,
  message,
  visible,
  handleCloseAlert,
}) => {
  if (visible) {
    return (
      <div className={'alert alert-' + kind} role="alert">
        <button
          type="button"
          className="close"
          onClick={() => handleCloseAlert()}
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
