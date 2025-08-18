// Function to remove duplicate page size fields and submit the form
function submitTableForm(form) {
  // Find all select elements with name PageSize
  const pageSizeSelects = form.querySelectorAll("select.page-size-select");

  // Remove `name` attribute from all but one to avoid duplicates in form submission
  pageSizeSelects.forEach((select, index) => {
    if (index !== 0) {
      select.removeAttribute("name");
    }
  });

  // Submit the form
  form.submit();
}

// Function to update PageSize, sync selects, and then submit the form
function updatePageSizeAndSubmit(selectElement) {
  // Find the closest form that contains the changed select element
  const form = selectElement.closest("form");

  // Find all select elements within the same form (pagination controls)
  const pageSizeSelects = form.querySelectorAll("select.page-size-select");

  // Sync all select elements' values within the same form
  pageSizeSelects.forEach((select) => (select.value = selectElement.value));

  // Submit the form using the utility function
  submitTableForm(form);
}

function rowClickNavigate(event) {
  const interactiveTags = ["BUTTON", "A", "INPUT", "SELECT", "TEXTAREA"];

  if (interactiveTags.includes(event.target.tagName)) return;

  const selection = window.getSelection();
  if (selection && selection.toString().length > 0) return;

  const row = event.target.closest("tr");

  const url = row.getAttribute("data-href");

  if (!url) return;

  window.location.href = url;
}
