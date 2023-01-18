/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// identity function for calling harmony imports with the correct context
/******/ 	__webpack_require__.i = function(value) { return value; };
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 3);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
exports['default'] = getAutoLoginUrl;
/**
* 自動ログインの制御を行う
* @param {string} callbackurl - ログイン後のコールバックURL
* @param {object} document - ブラウザオブジェクト
* @param {string|null} hasAltHashCookie - Cookie:has_althashの中身
* @param {string|null} hasAltTimeCookie - Cookie:has_alttimeの中身
* @param {'dev'|'stg'|'prd'} accountsEnv - リダイレクトするアカウントサービスの環境
* @return {string|null} リダイレクト先URL, リダイレクトしない場合はnull
*/
function getAutoLoginUrl(callbackUrl, document, hasAltHashCookie, hasAltTimeCookie) {
    var accountsEnv = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : "prd";

    var getPrefix = function getPrefix(env) {
        if (env === 'dev') {
            return 'dev-';
        }
        if (env === 'stg') {
            return 'stg-';
        }
        return '';
    };

    var getRedirectUrl = function getRedirectUrl(url, env) {
        var tld = document.domain.replace(/^.*\.(com|jp)$/, '$1');
        var domain = tld === 'jp' ? 'co.jp' : 'com';
        return 'https://' + String(getPrefix(env)) + 'accounts.dmm.' + domain + '/service/login/token/=/path=' + String(url);
    };

    var validateArgs = function validateArgs() {
        return typeof callbackUrl === 'string';
    };

    var canRedirect = function canRedirect(hasAltHash, hasAltTime) {
        return hasAltHash && !hasAltTime;
    };

    if (validateArgs() === false) {
        return null;
    }

    if (canRedirect(hasAltHashCookie, hasAltTimeCookie)) {
        return getRedirectUrl(callbackUrl, accountsEnv);
    }
    return null;
}

/***/ }),
/* 1 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


var _browserCookies = __webpack_require__(2);

var _browserCookies2 = _interopRequireDefault(_browserCookies);

var _index = __webpack_require__(0);

var _index2 = _interopRequireDefault(_index);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

(function () {
    var autoLoginElement = document.getElementById('auto-login');
    var path = autoLoginElement.getAttribute('data-encoded');
    var hasAltHash = _browserCookies2['default'].get('has_althash');
    var hasAltTime = _browserCookies2['default'].get('has_alttime');
    var href = (0, _index2['default'])(path, document, hasAltHash, hasAltTime);
    if (href != null) {
        location.href = href;
    }
})();

/***/ }),
/* 2 */
/***/ (function(module, exports) {

exports.defaults = {};

exports.set = function(name, value, options) {
  // Retrieve options and defaults
  var opts = options || {};
  var defaults = exports.defaults;

  // Apply default value for unspecified options
  var expires  = opts.expires  || defaults.expires;
  var domain   = opts.domain   || defaults.domain;
  var path     = opts.path     !== undefined ? opts.path     : (defaults.path !== undefined ? defaults.path : '/');
  var secure   = opts.secure   !== undefined ? opts.secure   : defaults.secure;
  var httponly = opts.httponly !== undefined ? opts.httponly : defaults.httponly;
  var samesite = opts.samesite !== undefined ? opts.samesite : defaults.samesite;

  // Determine cookie expiration date
  // If succesful the result will be a valid Date, otherwise it will be an invalid Date or false(ish)
  var expDate = expires ? new Date(
      // in case expires is an integer, it should specify the number of days till the cookie expires
      typeof expires === 'number' ? new Date().getTime() + (expires * 864e5) :
      // else expires should be either a Date object or in a format recognized by Date.parse()
      expires
  ) : 0;

  // Set cookie
  document.cookie = name.replace(/[^+#$&^`|]/g, encodeURIComponent)                // Encode cookie name
  .replace('(', '%28')
  .replace(')', '%29') +
  '=' + value.replace(/[^+#$&/:<-\[\]-}]/g, encodeURIComponent) +                  // Encode cookie value (RFC6265)
  (expDate && expDate.getTime() >= 0 ? ';expires=' + expDate.toUTCString() : '') + // Add expiration date
  (domain   ? ';domain=' + domain     : '') +                                      // Add domain
  (path     ? ';path='   + path       : '') +                                      // Add path
  (secure   ? ';secure'               : '') +                                      // Add secure option
  (httponly ? ';httponly'             : '') +                                      // Add httponly option
  (samesite ? ';samesite=' + samesite : '');                                       // Add samesite option
};

exports.get = function(name) {
  var cookies = document.cookie.split(';');
  
  // Iterate all cookies
  while(cookies.length) {
    var cookie = cookies.pop();

    // Determine separator index ("name=value")
    var separatorIndex = cookie.indexOf('=');

    // IE<11 emits the equal sign when the cookie value is empty
    separatorIndex = separatorIndex < 0 ? cookie.length : separatorIndex;

    var cookie_name = decodeURIComponent(cookie.slice(0, separatorIndex).replace(/^\s+/, ''));

    // Return cookie value if the name matches
    if (cookie_name === name) {
      return decodeURIComponent(cookie.slice(separatorIndex + 1));
    }
  }

  // Return `null` as the cookie was not found
  return null;
};

exports.erase = function(name, options) {
  exports.set(name, '', {
    expires:  -1,
    domain:   options && options.domain,
    path:     options && options.path,
    secure:   0,
    httponly: 0}
  );
};

exports.all = function() {
  var all = {};
  var cookies = document.cookie.split(';');

  // Iterate all cookies
  while(cookies.length) {
    var cookie = cookies.pop();

    // Determine separator index ("name=value")
    var separatorIndex = cookie.indexOf('=');

    // IE<11 emits the equal sign when the cookie value is empty
    separatorIndex = separatorIndex < 0 ? cookie.length : separatorIndex;

    // add the cookie name and value to the `all` object
    var cookie_name = decodeURIComponent(cookie.slice(0, separatorIndex).replace(/^\s+/, ''));
    all[cookie_name] = decodeURIComponent(cookie.slice(separatorIndex + 1));
  }

  return all;
};


/***/ }),
/* 3 */
/***/ (function(module, exports, __webpack_require__) {

__webpack_require__(0);
module.exports = __webpack_require__(1);


/***/ })
/******/ ]);