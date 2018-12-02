$(function() {
  var $menuHide = $("#menu-hide");
  var $menuShow = $("#menu-show");
  var $menuList = $(".pure-menu-list");

  $menuHide.click(function() {
    document.cookie = "hide=1; path=/";
    $menuHide.hide();
    $menuShow.show();
    $menuList.hide();
  });

  $menuShow.click(function() {
    document.cookie = "hide=0; path=/";
    $menuHide.show();
    $menuShow.hide();
    $menuList.show();
  });

  if (document.cookie == "hide=1") {
    $menuHide.hide();
    $menuShow.show();
    $menuList.hide();
  } else {
    $menuHide.show();
    $menuShow.hide();
    $menuList.show();
  }
});