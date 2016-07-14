'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import Header from './components/header.jsx';
import LoginForm from './components/login.jsx'
import FollowTeacherForm from './components/teacher.jsx'

let html = '';
const apiToken = getCookie('apiToken');

if (apiToken !== '') {
  html = (
    <div>
      <Header />
      <div className="container">
        <FollowTeacherForm />
      </div>
    </div>
  );
} else {
  html = (
    <div>
      <Header />
      <LoginForm />
    </div>
  )
}

ReactDOM.render(html, document.getElementById('app'));

// https://www.npmjs.com/package/cookie
function getCookie(name) {
  const cname = name + "=";
  const cookies = document.cookie.split(';');
  for (let cookie of cookies) {
    while (cookie.charAt(0) == ' ') {
      cookie = cookie.substring(1);
    }
    if (cookie.indexOf(cname) == 0) {
      return cookie.substring(cname.length, cookie.length);
    }
  }
  return "";
}

