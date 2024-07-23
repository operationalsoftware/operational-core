/*
 * getTheme checks if a theme is stored in local storage and returns it.
 * If no theme is stored, it returns "system-default".
 *
 * */
function getTheme() {
  const theme = localStorage.getItem("theme");
  if (theme === "dark") {
    return "dark";
  } else if (theme === "light") {
    return "light";
  } else {
    return "system-default";
  }
}

/*
 * setTheme sets the theme in local storage and updates the data-theme
 * attribute on the html element if the theme is "dark". If the theme is
 * "system-default", we will check if the user prefers dark mode and set
 * the theme accordingly.
 */
function updateThemeOnHtml() {
  const theme = localStorage.getItem("theme");
  const prefersDarkMode = window.matchMedia(
    "(prefers-color-scheme: dark)"
  ).matches;
  if (theme === "dark" || (!theme && prefersDarkMode)) {
    document.documentElement.setAttribute("data-theme", "dark");
  } else {
    document.documentElement.removeAttribute("data-theme");
  }
}

function setTheme(theme) {
  function saveThemeToLocalStorage(theme) {
    localStorage.setItem("theme", theme);
  }
  function removeThemeFromLocalStorage() {
    localStorage.removeItem("theme");
  }

  switch (theme) {
    case "dark":
      saveThemeToLocalStorage("dark");
      break;
    case "light":
      saveThemeToLocalStorage("light");
      break;
    case "system-default":
      removeThemeFromLocalStorage();
      break;
  }

  updateThemeOnHtml();
}

/*
 * Fullscreen utilities
 * */
function getIsFullScreen() {
  return Boolean(
    document.fullscreenElement ||
    document.mozFullScreenElement ||
    document.webkitFullscreenElement ||
    document.msFullscreenElement
  );
}

function requestFullScreen() {
  const element = document.documentElement;
  if (element.requestFullscreen) {
    element.requestFullscreen();
  } else if (element.mozRequestFullScreen) {
    element.mozRequestFullScreen();
  } else if (element.webkitRequestFullscreen) {
    element.webkitRequestFullscreen();
  } else if (element.msRequestFullscreen) {
    element.msRequestFullscreen();
  }
}

function exitFullScreen() {
  if (document.exitFullscreen) {
    document.exitFullscreen();
  } else if (document.mozCancelFullScreen) {
    document.mozCancelFullScreen();
  } else if (document.webkitExitFullscreen) {
    document.webkitExitFullscreen();
  } else if (document.msExitFullscreen) {
    document.msExitFullscreen();
  }
}

(function() {
  // deal with the theme
  // set on start
  updateThemeOnHtml();
  // listen for system theme changes
  // and update the theme accordingly
  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", () => {
      if (!localStorage.getItem("theme")) {
        updateThemeOnHtml();
      }
    });
})();

// show loading message when navigating away from the page if
// the navigation hasn't occurred after 0.5s
document.addEventListener('DOMContentLoaded', function() {
  let timeout;
  function addClassDelayed() {
    const element = document.getElementById('loading-message');
    if (element) {
      element.classList.add('show');
    }
  }
  function handleNavigation() {
    clearTimeout(timeout);
    timeout = setTimeout(addClassDelayed, 500);
  }
  window.addEventListener('beforeunload', handleNavigation);
});

// utility to preserve scroll height on navigation
function preserveScrollHeight() {
  const scrollHeight = window.scrollY || document.documentElement.scrollTop;
  localStorage.setItem('savedScrollHeight', scrollHeight.toString());
}

// restore scroll height when the DOM content has been parsed
document.addEventListener('DOMContentLoaded', function() {
  const savedScrollHeight = localStorage.getItem('savedScrollHeight');
  if (savedScrollHeight !== null) {
    window.scrollTo(0, parseInt(savedScrollHeight, 10));
    // delete from local storage
    localStorage.removeItem('savedScrollHeight');
  }
});
