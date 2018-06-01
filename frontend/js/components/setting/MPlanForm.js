import PropTypes from "prop-types";
import React from "react";

export default class MPlanForm extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <form className="form-horizontal">
        <div className="form-group">
          <div className="col-sm-3">
            <label htmlFor="plan" className="control-label">プラン</label>
          </div>
          <div className="col-sm-7">
            <p>{this.props.name}</p>
          </div>
        </div>
        <div className="form-group">
          <div className="col-sm-offset-2 col-sm-8">
          </div>
        </div>
      </form>
    );
  }
}

MPlanForm.propTypes = {
  id: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  dispatch: PropTypes.func.isRequired,
};
