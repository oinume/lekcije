//test.skip('<NotificationTimeSpanForm>', () => {})

import React from 'react';
import {configure, shallow} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { NotificationTimeSpanFormFC } from '../js/components/setting/NotificationTimeSpanForm';

configure({adapter: new Adapter()});

test('<NotificationTimeSpanFormFC>', () => {
  const wrapper = shallow(
    <NotificationTimeSpanFormFC
      editable={false}
      timeSpans={[]}
      handleAdd={() => {}}
      handleDelete={() => {}}
      handleUpdate={() => {}}
      handleOnChange={() => {}}
      handleSetEditable={() => {}}
    />
  );
  expect(wrapper.find('button')).toHaveLength(1);
  // expect(wrapper.find('[name="email"]')).toHaveLength(1);
});
