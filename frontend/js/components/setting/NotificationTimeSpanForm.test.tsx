import React from 'react';
import { NotificationTimeSpanForm } from './NotificationTimeSpanForm';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';

test('<NotificationTimeSpanForm>', () => {
  render(
    <NotificationTimeSpanForm
      timeSpans={[]}
      handleAdd={() => {}}
      handleDelete={() => {}}
      handleUpdate={() => {}}
      handleOnChange={() => {}}
    />
  );

  expect(screen.getByText('編集')).toBeInTheDocument();
  // TODO: write more assertions
});
