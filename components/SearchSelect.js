(function () {
  const closeIconString = `<svg class="remove-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z" /></svg>`;
  const searchSelect = me();
  const isMultiple = searchSelect.attribute("data-multiple");
  const hiddenInput = me("input[type='hidden']", searchSelect);
  const searchInput = me("input[type='search']", searchSelect);
  const selectedValues = [];
  let arrow;
  let selectedValuesContainer;

  // Conditionally rendered elements
  if (isMultiple) {
    selectedValuesContainer = me(".selected-values", searchSelect);
  } else {
    arrow = me(".arrow", searchSelect);
  }

  function createIcon(iconString) {
    const domParser = new DOMParser();
    const iconDoc = domParser.parseFromString(iconString, "image/svg+xml");
    return iconDoc.documentElement;
  }

  function replaceIcon() {
    const closeIcon = createIcon(closeIconString);
    if (hiddenInput.value !== "") {
      arrow.classList.remove("arrow");
      arrow.classList.add("remove-icon");
      arrow.innerHTML = closeIcon.outerHTML;
    } else {
      arrow.classList.remove("remove-icon");
      arrow.classList.add("arrow");
      arrow.innerHTML = "";
    }
  }

  searchInput.on("keyup", () => {
    if (isMultiple) {
      const defaultWidth = 50;
      const newWidth = defaultWidth + searchInput.value.length;

      searchInput.styles({
        width: `${newWidth + 1}px`,
      });
    }

    // Open the dropdown
    if (searchInput.value.length > 0) {
      searchSelect.classAdd("active");
    } else {
      searchSelect.classRemove("active");
    }
  });

  searchSelect.on("click", (e) => {
    e.stopPropagation();
    if (
      e.target.matches(".content-container") &&
      searchInput.value.length === 0
    ) {
      searchInput.focus();
    } else if (
      e.target.matches(".content-container") &&
      searchInput.value.length > 0
    ) {
      searchInput.focus();
      searchSelect.classAdd("active");
    }

    // Selecting Options
    if (e.target.matches(".option")) {
      if (isMultiple) {
        e.stopPropagation();
        const checkbox = me("input[type='checkbox']", e.target);
        const label = me("label", e.target).textContent;
        const value = checkbox.value;
        checkbox.checked = !checkbox.checked;
        const index = selectedValues.indexOf(value);

        const closeIcon = createIcon(closeIconString);
        // Create Badges
        const createSelectedValue = document.createElement("div");
        createSelectedValue.classList.add("selected-value");
        createSelectedValue.setAttribute("data-value", value);
        const createSelectedValueLabel = document.createElement("span");
        createSelectedValueLabel.textContent = label;
        createSelectedValue.appendChild(createSelectedValueLabel);
        createSelectedValue.appendChild(closeIcon.cloneNode(true));

        if (index === -1) {
          selectedValues.push(value);
          e.target.classList.add("checked");
          e.target.setAttribute("data-value", value);
          selectedValuesContainer.appendChild(createSelectedValue);
          searchInput.value = "";
        } else {
          e.target.classList.remove("checked");
          e.target.removeAttribute("data-value");
          selectedValues.splice(index, 1);
          const selectedValue = me(`.selected-value[data-value='${value}']`);
          selectedValuesContainer.removeChild(selectedValue);
        }
        hiddenInput.value = JSON.stringify(selectedValues);
      } else {
        e.stopPropagation();
        const checkbox = me("input[type='checkbox']", e.target);
        const label = me("label", e.target).textContent;
        const value = checkbox.value;
        checkbox.checked = !checkbox.checked;

        hiddenInput.value = value;
        replaceIcon();
        searchInput.value = label;
        searchSelect.classRemove("active");
      }
    }

    // Remove value for single select
    if (!isMultiple) {
      if (e.target.matches(".remove-icon")) {
        hiddenInput.value = "";
        searchInput.value = "";
        replaceIcon();
      }
    }
  });

  if (isMultiple) {
    // Remove selected values
    on(me(selectedValuesContainer), "click", (e) => {
      if (e.target.matches(".remove-icon")) {
        const value = e.target.parentNode.getAttribute("data-value");
        const index = selectedValues.indexOf(value);
        const option = me(`.option[data-value='${value}']`);
        option.classRemove("checked");
        selectedValues.splice(index, 1);
        selectedValuesContainer.removeChild(e.target.parentNode);
        hiddenInput.value = JSON.stringify(selectedValues);
      }
    });
  }

  // Close dropdown when clicked outside
  document.addEventListener("click", (e) => {
    if (!searchSelect.contains(e.target)) {
      searchSelect.classRemove("active");
    }
  });
})();
