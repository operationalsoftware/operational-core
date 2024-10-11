function submitUsersTableForm() {

  const formData = new FormData(document.getElementById("users-table-form")); // 'this' refers to the form element

  // To log form data
  for (const [key, value] of formData.entries()) {
    console.log(`${key}: ${value}`);
  }

  document.getElementById("users-table-form").submit();
}

