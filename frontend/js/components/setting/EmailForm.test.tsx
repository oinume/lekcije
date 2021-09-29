import React from 'react';
import { render, screen } from '@testing-library/react';
import { EmailForm } from './EmailForm';
import '@testing-library/jest-dom';

test('<EmailForm>', () => {
  const handleOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {};
  const handleUpdateEmail = (email: string) => {};
  render(<EmailForm email="oinume@gmail.com" handleOnChange={handleOnChange} handleUpdateEmail={handleUpdateEmail} />);
  expect(screen.getByDisplayValue('oinume@gmail.com')).toBeInTheDocument();
  // TODO: write more assertions
});
