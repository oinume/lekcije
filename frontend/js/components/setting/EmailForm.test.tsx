import React from 'react';
import { EmailForm } from './EmailForm';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';

test('<EmailForm>', () => {
  const handleOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {};
  const handleUpdateEmail = (email: string) => {};
  render(<EmailForm email="oinume@gmail.com" handleOnChange={handleOnChange} handleUpdateEmail={handleUpdateEmail} />);
  // expect(wrapper.find('button')).toHaveLength(1);
  // expect(wrapper.find('[name="email"]')).toHaveLength(1);
  expect(screen.getByDisplayValue('oinume@gmail.com')).toBeInTheDocument();
});
