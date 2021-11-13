import React, {useState} from 'react';

type Props = {
  email: string;
  handleOnChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handleUpdateEmail: (email: string) => void;
};

export const EmailForm: React.FC<Props> = ({email, handleOnChange, handleUpdateEmail}) => {
  const [buttonEnabled, setButtonEnabled] = useState<boolean>(false);

  return (
    <form
      onSubmit={event => {
        event.preventDefault();
        handleUpdateEmail(email);
      }}
    >
      <h5>通知先メールアドレス</h5>
      <div className="mb-3">
        <input
          required
          autoFocus
          type="email"
          className="form-control"
          name="email"
          id="email"
          placeholder="Email"
          autoComplete="on"
          value={email}
          onChange={event => {
            event.preventDefault();
            setButtonEnabled(event.currentTarget.value !== '');
            handleOnChange(event);
          }}
        />
      </div>
      <button
        type="button"
        disabled={!buttonEnabled}
        className="btn btn-primary"
        onClick={() => {
          handleUpdateEmail(email);
        }}
      >
        変更
      </button>
    </form>
  );
};
