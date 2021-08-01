import React from 'react';
import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { NotificationTimeSpanForm } from '../js/components/setting/NotificationTimeSpanForm';

configure({ adapter: new Adapter() });

test('<NotificationTimeSpanForm>', () => {
  const wrapper = shallow(
    <NotificationTimeSpanForm
      timeSpans={[]}
      handleAdd={() => {}}
      handleDelete={() => {}}
      handleUpdate={() => {}}
      handleOnChange={() => {}}
    />
  );
  expect(wrapper.find('UpdateButton')).toHaveLength(1);
  // expect(wrapper.find('[name="email"]')).toHaveLength(1);
});
