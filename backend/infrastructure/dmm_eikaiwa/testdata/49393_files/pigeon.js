var NaviApiPigeon = {};

/**
 * send data to pigeon system
 *
 * @param {string} endpoint endpoint for pigeon
 * @param {object} params POST payload
 * @param {object|null} options option for XHR (pass `null` if unnecessary)
 *
 * contents of `options`:
 * {number} timeout: time to timeout[ms]
 * {function} callbackOnTimeout: callback for timeout
 * {function} callbackOnUpload: callback for completion of uploading data
 */
NaviApiPigeon.request = function (endpoint, params, options) {
  if (navigator.sendBeacon) {
    navigator.sendBeacon(
      endpoint,
      new Blob([JSON.stringify(params)], { type: 'application/json' })
    );
    return;
  }
  // use XHR if sendBeacon is not supported
  var request = new XMLHttpRequest(); // eslint-disable-line vars-on-top
  request.open('POST', endpoint, true);
  request.setRequestHeader('Content-Type', 'application/json;charset=UTF-8');
  if (options) {
    if (options.timeout) {
      request.timeout = options.timeout;
    }
    if (options.callbackOnTimeout) {
      request.ontimeout = options.callbackOnTimeout;
    }
    if (options.callbackOnUpload) {
      request.upload.onloadend = options.callbackOnUpload;
    }
  }
  request.send(JSON.stringify(params));
};

/**
 * judge if the browser is iOS Safari (from BD department code)
 *
 * @method checkIOS
 * @private
 * @returns {boolean}
 */
NaviApiPigeon.checkIOS = function () {
  var ua = navigator.userAgent.toLowerCase();
  var isiPhone = ua.indexOf('iphone') > -1;
  var isiPod = ua.indexOf('ipod') > -1;
  var isiPad = ua.indexOf('ipad') > -1;
  var isSafari = ua.indexOf('safari') > -1;
  var isChrome = ua.indexOf('chrome') > -1;
  return (isiPhone || isiPod || isiPad) && isSafari && !isChrome;
};

/**
 * generate domain for Pigeon API endpoint
 *
 * @return {string}
 */
NaviApiPigeon.getDomain = function () {
  var env = document.getElementById('naviapi-pigeon').getAttribute('data-env');

  // ios
  if (/dmm\.co\.jp$/.test(window.location.hostname) // isR18
    && NaviApiPigeon.checkIOS()) {
    if (env === 'stg') {
      return 'stg-pigeon.i3.dmm.co.jp';
    }
    return 'pigeon.i3.dmm.co.jp';
  }

  // others
  if (env === 'stg') {
    return 'stg-pigeon.i3.dmm.com';
  }
  return 'pigeon.i3.dmm.com';
};

/**
 * generate Pigeon API endpoint
 *
 * @param {string} various path parameter
 * @param {string} action path parameter
 * @param {string} option path parameter
 * @return {string}
 */
NaviApiPigeon.getEndPoint = function (various, action, option) {
  var version = document.getElementById('naviapi-pigeon').getAttribute('data-version');
  return 'https://' + NaviApiPigeon.getDomain() + '/' + version + '/' + various + '/' + action + '/' + option;
};

/**
 * get cookie value for given `key`
 *
 * @param {string} key key for Cookie
 * @return {string} value of Cookie (empty string if not exist)
 */
NaviApiPigeon.getCookie = function (key) {
  var name = key + '=';
  var cookies;
  var cookie;
  var index;
  try {
    cookies = document.cookie.split(/;\s*/);
  } catch (e) {
    // sometimes accessing `document.cookie` raises `SecurityError` on Android Chrome
    return '';
  }
  for (index = 0; index < cookies.length; index += 1) {
    cookie = cookies[index];
    if (cookie.indexOf(name) === 0) {
      return cookie.substring(name.length, cookie.length);
    }
  }
  return '';
};

/**
 * generate common parameter for pigeon event
 *
 * @return {Object}
 */
NaviApiPigeon.createCommonParameters = function () {
  return {
    site_type: document.getElementById('naviapi-pigeon').getAttribute('data-site-type'),
    view_type: document.getElementById('naviapi-pigeon').getAttribute('data-view-type'),
    url: document.location.href,
    referer: document.referrer,
    segment: NaviApiPigeon.getCookie('i3_ab'),
    open_id: NaviApiPigeon.getCookie('i3_opnd'),
    cdp_id: NaviApiPigeon.getCookie('cdp_id')
  };
};

/**
 * callback for click event
 */
NaviApiPigeon.click = function (event) {
  var currentTarget = event.currentTarget;
  var options = null;
  var params = NaviApiPigeon.createCommonParameters();

  var dataExtraMap = currentTarget.getAttribute('data-extra-map');
  var map = null;
  if (dataExtraMap) {
    // create parameter from `data-extra-map` attribute (in JSON format)
    map = JSON.parse(dataExtraMap);
    Object.keys(map).forEach(function (key) {
      params[key] = map[key];
    });
  }

  try {
    NaviApiPigeon.request(
      NaviApiPigeon.getEndPoint(
        currentTarget.getAttribute('data-various'),
        'click',
        currentTarget.getAttribute('data-option')
      ),
      params,
      options
    );
  } catch (e) {
    // continue regardless of error
  }
};

/**
 * callback for opening navi
 */
NaviApiPigeon.open = function (event) {
  var currentTarget = event.currentTarget;
  var options = null;
  var params = NaviApiPigeon.createCommonParameters();

  try {
    NaviApiPigeon.request(
      NaviApiPigeon.getEndPoint(
        currentTarget.getAttribute('data-various'),
        'open',
        currentTarget.getAttribute('data-option')
      ),
      params,
      options
    );
  } catch (e) {
    // continue regardless of error
  }
};

/**
 * add event listeners
 */
NaviApiPigeon.clickElements = document.querySelectorAll('[data-navi-pigeon-click]');
Array.prototype.forEach.call(NaviApiPigeon.clickElements, function (element) {
  element.addEventListener('click', NaviApiPigeon.click);
});

NaviApiPigeon.openElements = document.querySelectorAll('[data-navi-pigeon-open]');
Array.prototype.forEach.call(NaviApiPigeon.openElements, function (element) {
  element.addEventListener('click', NaviApiPigeon.open);
});
