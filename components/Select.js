(function () {
  const closeIcon = "\u2715";

  const selectContainer = document.querySelector(".custom-select");
  const isMultiple = selectContainer.getAttribute("data-multiple") === "true";
  const arrow = selectContainer.querySelector(".arrow");

  // Conditionally set variables
  let multiSelectButton;
  let selectedValuesContainer;
  let defaultTextContainer;
  let defaultValue;

  if (isMultiple) {
    multiSelectButton = selectContainer.querySelector(".multi-select-button");
    selectedValuesContainer = selectContainer.querySelector(".selected-values");
  } else {
    defaultTextContainer = selectContainer.querySelector(".default-value");
    defaultValue = selectContainer.getAttribute("data-default-value");
  }

  let selectedValues = [];

  const defaultText = selectContainer.getAttribute("data-default-text");
  function updateDefaultText() {
    if (!isMultiple && defaultText && defaultTextContainer) {
      defaultTextContainer.textContent = defaultText;
    } else if (isMultiple && selectedValues.length === 0) {
      selectedValuesContainer.textContent = defaultText;
    }
  }

  function replaceIcon() {
    if (defaultTextContainer.textContent !== defaultText) {
      arrow.classList.remove("arrow");
      arrow.classList.add("remove-icon");
      arrow.innerHTML = closeIcon;
    } else {
      arrow.classList.remove("remove-icon");
      arrow.classList.add("arrow");
      arrow.innerHTML = "";
    }
  }

  // Set default option
  function selectDefaultOption() {
    const options = selectContainer.querySelectorAll("input[type='radio']");
    options.forEach((option) => {
      if (option.value === defaultValue) {
        option.checked = true;
        defaultTextContainer.textContent =
          option.parentElement.querySelector("label").textContent;
        replaceIcon();
      }
    });
  }
  document.addEventListener("DOMContentLoaded", selectDefaultOption);

  selectContainer.addEventListener("click", (e) => {
    e.stopPropagation();
    // Simple Select
    if (e.target.matches(".select-button")) {
      selectContainer.classList.toggle("active");
    }
    if (e.target.matches(".select-dropdown li")) {
      const parent = e.target.children[0];
      if (parent) {
        const label = parent.children[0].textContent;
        defaultTextContainer.textContent = label;
        if (selectContainer) {
          selectContainer.classList.toggle("active");
        }
        replaceIcon();
      }
    }

    if (!isMultiple) {
      if (e.target.matches(".remove-icon")) {
        defaultTextContainer.textContent = defaultText;
        replaceIcon();
      }
    }

    // Multi Select
    if (e.target.matches(".multi-select-button")) {
      selectContainer.classList.toggle("active");
    }

    if (e.target.matches(".option")) {
      const checkbox = e.target.querySelector("input[type='checkbox']");
      const label = e.target.querySelector("label").textContent;
      const value = checkbox.value;
      checkbox.checked = String(!checkbox.checked);
      const index = selectedValues.indexOf(value);
      // Create Badges
      const createSelectedValue = document.createElement("span");
      const closeIconElement = document.createElement("span");
      closeIconElement.classList.add("remove-icon");
      createSelectedValue.classList.add("selected-value");
      createSelectedValue.setAttribute("data-value", value);
      createSelectedValue.textContent = label;
      closeIconElement.textContent = closeIcon;
      createSelectedValue.appendChild(closeIconElement);

      if (index === -1) {
        // Remove default text
        if (selectedValues.length === 0) {
          selectedValuesContainer.textContent = "";
        }
        selectedValues.push(value);
        e.target.classList.add("checked");
        e.target.setAttribute("data-value", value);
        selectedValuesContainer.appendChild(createSelectedValue);
      } else {
        e.target.classList.remove("checked");
        e.target.removeAttribute("data-value");
        selectedValues.splice(index, 1);
        const selectedValue = selectedValuesContainer.querySelector(
          `.selected-value[data-value='${value}']`
        );
        selectedValuesContainer.removeChild(selectedValue);
      }
      updateDefaultText();
    }
  });

  // if (selectedValuesContainer) {
  //   selectedValuesContainer.addEventListener("click", (e) => {
  //     const removeIcon = e.target.closest(".remove-icon");
  //     if (removeIcon) {
  //       const parent = removeIcon.parentNode;
  //       const value = parent.getAttribute("data-value");
  //       const index = selectedValues.indexOf(value);
  //       const option = document.querySelector(`.option[data-value='${value}']`);
  //       option.classList.remove("checked");
  //       option.removeAttribute("data-value");
  //       selectedValues.splice(index, 1);
  //       selectedValuesContainer.removeChild(parent);
  //       multiSelectHiddenInput.value = JSON.stringify(selectedValues);
  //       updateDefaultText();
  //     }
  //   });
  // }

  // Close dropdown when clicked outside
  document.addEventListener("click", (e) => {
    if (!selectContainer.contains(e.target)) {
      selectContainer.classList.remove("active");
    }
  });

  updateDefaultText();
})();
