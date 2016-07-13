'use strict';

import React from 'react'

export default class LoginForm extends React.Component {

  render() {
    return (
      <div className="container">
        <div className="starter-template">
          <a href="/oauth/google" className="btn btn-primary">Google Sign in</a>
        </div>
      </div>
    );
  }
}
