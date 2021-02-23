import React, { useState } from 'react';

type Props = {
  email: string;
  handleOnChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handleUpdateEmail: (email: string) => void;
};

export const EmailForm: React.FC<Props> = ({
  email,
  handleOnChange,
  handleUpdateEmail,
}) => {
  const [buttonEnabled, setButtonEnabled] = useState<boolean>(false);
  return (
    <form className="form-horizontal">
      <div className="form-group">
        <div className="col-sm-3">
          <label htmlFor="email" className="control-label">
            Email
          </label>
        </div>
        <div className="col-sm-7">
          <input
            type="email"
            className="form-control"
            name="email"
            id="email"
            placeholder="Email"
            required
            autoFocus
            autoComplete="on"
            value={email}
            onChange={(e) => {
              setButtonEnabled(e.currentTarget.value !== '');
              handleOnChange(e);
            }}
          />
        </div>
      </div>
      <div className="form-group">
        <div className="col-sm-offset-2 col-sm-8">
          <button
            type="button"
            disabled={!buttonEnabled}
            className="btn btn-primary"
            onClick={() => handleUpdateEmail(email)}
          >
            変更
          </button>
        </div>
      </div>
    </form>
  );
};
