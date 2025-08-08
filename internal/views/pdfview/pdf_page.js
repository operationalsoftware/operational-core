function handleTemplateNameChange(e) {
  const templateNameSelect = e.target;

  const form = templateNameSelect.closest("form");
  if (!form) return;

  e.preventDefault();

  // Build query string from form data
  const params = new URLSearchParams(new FormData(form));

  // Get current URL without query string
  const baseUrl = window.location.pathname;

  // Redirect to current URL with new query parameters
  window.location.href = `${baseUrl}?${params.toString()}`;
}
