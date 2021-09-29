import React from 'react';

type Props = {
  kind: string;
  message: string;
  visible: boolean;
  handleCloseAlert: () => void;
};

export const ToggleAlert: React.FC<Props> = ({ kind, message, visible, handleCloseAlert }) => {
  if (visible) {
    return (
      <div className={`alert alert-dismissible alert-${kind}`} role="alert">
        <button
          type="button"
          className="btn-close"
          data-bs-dismiss="alert"
          aria-label="Close"
          onClick={() => handleCloseAlert()}
        />
        {message}
      </div>
    );
  }
  return <div />;
};
