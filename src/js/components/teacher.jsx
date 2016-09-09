'use strict';

import React from 'react'

export default class TeacherMain extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      teachers: [
        { id: 1, name: "Xai" },
        { id: 2, name: "Emina" },
        { id: 3, name: "c" },
      ]
    };
  }

  render() {
    return (
      <div className="container">
        <FollowTeacherForm />
        <FollowingTeacherList teachers={this.state.teachers} />
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
            <label htmlFor="teacherIdsOrUrl">Teacher's URL or ID</label>
            <input type="text"
                   className="form-control"
                   id="teacherIdsOrUrl"
                   name="teacherIdsOrUrl"
                   placeholder="URL or ID" />
          </div>
          <div className="form-group">
            <button type="submit" className="btn btn-default">Submit</button>
          </div>
        </form>
      </div>
    );
  }
}

export class FollowingTeacherList extends React.Component {

  render() {
    const teachers = [];
    for (let teacher of this.props.teachers) {
      teachers.push(<Teacher id={teacher.id} name={teacher.name} />);
    }
    return (
      <div>
        <table className="table table-striped table-hover">
          <thead>
          <tr><th>ID</th><th>Name</th></tr>
          </thead>
          <tbody>
            {teachers}
          </tbody>
        </table>
      </div>
    );
  }
}

export class Teacher extends React.Component {
  render() {
    return (
      <tr><td>{this.props.id}</td><td>{this.props.name}</td></tr>
    )
  }
}
