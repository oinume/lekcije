import React from 'react';
import RaisedButton from 'material-ui/lib/raised-button';

class FollowingTeacherForm extends React.Component {

  constructor(props, context) {
    super(props, context);
  }

  handleTouchTap() {
    alert('handleTouchTap()');
  }

  render() {
    return (
      <RaisedButton
        label="Press this"
        primary={true}
        onTouchTap={this.handleTouchTap}
      />
    );
  }

}

export default FollowingTeacherForm;
