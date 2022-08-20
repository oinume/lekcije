import React from 'react';

type Props = {
  kind: string;
  message: string;
  isVisible: boolean;
  handleCloseAlert: () => void;
};

export const ToggleAlert: React.FC<Props> = ({kind, message, isVisible, handleCloseAlert}) => {
  if (isVisible) {
    return (
      <div className={'alert alert-dismissible alert-' + kind} role="alert">
        <button
          type="button"
          className="btn-close"
          data-bs-dismiss="alert"
          aria-label="Close"
          onClick={() => {
            handleCloseAlert();
          }}
        />
        {message}
      </div>
    );
  }

  return <div/>;
};
