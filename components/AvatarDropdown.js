(function() {
  const LIGHT_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M3.55 19.09L4.96 20.5L6.76 18.71L5.34 17.29M12 6C8.69 6 6 8.69 6 12S8.69 18 12 18 18 15.31 18 12C18 8.68 15.31 6 12 6M20 13H23V11H20M17.24 18.71L19.04 20.5L20.45 19.09L18.66 17.29M20.45 5L19.04 3.6L17.24 5.39L18.66 6.81M13 1H11V4H13M6.76 5.39L4.96 3.6L3.55 5L5.34 6.81L6.76 5.39M1 13H4V11H1M13 20H11V23H13" /></svg>`;
  const DARK_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 2A9.91 9.91 0 0 0 9 2.46A10 10 0 0 1 9 21.54A10 10 0 1 0 12 2Z" /></svg>`;

  const FULLSCREEN_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M5,5H10V7H7V10H5V5M14,5H19V10H17V7H14V5M17,14H19V19H14V17H17V14M10,17V19H5V14H7V17H10Z" /></svg>`;

  const FULLSCREEN_EXIT_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M14,14H19V16H16V19H14V14M5,14H10V19H8V16H5V14M8,5H10V10H5V8H8V5M19,8V10H14V5H16V8H19Z" /></svg>`;

  let initialIsFullScreen = Boolean(
    document.fullscreenElement ||
    document.mozFullScreenElement ||
    document.webkitFullscreenElement ||
    document.msFullscreenElement
  );

  const avatarDropdown = me(".avatar-dropdown");
  const avatarDropdownContent = me(".dropdown", avatarDropdown);
  const themeToggle = me(".theme-switcher", avatarDropdownContent);
  const fullscreenToggle = me(".fullscreen-switcher", avatarDropdownContent);

  function initThemeToggle() {
    const theme = document.documentElement.getAttribute("data-theme");
    const themeIcon = theme === "dark" ? LIGHT_ICON : DARK_ICON;
    const themeSvg = createIcon(themeIcon);
    themeToggle.innerHTML = "";
    themeToggle.appendChild(themeSvg);
  }

  const toggleFullScreen = () => {
    if (initialIsFullScreen) {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen();
      } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen();
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen();
      }
    } else {
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
    initialIsFullScreen = !initialIsFullScreen;
  };

  function initFullscreenToggle() {
    const fullscreenIcon = initialIsFullScreen
      ? FULLSCREEN_EXIT_ICON
      : FULLSCREEN_ICON;
    const fullscreenSvg = createIcon(fullscreenIcon);
    fullscreenToggle.innerHTML = "";
    fullscreenToggle.appendChild(fullscreenSvg);
  }

  // Toggle dark mode
  themeToggle.addEventListener("click", (e) => {
    e.stopPropagation();
    const currentTheme = document.documentElement.getAttribute("data-theme");
    const newTheme = currentTheme === "dark" ? "light" : "dark";
    document.documentElement.setAttribute("data-theme", newTheme);
    document.cookie = `theme=${newTheme};path=/;max-age=31536000`;
    initThemeToggle();
  });

  // Toggle fullscreen
  fullscreenToggle.on("click", (e) => {
    e.stopPropagation();
    toggleFullScreen();
    initFullscreenToggle();
  });

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

  updateAvatarDropdownPosition();
  initThemeToggle();
  initFullscreenToggle();

  window.addEventListener("resize", updateAvatarDropdownPosition);

  document.addEventListener("click", (e) => {
    if (!e.target.closest(".avatar-dropdown")) {
      avatarDropdownContent.classList.remove("show");
    }
  });
})();
