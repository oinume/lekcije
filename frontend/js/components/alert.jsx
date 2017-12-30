import React from 'react';

export default class Alert extends React.Component {

  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.visible) {
      return (
        <div className={"alert alert-" + this.props.kind} role="alert">
          <button
            type="button" className="close"
            onClick={() => this.props.dispatch('hideAlert')}
          >
            &times;
          </button>
          {this.props.message}
        </div>
      );
    } else {
      return (
        <div/>
      );
    }
  }
}
