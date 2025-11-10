document.addEventListener("DOMContentLoaded", () => {
  const archiveForms = document.querySelectorAll(
    ".archive-service-schedule-form"
  );

  archiveForms.forEach((form) => {
    form.addEventListener("submit", (event) => {
      const message =
        form.dataset.confirmMessage ||
        "Are you sure you want to archive this schedule?";

      const confirmed = window.confirm(message);
      if (!confirmed) {
        event.preventDefault();
      }
    });
  });
});
