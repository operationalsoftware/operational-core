(function() {
  const localMe = generateMe(me());
  const modal = localMe("dialog.modal");

  // Modal Logic
  localMe("#open-modal").on("click", () => {
    openModal(modal);
  });

  localMe("#close-btn").on("click", () => {
    closeModal(modal);
  });
})();
