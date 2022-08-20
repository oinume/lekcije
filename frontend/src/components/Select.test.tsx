import React from 'react';
import {render, screen} from '@testing-library/react';
import {Select} from './Select';
import '@testing-library/jest-dom';

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
