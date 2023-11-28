(function () {
  const multiSelectContainer = me(".custom-multi-select");
  const selectedValuesContainer = me(".selected-values", multiSelectContainer);
  const multiSelect = me(".multi-select-button");
  const hiddenInput = me("input[type='hidden']", multiSelectContainer);
  const selectedValues = [];

  function updateDefaultText() {
    const textNode = multiSelect.childNodes[0];
    if (selectedValues.length > 0) {
      if (textNode.nodeType === 3) {
        multiSelect.removeChild(textNode);
      }
    } else {
      multiSelect.insertAdjacentHTML("afterbegin", "Select Options");
    }
  }

  updateDefaultText();

  multiSelectContainer.on("click", (e) => {
    e.stopPropagation();
    if (e.target.matches(".multi-select-button")) {
      setAriaAttribute(multiSelect);
      multiSelectContainer.classToggle("active");
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
      hiddenInput.value = JSON.stringify(selectedValues);
      updateDefaultText();
    }
  });

  on(me(selectedValuesContainer), "click", (e) => {
    if (e.target.matches(".remove-icon")) {
      const value = e.target.parentNode.getAttribute("data-value");
      const index = selectedValues.indexOf(value);
      const option = me(`.option[data-value='${value}']`);
      option.classRemove("checked");
      selectedValues.splice(index, 1);
      selectedValuesContainer.removeChild(e.target.parentNode);
      hiddenInput.value = JSON.stringify(selectedValues);
      updateDefaultText();
    }
  });

  me("body").on("click", () => {
    setAriaAttribute(multiSelect);
    multiSelectContainer.classRemove("active");
  });
})();
