(function() {
  const container = document.currentScript.closest(".copy-button");
  const btn = container?.querySelector("button.button");
  const statusSpan = container?.querySelector("span.status");

  const textToCopy = container?.dataset?.text || "";

  if (btn && statusSpan) {
    btn.addEventListener("click", function() {
      navigator.clipboard.writeText(textToCopy)
        .then(() => {
          statusSpan.classList.remove("hidden");
          setTimeout(() => {
            statusSpan.classList.add("hidden");
          }, 2000);
        })
        .catch(err => {
          console.error("Clipboard copy failed", err);
        });
    });
  }
})();

