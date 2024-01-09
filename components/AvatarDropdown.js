(function () {
  const avatarDropdown = me(".avatar-dropdown");
  const avatarDropdownContent = me(".content-container", avatarDropdown);

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

  window.addEventListener("resize", updateAvatarDropdownPosition);

  document.addEventListener("click", (e) => {
    if (!e.target.closest(".avatar-dropdown")) {
      avatarDropdownContent.classList.remove("show");
    }
  });
})();
