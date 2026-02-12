"use strict";
// block scoping
{
  const closeButtons = document.querySelectorAll("[data-alert-close]");
  closeButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const alert = button.closest(".alert");
      if (alert) {
        alert.classList.add("hide");
      }
    });
  });
}
