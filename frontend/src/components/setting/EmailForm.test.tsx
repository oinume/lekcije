import React from 'react';
import {render, screen} from '@testing-library/react';
import '@testing-library/jest-dom';
import {EmailForm} from './EmailForm';

test('<EmailForm>', () => {
  const handleOnChange = (_: React.ChangeEvent<HTMLInputElement>) => { /* Nop */ };
  const handleUpdateEmail = (_: string) => true;
  render(<EmailForm email="oinume@gmail.com" handleOnChange={handleOnChange} handleUpdateEmail={handleUpdateEmail}/>);
  expect(screen.getByDisplayValue('oinume@gmail.com')).toBeInTheDocument();
  // TODO: write more assertions
});
