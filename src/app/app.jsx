'use strict';

import React from 'react';
import ReactDOM from 'react-dom';

let html = (
  <a href="/oauth/google" className="btn btn-primary">Google Sign in</a>
);
const apiToken = getCookie('apiToken');
if (apiToken !== '') {
  html = (
    <div className="container">
      <nav className="navbar navbar-inverse navbar-fixed-top">
        <div className="container">
          <div className="navbar-header">
            <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
              <span className="sr-only">Toggle navigation</span>
              <span className="icon-bar"></span>
              <span className="icon-bar"></span>
              <span className="icon-bar"></span>
            </button>
            <a className="navbar-brand" href="#">Project name</a>
          </div>
          <div id="navbar" className="collapse navbar-collapse">
            <ul className="nav navbar-nav">
              <li className="active"><a href="#">Home</a></li>
              <li><a href="#about">About</a></li>
              <li><a href="#contact">Contact</a></li>
            </ul>
          </div>
        </div>
      </nav>
      <hr />
      <div className="container">
        <div className="starter-template">
        <h1>Bootstrap starter template</h1>
        <p className="lead">Use this document as a way to quickly start any new project.<br /> All you get is this text and a mostly barebones HTML document.</p>
        </div>
      </div>
    </div>
  );
} else {
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

