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

export const Select: React.FC<Props> = ({name, value, className, onChange, options}) => (
  <select
    name={name}
    value={value}
    className={className}
    style={{width: 'auto'}}
    data-testid={'select-' + name}
    onChange={onChange}
  >
    {options.map((o: Option) => (
      <option key={o.value} value={o.value}>
        {o.label}
      </option>
    ))}
  </select>
);
