(function () {
  const closeIcon = "\u2715";

  const selectContainer = document.querySelector(".custom-select");
  const isMultiple = selectContainer.getAttribute("data-multiple") === "true";
  const arrow = selectContainer.querySelector(".arrow");
  const placeholderTextContainer =
    selectContainer.querySelector(".placeholder");
  const selectedValueContainer = selectContainer.querySelector(
    ".select-value-container"
  );

  const placeholderText = selectContainer.getAttribute("data-placeholder");
  const selectName = selectContainer.getAttribute("data-name");

  let multiSelectValues = [];

  function replaceIcon() {
    if (!selectedValueContainer.classList.contains("hidden")) {
      arrow.classList.remove("arrow");
      arrow.classList.add("remove-icon");
      arrow.innerHTML = closeIcon;
    } else {
      arrow.classList.remove("remove-icon");
      arrow.classList.add("arrow");
      arrow.innerHTML = "";
    }
  }

  selectContainer.addEventListener("click", (e) => {
    e.preventDefault();
    e.stopPropagation();
    // Toggle dropdown
    if (e.target.matches(".select-button")) {
      selectContainer.classList.toggle("active");
    }

    if (e.target.matches(".select-dropdown .option")) {
      if (isMultiple) {
        const label = e.target.childNodes[0].textContent;
        const value = e.target.getAttribute("data-value");
        const index = multiSelectValues.indexOf(value);

        // Create Badges
        const createSelectedValue = document.createElement("span");
        const closeIconElement = document.createElement("span");

        closeIconElement.classList.add("remove-icon");
        createSelectedValue.classList.add("selected-value");
        createSelectedValue.setAttribute("data-value", value);
        createSelectedValue.textContent = label;
        closeIconElement.textContent = closeIcon;
        createSelectedValue.appendChild(closeIconElement);

        // Create Checkbox
        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.name = selectName;
        checkbox.value = value;
        checkbox.checked = true;

        if (index === -1) {
          multiSelectValues.push(value);
          e.target.classList.add("checked");
          createSelectedValue.appendChild(checkbox);
          selectedValueContainer.appendChild(createSelectedValue);
        } else {
          multiSelectValues.splice(index, 1);
          e.target.classList.remove("checked");
          const selectedValue = selectedValueContainer.querySelector(
            `.selected-value[data-value='${value}']`
          );
          selectedValueContainer.removeChild(selectedValue);
        }
        selectedValueContainer.classList.remove("hidden");
        placeholderTextContainer.classList.add("hidden");
      } else {
        if (!selectedValueContainer.classList.contains("hidden")) {
          return;
        }
        const label = e.target.childNodes[0].textContent;
        const value = e.target.getAttribute("data-value");
        // create a radio button
        const radio = document.createElement("input");
        radio.type = "radio";
        radio.name = "select";
        radio.value = value;
        radio.checked = true;

        // create a label
        const labelElement = document.createElement("label");
        labelElement.textContent = label;

        labelElement.appendChild(radio);
        selectedValueContainer.textContent = "";
        selectedValueContainer.appendChild(labelElement);
        e.target.classList.add("checked");
        placeholderTextContainer.classList.add("hidden");
        selectedValueContainer.classList.remove("hidden");
        replaceIcon();
      }
    }

    // Removing selected value
    if (e.target.matches(".remove-icon")) {
      if (isMultiple) {
        const selectedValue = e.target.parentNode;
        const value = selectedValue.getAttribute("data-value");
        const index = multiSelectValues.indexOf(value);
        const option = selectContainer.querySelector(
          `.option[data-value='${value}']`
        );
        option.classList.remove("checked");
        multiSelectValues.splice(index, 1);
        selectedValueContainer.removeChild(selectedValue);
        if (multiSelectValues.length < 1) {
          selectedValueContainer.classList.add("hidden");
          placeholderTextContainer.classList.remove("hidden");
        }
      } else {
        const selectedOption = selectedValueContainer.querySelector(
          "input[type='radio']"
        );
        const option = selectContainer.querySelector(
          `.option[data-value='${selectedOption.value}']`
        );
        option.classList.remove("checked");
        selectedValueContainer.textContent = "";
        placeholderTextContainer.classList.remove("hidden");
        selectedValueContainer.classList.add("hidden");
        replaceIcon();
      }
    }
  });

  // Close dropdown when clicked outside
  document.addEventListener("click", (e) => {
    if (!selectContainer.contains(e.target)) {
      selectContainer.classList.remove("active");
    }
  });
})();
