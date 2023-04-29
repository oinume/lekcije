import React from 'react';
import {render, screen} from '@testing-library/react';
import '@testing-library/jest-dom';
import {Select} from './Select';

test('<Select> initial state', () => {
  const options = [
    {value: 'japan', label: 'Japan'},
    {value: 'china', label: 'China'},
  ];
  render(
    <Select
      name="country"
      value=""
      className=""
      options={options}
      onChange={(_: React.ChangeEvent<HTMLSelectElement>) => {
        // Nop
      }}
    />,
  );

  expect(screen.getByDisplayValue('Japan')).toBeInTheDocument();
  // TODO: write more assertions
});
