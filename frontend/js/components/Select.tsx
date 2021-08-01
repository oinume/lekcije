import React from 'react';

export type Option = {
  value: string;
  label: string;
};

type Props = {
  name: string;
  value: string;
  className: string;
  onChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
  options: Option[];
};

export const Select: React.FC<Props> = ({ name, value, className, onChange, options }) => {
  return (
    <select
      name={name}
      value={value}
      className={className}
      onChange={onChange}
      style={{ width: 'auto' }}
      data-testid={'select'}
    >
      {options.map((o: Option) => {
        return (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        );
      })}
    </select>
  );
};
