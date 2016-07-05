'use strict';

import React from 'react';
import ReactDOM from 'react-dom';

let buttonsInstance = (
    <a href="/oauth/google" className="btn btn-primary">Google Sign in</a>
);

const apiToken = getCookie('apiToken');
if (apiToken !== '') {
    buttonsInstance = (
        <a href="/oauth/google" className="btn btn-primary">You're logged in</a>
    );
} else {
    // TODO
}

ReactDOM.render(buttonsInstance, document.getElementById('app'));

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

