'use strict';

import React from 'react';
import {EmailForm} from "../js/components/setting/EmailForm.tsx";
import {configure, shallow} from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import MicroContainer from 'react-micro-container';

configure({adapter: new Adapter()});

test('<EmailForm>', () => {
  //class Container extends MicroContainer {}

  // let container = new Container();
  // container.subscribe({
  //   onChangeEmail: () => {},
  //   updateEmail: () => {},
  // });
 // container.dispatch('onChangeEmail', 'updateEmail');

  const handleOnChange = (e) => {};
  const handleUpdateEmail = (email) => {};
  const wrapper = shallow(
    <EmailForm
      email="oinume@gmail.com"
      handleOnChange={handleOnChange}
      handleUpdateEmail={handleUpdateEmail}
    />
  );
  expect(wrapper.find('button')).toHaveLength(1);
  expect(wrapper.find('[name="email"]')).toHaveLength(1);
});
