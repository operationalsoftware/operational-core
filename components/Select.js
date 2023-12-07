(function () {
  const CLOSE_ICON = `<svg class="remove-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z" /></svg>`;

  const selectContainer = me();
  const isMultiple = selectContainer.attribute("data-multiple");
  const multiSelectHiddenInput = me("input[type='hidden']", selectContainer);
  const selectHiddenInput = me(".hidden-input", selectContainer);
  const arrow = me(".arrow", selectContainer);

  // Conditionally set variables
  let multiSelectButton;
  let selectedValuesContainer;
  let optionText;

  if (isMultiple) {
    multiSelectButton = me(".multi-select-button", selectContainer);
    selectedValuesContainer = me(".selected-values", selectContainer);
  } else {
    optionText = me(".select-button span", selectContainer);
  }

  const selectedValues = [];

  function updateDefaultText() {
    const textNode = isMultiple ? multiSelectButton.childNodes[0] : optionText;
    if (selectHiddenInput.value === "" && !isMultiple) {
      optionText.textContent = "Select an option";
    }

    if (isMultiple) {
      if (selectedValues.length > 0) {
        if (textNode.nodeType === 3) {
          multiSelectButton.removeChild(textNode);
        }
      } else {
        multiSelectButton.insertAdjacentHTML("afterbegin", "Select Options");
      }
    }
  }

  function replaceIcon() {
    // Create Remove Icon
    const closeIcon = createIcon(CLOSE_ICON);
    if (selectHiddenInput.value !== "") {
      arrow.classList.remove("arrow");
      arrow.classList.add("remove-icon");
      arrow.innerHTML = closeIcon.outerHTML;
    } else {
      arrow.classList.remove("remove-icon");
      arrow.classList.add("arrow");
      arrow.innerHTML = "";
    }
  }

  selectContainer.on("click", (e) => {
    e.stopPropagation();
    // Simple Select
    if (e.target.matches(".select-button")) {
      selectContainer.classToggle("active");
    }
    if (e.target.matches(".select-dropdown li")) {
      const parent = e.target.children[0];
      if (parent) {
        const label = parent.children[0].textContent;
        const value = parent.children[0].getAttribute("for");
        selectHiddenInput.value = value;
        optionText.textContent = label;
        if (selectContainer) {
          selectContainer.classToggle("active");
        }
        replaceIcon();
      }
    }

    if (!isMultiple) {
      if (e.target.matches(".remove-icon")) {
        selectHiddenInput.value = "";
        optionText.textContent = "Select an option";
        replaceIcon();
      }
    }

    // Multi Select
    if (e.target.matches(".multi-select-button")) {
      selectContainer.classToggle("active");
    }

    if (e.target.matches(".option")) {
      const checkbox = me("input[type='checkbox']", e.target);
      const label = me("label", e.target).textContent;
      const value = checkbox.value;
      checkbox.checked = !checkbox.checked;
      const index = selectedValues.indexOf(value);
      // Create Remove Icon
      const closeIconString = `<svg class="remove-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z" /></svg>`;
      const closeIcon = createIcon(closeIconString);
      // Create Badges
      const createSelectedValue = document.createElement("span");
      createSelectedValue.classList.add("selected-value");
      createSelectedValue.setAttribute("data-value", value);
      createSelectedValue.textContent = label;
      createSelectedValue.appendChild(closeIcon.cloneNode(true));
      if (index === -1) {
        selectedValues.push(value);
        e.target.classList.add("checked");
        e.target.setAttribute("data-value", value);
        selectedValuesContainer.appendChild(createSelectedValue);
      } else {
        e.target.classList.remove("checked");
        e.target.removeAttribute("data-value");
        selectedValues.splice(index, 1);
        const selectedValue = me(`.selected-value[data-value='${value}']`);
        selectedValuesContainer.removeChild(selectedValue);
      }
      multiSelectHiddenInput.value = JSON.stringify(selectedValues);
      updateDefaultText();
    }
  });

  if (selectedValuesContainer) {
    selectedValuesContainer.on("click", (e) => {
      if (e.target.matches(".remove-icon")) {
        const parent = e.target.parentNode;
        const value = parent.getAttribute("data-value");
        const index = selectedValues.indexOf(value);
        const option = me(`.option[data-value='${value}']`);
        option.classList.remove("checked");
        option.removeAttribute("data-value");
        selectedValues.splice(index, 1);
        selectedValuesContainer.removeChild(parent);
        multiSelectHiddenInput.value = JSON.stringify(selectedValues);
        updateDefaultText();
      }
    });
  }

  // Close dropdown when clicked outside
  document.addEventListener("click", (e) => {
    if (!selectContainer.contains(e.target)) {
      selectContainer.classRemove("active");
    }
  });

  updateDefaultText();
})();
