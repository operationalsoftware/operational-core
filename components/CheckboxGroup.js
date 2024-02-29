(function () {
  const checkboxGroup = document.querySelector(".checkbox-group");
  const checkboxOptions = checkboxGroup.querySelector(".checkbox-options");
  const hiddenInput = checkboxGroup.querySelector(".hidden-input");
  const options = [];

  checkboxOptions.addEventListener("change", (e) => {
    if (e.target.tagName === "INPUT") {
      const option = e.target.value;
      if (e.target.checked) {
        options.push(option);
        hiddenInput.value = JSON.stringify(options);
        if (e.target.hasAttribute("checked")) {
          e.target.setAttribute("checked", "true");
        } else {
          e.target.setAttribute("checked", "true");
        }
      } else {
        options.splice(options.indexOf(option), 1);
        hiddenInput.value = JSON.stringify(options);
        if (e.target.hasAttribute("checked")) {
          e.target.setAttribute("checked", "false");
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
