function generateMe(startEl) {
  return (selector) => {
    return me(selector, startEl);
  };
}

function setAriaAttribute(el) {
  el.setAttribute(
    "aria-expanded",
    el.getAttribute("aria-expanded") === "true" ? "false" : "true"
  );
}

me().on("keydown", (e) => {
  if (e.key === "Escape") {
    const openModal = me("dialog.modal.open");
    if (openModal) {
      halt(e);
      closeModal(openModal);
    }
  }
});

function getUrlParams(specificParams) {
  const url = window.location.href;
  const params = {};
  const queryString = url ? url.split("?")[1] : window.location.search.slice(1);

  if (queryString) {
    const keyValuePairs = queryString.split("&");

    keyValuePairs.forEach((keyValue) => {
      const [key, value] = keyValue.split("=");
      params[key] = decodeURIComponent(value.replace(/\+/g, " "));
    });

    // Filter specific parameters if specified
    if (specificParams && specificParams.length > 0) {
      const filteredParams = {};
      specificParams.forEach((param) => {
        if (params.hasOwnProperty(param)) {
          filteredParams[param] = params[param];
        }
      });
      return filteredParams;
    }
  }

  return params;
}

function addUrlParams(params) {
  let url = new URL(window.location.href);
  let queryParams = new URLSearchParams(window.location.search);

  params.forEach((param) => {
    queryParams.set(param.name, param.value);
  });
  url.search = queryParams.toString();

  // Use pushState to update the browser URL without reloading the page
  window.history.pushState({ path: url.href }, "", url.href);
}

function removeUrlParams(params) {
  let url = new URL(window.location.href);
  let queryParams = new URLSearchParams(window.location.search);

  params.forEach((param) => {
    queryParams.delete(param);
  });
  url.search = queryParams.toString();

  // Use pushState to update the browser URL without reloading the page
  window.history.pushState({ path: url.href }, "", url.href);
}

function openModal(el) {
  me("body").styles({ overflow: "hidden" });
  el.classRemove("hidden");
  el.classAdd("open");
  el.showModal();
}

function closeModal(el) {
  el.classRemove("open");
  el.classAdd("hidden");
  setTimeout(() => {
    me("body").styles({ overflow: "auto" });
    el.close();
  }, 250);
}

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

(function () {
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
