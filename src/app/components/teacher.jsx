'use strict';

import React from 'react'

export default class TeacherMain extends React.Component {
  render() {
    return (
      <div className="container">
        <FollowTeacherForm />
        <FollowingTeacherList />
      </div>
    );
  }
}

export class FollowTeacherForm extends React.Component {
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

export class FollowingTeacherList extends React.Component {

  render() {
    return (
      <div>
        <table className="table table-striped table-hover">
          <thead>
            <tr><th>ID</th><th>Name</th></tr>
          </thead>
          <tbody>
            <tr><td>1</td><td>Name1</td></tr>
            <tr><td>2</td><td>Name2</td></tr>
          </tbody>
        </table>
      </div>
    );
  }
}
