/**
 * @jest-environment jsdom
 */
import React from 'react';
import { Select } from '../js/components/Select';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';

test('Select initial state', () => {
  const options = [
    { value: 'japan', label: 'Japan' },
    { value: 'china', label: 'China' },
  ];
  render(
    <Select
      name="country"
      value=""
      className=""
      onChange={(event: React.ChangeEvent<HTMLSelectElement>) => {}}
      options={options}
    />
  );

  expect(screen.getByTestId('select')).toBeInTheDocument();
});
