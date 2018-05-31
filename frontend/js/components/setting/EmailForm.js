import React from "react";
import PropTypes from "prop-types";

export default class EmailForm extends React.Component {

  constructor(props) {
    super(props);
    this.onChange = this.onChange.bind(this);
  }

  onChange(e) {
    this.props.dispatch('onChangeEmail', e.target.value);
  }

  render() {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="email" className="control-label">Email</label>
          </div>
          <div className="col-sm-7">
            <input
              type="email" className="form-control" name="email" id="email"
              placeholder="Email" required autoFocus autoComplete="on"
              value={this.props.value} onChange={this.onChange}/>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
            <button
              type="button"
              disabled={!this.props.value}
              className="btn btn-primary"
              onClick={() => this.props.dispatch('updateEmail', this.props.value)}
            >
              変更
            </button>
          </div>
        </div>
      </form>
    );
  }
}

EmailForm.propTypes = {
  value: PropTypes.string.isRequired,
  dispatch: PropTypes.func.isRequired,
};
