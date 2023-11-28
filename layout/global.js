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
