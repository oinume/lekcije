import React from 'react';
import { Select } from '../js/components/Select';
import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

configure({ adapter: new Adapter() });

test('Select initial state', () => {
  const options = [
    { value: 'japan', label: 'Japan' },
    { value: 'china', label: 'China' },
  ];
  const wrapper = shallow(
    <Select
      name="country"
      value=""
      className=""
      onChange={(event: React.ChangeEvent<HTMLSelectElement>) => {}}
      options={options}
    />
  );

  expect(wrapper.find('[name="country"]')).toHaveLength(1);
  expect(wrapper.find('[name="country"]').children()).toHaveLength(2);
});
