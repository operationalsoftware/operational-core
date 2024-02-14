(function () {
  const checkboxGroup = me();
  const checkboxOptions = me(".checkbox-options");
  const hiddenInput = me(".hidden-input", checkboxGroup);
  const options = [];

  checkboxOptions.on("change", (e) => {
    if (e.target.tagName === "INPUT") {
      const option = e.target.value;
      if (e.target.checked) {
        options.push(option);
        hiddenInput.value = JSON.stringify(options);
        if (e.target.attributes["checked"]) {
          e.target.attributes["checked"].value = "true";
        } else {
          e.target.setAttribute("checked", "true");
        }
      } else {
        options.splice(options.indexOf(option), 1);
        hiddenInput.value = JSON.stringify(options);
        if (e.target.attributes["checked"]) {
          e.target.attributes["checked"].value = "false";
        } else {
          e.target.removeAttribute("checked");
        }
      }
    }
  });

  document.addEventListener("DOMContentLoaded", () => {
    const checkboxOption = checkboxOptions.querySelectorAll("input");
    checkboxOption.forEach((option) => {
      if (option.checked) {
        options.push(option.value);
        hiddenInput.value = JSON.stringify(options);
      }
    });
  });
})();
