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

function createIcon(iconString) {
  const domParser = new DOMParser();
  const iconDoc = domParser.parseFromString(iconString, "image/svg+xml");
  return iconDoc.documentElement;
}

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

(function () {
  // find if theme cookie exists
  const cookies = document.cookie;
  const themeCookie = cookies.split(";").find((cookie) => {
    return cookie.trim().startsWith("theme=");
  });

  if (themeCookie) {
    const theme = themeCookie.split("=")[1];
    document.documentElement.setAttribute("data-theme", theme);
  } else {
    // if theme cookie doesn't exist, set it to dark
    let theme = window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
    document.documentElement.setAttribute("data-theme", theme);
    document.cookie = `theme=${theme};path=/;max-age=31536000`;
  }
})();
