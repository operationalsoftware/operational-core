(function() {
  /*
   * Avatar dropdown positioning
   * */
  const avatarDropdown = me(".avatar-dropdown");
  const avatarDropdownContent = me(".dropdown", avatarDropdown);
  function updateAvatarDropdownPosition() {
    const buttonRect = avatarDropdown.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();
    const avatarDropdownWidth = avatarDropdownContent.offsetWidth;

    let left = buttonRect.left - bodyRect.left - avatarDropdownWidth;
    const top = buttonRect.top - bodyRect.top + buttonRect.height;

    avatarDropdownContent.styles({
      left: `${left}px`,
      top: `${top}px`,
    });
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
        me(".theme-dark").classAdd("show");
        me(".theme-light").classRemove("show");
        me(".theme-system-default").classRemove("show");
        break;
      case "light":
        me(".theme-light").classAdd("show");
        me(".theme-dark").classRemove("show");
        me(".theme-system-default").classRemove("show");
        break;
      default:
        me(".theme-system-default").classAdd("show");
        me(".theme-dark").classRemove("show");
        me(".theme-light").classRemove("show");
    }
  }
  // initialise the theme icon
  setThemeIcon();
  // listen for theme changes and update the theme and the icon
  const themeToggle = me(".theme-toggle", avatarDropdownContent);
  themeToggle.on("click", (e) => {
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
      me(".fullscreen-exit").classAdd("show");
      me(".fullscreen").classRemove("show");
    } else {
      me(".fullscreen").classAdd("show");
      me(".fullscreen-exit").classRemove("show");
    }
  }
  // initialise the fullscreen icon
  setFullscreenIcon();
  // listen for fullscreen changes and update the icon
  const fullscreenToggle = me(".fullscreen-toggle", avatarDropdownContent);
  fullscreenToggle.on("click", (e) => {
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
