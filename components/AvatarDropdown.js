(function () {
  const LIGHT_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M3.55 19.09L4.96 20.5L6.76 18.71L5.34 17.29M12 6C8.69 6 6 8.69 6 12S8.69 18 12 18 18 15.31 18 12C18 8.68 15.31 6 12 6M20 13H23V11H20M17.24 18.71L19.04 20.5L20.45 19.09L18.66 17.29M20.45 5L19.04 3.6L17.24 5.39L18.66 6.81M13 1H11V4H13M6.76 5.39L4.96 3.6L3.55 5L5.34 6.81L6.76 5.39M1 13H4V11H1M13 20H11V23H13" /></svg>`;
  const DARK_ICON = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 2A9.91 9.91 0 0 0 9 2.46A10 10 0 0 1 9 21.54A10 10 0 1 0 12 2Z" /></svg>`;

  const avatarDropdown = me(".avatar-dropdown");
  const avatarDropdownContent = me(".content-container", avatarDropdown);
  const themeToggle = me(".theme-switcher", avatarDropdownContent);

  function initThemeToggle() {
    const theme = document.documentElement.getAttribute("data-theme");
    const themeIcon = theme === "dark" ? LIGHT_ICON : DARK_ICON;
    const themeSvg = createIcon(themeIcon);
    themeToggle.innerHTML = "";
    themeToggle.appendChild(themeSvg);
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

  window.addEventListener("resize", updateAvatarDropdownPosition);

  document.addEventListener("click", (e) => {
    if (!e.target.closest(".avatar-dropdown")) {
      avatarDropdownContent.classList.remove("show");
    }
  });
})();
