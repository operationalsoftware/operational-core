(function () {
  const form = me("form");

  const reviver = (key, value) => {
    // Check if the value is a string and represents an integer
    if (typeof value === "string" && /^\d+$/.test(value)) {
      return parseInt(value, 10);
    }

    // If not, return the original value
    return value;
  };

  form.on("submit", (e) => {
    e.preventDefault();
    // capture the form data
    const formData = new FormData(form);
    // convert the form data into a JSON object
    const data = Object.fromEntries(formData);
    // Get the multi-select values
    const multiSelectValues = JSON.parse(data["multi-select"], reviver);
    // Remove the multi-select from the form data
    delete data["multi-select"];
    // Add the multi-select values to the form data
    data["multi-select"] = multiSelectValues;
    // Log the form data
    console.log(data);
  });
})();
