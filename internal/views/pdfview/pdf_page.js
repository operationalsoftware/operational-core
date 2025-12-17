function handleTemplateNameChange(e) {
  const templateNameSelect = e.target;

  const form = templateNameSelect.closest("form");
  if (!form) return;

  e.preventDefault();

  const params = new URLSearchParams(new FormData(form));
  const baseUrl = window.location.pathname;
  window.location.href = `${baseUrl}?${params.toString()}`;
}
