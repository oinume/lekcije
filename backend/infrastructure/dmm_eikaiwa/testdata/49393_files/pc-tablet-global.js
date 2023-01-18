// open/close control for iPad
var NaviApiPcTabletGlobal = function () {};

var otherService = document.getElementsByClassName('_n4v1-global-other-service')[0];
if (typeof otherService !== 'undefined') {
  otherService.addEventListener('click', function () {
    document.getElementsByClassName('_n4v1-global-navi')[0].style.visibility = 'hidden';
    document.getElementsByClassName('_n4v1-global-navi')[0].classList.remove('_n4v1-global-navi');
    document.getElementsByClassName('_n4v1-global-navi-extention')[0].classList.add('_n4v1-global-navi');
    document.getElementsByClassName('_n4v1-global-navi')[0].classList.remove('_n4v1-global-navi-extention');
  });
}

// open global navi
NaviApiPcTabletGlobal.openGlobal = function () {
  NaviApiPcTabletGlobal.main.style.visibility = 'visible';
  NaviApiPcTabletGlobal.main.style.opacity = 1;
  NaviApiPcTabletGlobal.main.style.zIndex = 1100;
  NaviApiPcTabletGlobal.isOpened = true;
};

// close global navi
NaviApiPcTabletGlobal.closeGlobal = function () {
  NaviApiPcTabletGlobal.main.style.visibility = 'hidden';
  NaviApiPcTabletGlobal.main.style.opacity = 0;
  NaviApiPcTabletGlobal.main.style.zIndex = 1;
  NaviApiPcTabletGlobal.isOpened = false;
};

// judge if `openNavigationElement` contains `tapTargetElement`
NaviApiPcTabletGlobal.checkTarget = function (tapTargetElement, openNavigationElement) {
  if (tapTargetElement === openNavigationElement) {
    return true;
  }
  if (tapTargetElement.parentNode === null) {
    return false;
  }
  return NaviApiPcTabletGlobal.checkTarget(tapTargetElement.parentNode, openNavigationElement);
};

// script attribute defer is mandatory
NaviApiPcTabletGlobal.isOpened = false;
// TODO replce `class` to `id`
NaviApiPcTabletGlobal.root = document.getElementsByClassName('_n4v1-global')[0];
NaviApiPcTabletGlobal.icon = document.getElementsByClassName('_n4v1-global-icon')[0];
NaviApiPcTabletGlobal.main = document.getElementsByClassName('_n4v1-global-navi')[0];

// add event listener for tapping global navi
NaviApiPcTabletGlobal.icon.addEventListener('touchend', function () {
  if (NaviApiPcTabletGlobal.isOpened) {
    NaviApiPcTabletGlobal.closeGlobal();
  } else {
    NaviApiPcTabletGlobal.openGlobal();
  }
});

// add event listener for outside of global navi
document.body.addEventListener('touchend', function (e) {
  if (!NaviApiPcTabletGlobal.checkTarget(e.target, NaviApiPcTabletGlobal.root)) {
    NaviApiPcTabletGlobal.closeGlobal();
  }
});
