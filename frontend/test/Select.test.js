'use strict';

import React from 'react';
//import renderer from 'react-test-renderer';
import Select from "../js/components/Select";
import {shallow} from 'enzyme';
import { configure } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

configure({ adapter: new Adapter() });

test('Select initial state', () => {
  const options = [
    {value: 'japan', label: 'Japan'},
    {value: 'china', label: 'China'},
  ];
  const wrapper = shallow(
    <Select
      name="country"
      value=""
      onChange={() => {}}
      options={options}
    />
  );

  expect(wrapper.find('[name="country"]')).toHaveLength(1);
  expect(wrapper.find('[name="country"]').children()).toHaveLength(2);
});
