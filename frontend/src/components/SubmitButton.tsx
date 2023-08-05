import React from 'react';

type Props = {
  readonly disabled: boolean; // eslint-disable-line react/boolean-prop-naming
  readonly loading: boolean; // eslint-disable-line react/boolean-prop-naming
  readonly children?: React.ReactNode;
};

export const SubmitButton: React.FC<Props> = ({disabled, loading, children}) => (
  <button
    type="submit"
    className="btn btn-primary"
    disabled={disabled}
  >
    { children ?? children }
    { loading ? <span className="spinner-border spinner-border-sm mx-1" role="status" aria-hidden="true"/> : undefined }
  </button>
);
