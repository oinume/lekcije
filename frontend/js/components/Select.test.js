import React from 'react';
import renderer from 'react-test-renderer';
import Select from "./Select";

test('hogehoge', () => {
  const options = [
    {value1: 'value1', label: 'label1'},
  ];
  const component = renderer.create(
    <Select
      name="hoge"
      value="fuga"
      onChange={() => {}}
      options={options}
    />
  );
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});
