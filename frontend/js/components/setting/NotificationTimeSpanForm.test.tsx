import React from 'react';
import {render, screen} from '@testing-library/react';
import {NotificationTimeSpanForm} from './NotificationTimeSpanForm';
import '@testing-library/jest-dom';

test('<NotificationTimeSpanForm>', () => {
  render(
    <NotificationTimeSpanForm
      timeSpans={[]}
      handleAdd={() => { /* Nop */ }}
      handleDelete={() => { /* Nop */ }}
      handleUpdate={() => { /* Nop */ }}
      handleOnChange={() => { /* Nop */ }}
    />,
  );

  expect(screen.getByText('編集')).toBeInTheDocument();
  // TODO: write more assertions
});
