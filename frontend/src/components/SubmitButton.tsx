import React from 'react';

type Props = {
  disabled: boolean;
  loading: boolean;
  buttonProps?: any;
  children?: React.ReactNode
};

export const SubmitButton: React.FC<Props> = ({disabled, loading, children}) => (
    <button
      type="submit"
      className="btn btn-primary"
      disabled={disabled}
    >
      { children ? children : <></> }
      { loading ?
        (<span
          className="spinner-border spinner-border-sm mx-1"
          role="status"
          aria-hidden="true"
        />) : <></> }
    </button>
);
