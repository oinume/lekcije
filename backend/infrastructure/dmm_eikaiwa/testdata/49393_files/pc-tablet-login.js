var NaviApiTabletLogin = function () {};

// Opened
NaviApiTabletLogin.openLogin = function () {
  NaviApiTabletLogin.main.style.visibility = 'visible';
  NaviApiTabletLogin.main.style.opacity = 1;
  NaviApiTabletLogin.main.style.zIndex = 1100;
  NaviApiTabletLogin.isOpened = true;
};

// Closed
NaviApiTabletLogin.closeLogin = function () {
  NaviApiTabletLogin.main.style.visibility = 'hidden';
  NaviApiTabletLogin.main.style.opacity = 0;
  NaviApiTabletLogin.main.style.zIndex = 1;
  NaviApiTabletLogin.isOpened = false;
};

// Validate
NaviApiTabletLogin.isLogin = function (tapTargetElement, openNavigationElement) {
  if (tapTargetElement === openNavigationElement) {
    return true;
  }
  if (tapTargetElement.parentNode === null) {
    return false;
  }
  return NaviApiTabletLogin.isLogin(tapTargetElement.parentNode, openNavigationElement);
};

// script attribute defer is mandatory
NaviApiTabletLogin.isOpened = false;
// TODO replace `class` to `id`
NaviApiTabletLogin.root = document.getElementsByClassName('_n4v1-login')[0];
NaviApiTabletLogin.icon = document.getElementsByClassName('_n4v1-login-icon')[0];
NaviApiTabletLogin.main = document.getElementsByClassName('_n4v1-login-main')[0];

// Add Listner:Open
NaviApiTabletLogin.icon.addEventListener('touchend', function () {
  if (NaviApiTabletLogin.isOpened) {
    NaviApiTabletLogin.closeLogin();
  } else {
    NaviApiTabletLogin.openLogin();
  }
});

// Add Listner:Close
document.body.addEventListener('touchend', function (e) {
  if (!NaviApiTabletLogin.isLogin(e.target, NaviApiTabletLogin.root)) {
    NaviApiTabletLogin.closeLogin();
  }
});
