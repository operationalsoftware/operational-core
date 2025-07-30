document.addEventListener("DOMContentLoaded", function() {
  const container = document.currentScript.closest(".clipboard-container");
  const btn = container?.querySelector(".clipboard-btn");
  const status = container?.querySelector(".clipboard-status");

  const textToCopy = btn?.dataset?.text || "";

  if (btn && status) {
    btn.addEventListener("click", function() {
      navigator.clipboard.writeText(textToCopy)
        .then(() => {
          status.classList.remove("hidden");
          setTimeout(() => {
            status.classList.add("hidden");
          }, 2000);
        })
        .catch(err => {
          console.error("Clipboard copy failed", err);
        });
    });
  }
});
