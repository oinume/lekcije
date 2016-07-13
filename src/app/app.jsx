'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import Header from './components/header.jsx';
import LoginForm from './components/login.jsx'

let html = '';
const apiToken = getCookie('apiToken');

if (apiToken !== '') {
  html = (
    <div>
      <Header />
      <div className="container">
        <div className="starter-template">
        <h1>Bootstrap starter template</h1>
        <p className="lead">Use this document as a way to quickly start any new project.<br /> All you get is this text and a mostly barebones HTML document.</p>
        </div>
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

