'use strict';

import React from 'react'

export default class FollowTeacherForm extends React.Component {
  render() {
    return (
      <div>
        <h2>Follow a teacher!</h2>
        <form>
          <div className="form-group">
            <label for="teacherUrlOrId">Teacher's URL or ID</label>
            <input type="text"
                   className="form-control"
                   id="teacherIdOrUrl"
                   name="teacherUrlOrId"
                   placeholder="URL or ID" />
          </div>
        </form>
      </div>
    );
  }
}
