import React from 'react';

export type Option = {
  value: string;
  label: string;
};

type Props = {
  readonly name: string;
  readonly value: string;
  readonly className: string;
  readonly onChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
  readonly options: Option[];
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
