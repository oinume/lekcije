'use strict';

import React from 'react'

export default class FollowTeacherForm extends React.Component {
  render() {
    return (
      <div className="starter-template">
        <h1>Follow a teacher!</h1>
        <form>
          <input type="text" name="teacherUrl" value="" />
        </form>
      </div>
    );
  }
}
