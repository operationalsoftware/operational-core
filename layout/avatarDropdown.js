(function () {
  /*
   * Avatar dropdown positioning
   * */
  const avatarDropdown = document.querySelector(".avatar-dropdown");
  const avatarDropdownContent = avatarDropdown.querySelector(".dropdown");
  function updateAvatarDropdownPosition() {
    const buttonRect = avatarDropdown.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();
    const avatarDropdownWidth = avatarDropdownContent.offsetWidth;

    let left = buttonRect.left - bodyRect.left - avatarDropdownWidth;
    const top = buttonRect.top - bodyRect.top + buttonRect.height;

    avatarDropdownContent.style.left = `${left}px`;
    avatarDropdownContent.style.top = `${top}px`;
  }
  // initialise the position
  updateAvatarDropdownPosition();
  // listen for window resize and update the position
  window.addEventListener("resize", updateAvatarDropdownPosition);

  /*
   * Theme toggling
   * */
  function setThemeIcon() {
    // Deal with theme
    const theme = getTheme();
    switch (theme) {
      case "dark":
        document.querySelector(".theme-dark").classList.add("show");
        document.querySelector(".theme-light").classList.remove("show");
        document
          .querySelector(".theme-system-default")
          .classList.remove("show");
        break;
      case "light":
        document.querySelector(".theme-light").classList.add("show");
        document.querySelector(".theme-dark").classList.remove("show");
        document
          .querySelector(".theme-system-default")
          .classList.remove("show");
        break;
      default:
        document.querySelector(".theme-system-default").classList.add("show");
        document.querySelector(".theme-dark").classList.remove("show");
        document.querySelector(".theme-light").classList.remove("show");
    }
  }
  // initialise the theme icon
  setThemeIcon();
  // listen for theme changes and update the theme and the icon
  const themeToggle = avatarDropdownContent.querySelector(".theme-toggle");
  themeToggle.addEventListener("click", (e) => {
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
      document.querySelector(".fullscreen-exit").classList.add("show");
      document.querySelector(".fullscreen").classList.remove("show");
    } else {
      document.querySelector(".fullscreen").classList.add("show");
      document.querySelector(".fullscreen-exit").classList.remove("show");
    }
  }
  // initialise the fullscreen icon
  setFullscreenIcon();
  // listen for fullscreen changes and update the icon
  const fullscreenToggle =
    avatarDropdownContent.querySelector(".fullscreen-toggle");
  fullscreenToggle.addEventListener("click", (e) => {
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

  /*
   * Close the dropdown when clicking outside
   * */
  document.addEventListener("click", (e) => {
    if (!e.target.closest(".avatar-dropdown")) {
      avatarDropdownContent.classList.remove("show");
    }
  });
})();
