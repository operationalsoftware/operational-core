document.addEventListener("DOMContentLoaded", () => {
  const unassignForms = document.querySelectorAll(
    ".unassign-service-schedule-form"
  );

  unassignForms.forEach((form) => {
    form.addEventListener("submit", (event) => {
      const message =
        form.dataset.confirmMessage ||
        "Are you sure you want to unassign this schedule?";

      const confirmed = window.confirm(message);
      if (!confirmed) {
        event.preventDefault();
      }
    });
  });
});
