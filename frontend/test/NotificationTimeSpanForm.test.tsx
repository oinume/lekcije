import React from 'react';
import {configure, shallow} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { NotificationTimeSpanForm } from '../js/components/setting/NotificationTimeSpanForm';

configure({adapter: new Adapter()});

test('<NotificationTimeSpanFormFC>', () => {
  const wrapper = shallow(
    <NotificationTimeSpanForm
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
