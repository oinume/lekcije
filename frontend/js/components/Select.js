import PropTypes from 'prop-types';
import React from "react";

export default class Select extends React.Component {

  constructor(props) {
    super(props);
  }

  render() {
    let options = [];
    for (const o of this.props.options) {
      options.push(<option value={o.value}>{o.label}</option>);
    }

    return (
      <select
        name={this.props.name}
        value={this.props.value}
        className={this.props.className}
        onChange={this.props.onChange}
      >
        {options}
      </select>
    );
  }
}

Select.propTypes = {
  name: PropTypes.string.isRequired,
  value: PropTypes.any.isRequired,
  options: PropTypes.array.isRequired,
  onChange: PropTypes.func.isRequired,
  className: PropTypes.string,
};
