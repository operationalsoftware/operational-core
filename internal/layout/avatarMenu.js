"use strict";
// block scoping
{
  const navbarEl = document.getElementById("navbar");
  const buttonEl = document.getElementById("navbar-avatar-menu-button");
  const panelEl = document.getElementById("navbar-avatar-menu");
  const navbarCollapseEl = document.getElementById(
    "navbar-collapse-menu-button"
  );
  const navbarExpandEl = document.getElementById("navbar-expand-menu-button");

  document.addEventListener("DOMContentLoaded", () => {
    const makeNavbarCollapse = JSON.parse(
      localStorage.getItem("navbar-collapse")
    );

    if (makeNavbarCollapse) {
      navbarEl.classList.add("hidden");
      navbarExpandEl.classList.add("hidden");
    }
  });

  buttonEl.addEventListener("click", showPanel);

  function showPanel() {
    panelEl.classList.toggle("show");
  }

  // Add click event listener to the document to close the panel on click outside
  document.addEventListener("click", closePanel);

  function closePanel(event) {
    // Check if the click is outside both the button and the panel
    if (!buttonEl.contains(event.target) && !panelEl.contains(event.target)) {
      panelEl.classList.remove("show");
    }
  }

  /*
   * Theme toggling
   * */
  function setThemeIcon() {
    // Deal with theme
    const theme = getTheme();
    switch (theme) {
      case "dark":
        document.querySelector(".theme-dark-icon").classList.add("show");
        document.querySelector(".theme-light-icon").classList.remove("show");
        document
          .querySelector(".theme-system-default-icon")
          .classList.remove("show");
        break;
      case "light":
        document.querySelector(".theme-light-icon").classList.add("show");
        document.querySelector(".theme-dark-icon").classList.remove("show");
        document
          .querySelector(".theme-system-default-icon")
          .classList.remove("show");
        break;
      default:
        document
          .querySelector(".theme-system-default-icon")
          .classList.add("show");
        document.querySelector(".theme-dark-icon").classList.remove("show");
        document.querySelector(".theme-light-icon").classList.remove("show");
    }
  }
  // initialise the theme icon
  setThemeIcon();
  // listen for theme changes and update the theme and the icon
  const themeToggleButtonEl = document.getElementById("theme-toggle-button");
  themeToggleButtonEl.addEventListener("click", (e) => {
    e.stopPropagation();
    const currentTheme = getTheme();

    // system goes to dark, dark goes to light, light goes to system
    switch (currentTheme) {
      case "dark":
        setTheme("light");
        break;
      case "light":
        setTheme("system-default");
        break;
      default:
        setTheme("dark");
    }

    setThemeIcon();
  });

  /*
   * Fullscreen toggling
   * */
  function setFullscreenIcon() {
    const isFullScreen = getIsFullScreen();
    if (isFullScreen) {
      document.querySelector(".fullscreen-exit-icon").classList.add("show");
      document.querySelector(".fullscreen-icon").classList.remove("show");
    } else {
      document.querySelector(".fullscreen-icon").classList.add("show");
      document.querySelector(".fullscreen-exit-icon").classList.remove("show");
    }
  }
  // initialise the fullscreen icon
  setFullscreenIcon();
  // listen for fullscreen changes and update the icon
  const fullscreenToggleButtonEl = document.getElementById(
    "fullscreen-toggle-button"
  );
  fullscreenToggleButtonEl.addEventListener("click", (e) => {
    e.stopPropagation();
    const isFullScreen = getIsFullScreen();
    if (isFullScreen) {
      exitFullScreen();
    } else {
      requestFullScreen();
    }
    setFullscreenIcon();
  });
  // listen for changes to fullscreen and update the icon
  document.addEventListener("fullscreenchange", setFullscreenIcon);

  document.addEventListener("keydown", (event) => {
    if (event.ctrlKey && event.key === "/") {
      event.preventDefault();
      window.location.href = "/search";
    }
  });

  navbarCollapseEl.addEventListener("click", (e) => {
    if (navbarEl) {
      navbarEl.classList.add("hidden");
      localStorage.setItem("navbar-collapse", true);
    }
  });

  navbarExpandEl.addEventListener("click", (e) => {
    if (navbarEl) {
      navbarEl.classList.remove("hidden");
      localStorage.setItem("navbar-collapse", false);
      navbarExpandEl.classList.add("hidden");
    }
  });

  document.addEventListener("mousemove", (e) => {
    const widthThreshold = window.outerWidth - e.clientX;
    if (
      e.clientY <= 150 &&
      widthThreshold < 150 &&
      navbarEl.classList.contains("hidden")
    ) {
      navbarExpandEl.classList.remove("hidden");
    } else {
      navbarExpandEl.classList.add("hidden");
    }
  });
}
