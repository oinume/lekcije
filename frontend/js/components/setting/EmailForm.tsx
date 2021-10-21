import React, {useState} from 'react';

type Props = {
  email: string;
  handleOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleUpdateEmail: (email: string) => void;
};

export const EmailForm: React.FC<Props> = ({email, handleOnChange, handleUpdateEmail}) => {
  const [buttonEnabled, setButtonEnabled] = useState<boolean>(false);

  return (
    <form
      onSubmit={e => {
        e.preventDefault();
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
          onChange={e => {
            e.preventDefault();
            setButtonEnabled(e.currentTarget.value !== '');
            handleOnChange(e);
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
