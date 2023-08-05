import React, {useState} from 'react';
import {type SubmitHandler, useForm} from 'react-hook-form';
import {InputError} from '../InputError';

type Props = {
  readonly email: string;
  readonly handleOnChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  readonly handleUpdateEmail: (email: string) => boolean;
};

type FormValues = {
  email: string;
};

const formValidationRules = {
  required: {value: true, message: '入力してください'},
  maxLength: {value: 100, message: '100文字以内で入力してください'},
};

export const EmailForm: React.FC<Props> = ({email, handleOnChange, handleUpdateEmail}) => {
  const {
    register,
    handleSubmit,
    formState: {errors, isDirty, isValid},
  } = useForm<FormValues>({mode: 'onChange'});
  const emailField = register('email', {
    ...formValidationRules,
    onChange(event: React.ChangeEvent<HTMLInputElement>) {
      event.preventDefault();
      handleOnChange(event);
    },
  });

  const onSubmit: SubmitHandler<FormValues> = data => {
    const success = handleUpdateEmail(email);
  };

  return (
    <form
      className="needs-validation"
      onSubmit={handleSubmit(onSubmit)}
    >
      <h5>通知先メールアドレス</h5>
      <div className="mb-3">
        <input
          required
          autoFocus
          type="email"
          className={`form-control ${errors.email ? 'is-invalid' : ''}`}
          id="email"
          placeholder="Email"
          autoComplete="on"
          value={email}
          {...emailField}
        />
        { errors.email && <InputError message={errors.email.message}/> }
      </div>
      <button
        type="submit"
        disabled={!isDirty || !isValid}
        className="btn btn-primary"
      >
        変更
      </button>
    </form>
  );
};
