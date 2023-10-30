(function () {
  // Modal Logic
  function closeModal() {
    me(".modal").classRemove("open");
    me(".modal").classAdd("hidden");
    setTimeout(() => {
      me("body").styles({ overflow: "auto" });
      me(".modal").close();
    }, 250);
  }

  me("#open-modal").on("click", () => {
    me("body").styles({ overflow: "hidden" });
    me(".modal").classRemove("hidden");
    me(".modal").classAdd("open");
    me(".modal").showModal();
  });

  me("#close-btn").on("click", () => {
    closeModal();
  });

  me("body").on("keydown", (e) => {
    halt(e);
    if (e.key === "Escape") {
      closeModal();
    }
  });
})();
