(function () {
  const appGalleryButton = document.querySelector(".app-gallery__button");
  const appGalleryContent = document.querySelector(
    ".app-gallery-content__container"
  );
  const navbar = document.getElementById("navbar");

  function updateAppGalleryPosition() {
    const buttonRect = appGalleryButton.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();
    const navbarRect = navbar.getBoundingClientRect();
    const appGalleryWidth = appGalleryContent.offsetWidth;

    let left = buttonRect.left - bodyRect.left - appGalleryWidth;
    const top = navbarRect.bottom + 10;

    appGalleryContent.style.left = `${left}px`;
    appGalleryContent.style.top = `${top}px`;
  }

  updateAppGalleryPosition();

  window.addEventListener("resize", updateAppGalleryPosition);

  document.addEventListener("click", (e) => {
    if (!e.target.closest(".app-gallery")) {
      appGalleryContent.classList.remove("show");
    }
  });
})();
