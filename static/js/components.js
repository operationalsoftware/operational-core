/*
 * START: SELECT
 * */
(function () {
  // manage open state
  function toggleOpenSelect(selectEl) {
    selectEl.classList.toggle("open");
  }
  function closeSelect(selectEl) {
    selectEl.classList.remove("open");
  }

  // property getters
  function getSelectIsMultiple(selectEl) {
    return selectEl.getAttribute("data-multiple") === "true";
  }
  function getSelectName(selectEl) {
    return selectEl.getAttribute("data-name");
  }

  // sub element getters
  function getSelectedValuesEl(selectEl) {
    return selectEl.querySelector("span.selected-values");
  }
  function getPlaceholderEl(selectEl) {
    return selectEl.querySelector("span.placeholder");
  }
  function getDropdownArrowEl(selectEl) {
    return selectEl.querySelector("span.dropdown-arrow");
  }
  function getClearSelectEl(selectEl) {
    return selectEl.querySelector("span.clear-select");
  }

  // option value and label getters
  function getOptionValue(optionEl) {
    return optionEl.getAttribute("data-value");
  }
  function getOptionLabel(optionEl) {
    return optionEl.querySelector("span.option-label").innerHTML;
  }

  // selected option chip value getter
  function getChipValue(chipEl) {
    const checkboxEl = chipEl.querySelector("input[type=checkbox]");
    return checkboxEl.value;
  }

  // get an option element by a specified value
  function getOptionElByValue(selectEl, value) {
    let optionEl;
    selectEl.querySelectorAll("div.option").forEach(function (el) {
      const currOptionValue = getOptionValue(el);
      if (currOptionValue === value) optionEl = el;
    });

    return optionEl;
  }

  function getChipElByValue(selectEl, value) {
    let checkboxEl;
    selectEl.querySelectorAll("input[type=checkbox]").forEach(function (cb) {
      if (cb.value === value) checkboxEl = cb;
    });
    const chipEl = checkboxEl.closest("span.selected-value-chip");
    return chipEl;
  }

  function showAsCleared(selectEl) {
    // show placeholder
    const placeholderEl = getPlaceholderEl(selectEl);
    placeholderEl.classList.remove("hidden");

    // hide clear select
    const clearSelectEl = getClearSelectEl(selectEl);
    clearSelectEl.classList.add("hidden");

    // show dropdown arrow
    const dropdownArrowEl = getDropdownArrowEl(selectEl);
    dropdownArrowEl.classList.remove("hidden");
  }

  function showAsHasValue(selectEl) {
    // hide placeholder
    const placeholderEl = getPlaceholderEl(selectEl);
    placeholderEl.classList.add("hidden");

    // show clear select
    const clearSelectEl = getClearSelectEl(selectEl);
    clearSelectEl.classList.remove("hidden");

    // hide dropdown arrow
    const dropdownArrowEl = getDropdownArrowEl(selectEl);
    dropdownArrowEl.classList.add("hidden");
  }

  function getSelectedValue(selectEl) {
    const radioEl = selectEl.querySelector("input[type=radio]");
    if (!radioEl) return null;
    return radioEl.getAttribute("value");
  }
  function getMultiSelectValues(selectEl) {
    const values = [];
    selectEl
      .querySelectorAll("input[type=checkbox]")
      .forEach(function (checkboxEl) {
        const checkboxVal = checkboxEl.getAttribute("value");
        if (checkboxVal) values.push(checkboxVal);
      });
    return values;
  }

  function selectOption(selectEl, optionEl) {
    const isMultiple = getSelectIsMultiple(selectEl);
    const selectName = getSelectName(selectEl);

    // firstly, if we are in single select mode, remove the selected class
    // from all options
    if (!isMultiple) {
      selectEl.querySelectorAll("div.option").forEach(function (optionEl) {
        optionEl.classList.remove("selected");
      });
    }

    // now we can add the selected class to the option element
    optionEl.classList.add("selected");

    const optionValue = getOptionValue(optionEl);
    const optionLabel = getOptionLabel(optionEl);

    const selectedValuesEl = getSelectedValuesEl(selectEl);

    if (isMultiple) {
      // we need to create a chip and append it to the innerHTML of the
      // selected values span
      const chipEl = document.createElement("span");
      chipEl.classList.add("selected-value-chip");
      chipEl.innerHTML = optionLabel;
      // append a span with a remove icon
      const removeEl = document.createElement("span");
      removeEl.classList.add("remove");
      removeEl.textContent = "\u2715"; // small 'clear' icon
      chipEl.appendChild(removeEl);

      // Create a hidden checkbox and append to our chip span
      const checkboxInputEl = document.createElement("input");
      checkboxInputEl.type = "checkbox";
      checkboxInputEl.name = selectName;
      checkboxInputEl.value = optionValue;
      checkboxInputEl.checked = true;
      chipEl.appendChild(checkboxInputEl);

      // finally append out chip
      selectedValuesEl.appendChild(chipEl);
    } else {
      selectedValuesEl.innerHTML = optionLabel;
      const radioInputEl = document.createElement("input");
      radioInputEl.type = "radio";
      radioInputEl.name = selectName;
      radioInputEl.value = optionValue;
      radioInputEl.checked = true;

      selectedValuesEl.appendChild(radioInputEl);
    }

    showAsHasValue(selectEl);

    if (!isMultiple) {
      closeSelect(selectEl);
    }
  }

  function deselectOptionByOptionEl(selectEl, optionEl) {
    console.log("deselect by option el");
    const isMultiple = getSelectIsMultiple(selectEl);

    // firstly, we can remove the selected class from the option element
    optionEl.classList.remove("selected");

    const optionValue = getOptionValue(optionEl);
    const selectedValuesEl = getSelectedValuesEl(selectEl);

    console.log("deselecting value", optionValue);

    if (isMultiple) {
      // find the checkbox element with the same value
      const chipEl = getChipElByValue(selectEl, optionValue);
      console.log("chipEl", chipEl);
      chipEl.remove();
    } else {
      selectedValuesEl.innerHTML = "";
    }

    if (selectedValuesEl.innerHTML === "") {
      showAsCleared(selectEl);
    }
  }

  function deselectOptionByChipEl(selectEl, chipEl) {
    const chipValue = getChipValue(chipEl);

    const optionEl = getOptionElByValue(selectEl, chipValue);
    optionEl.classList.remove("selected");

    chipEl.remove();

    const selectedValuesEl = getSelectedValuesEl(selectEl);
    if (selectedValuesEl.innerHTML === "") {
      showAsCleared(selectEl);
    }
  }

  function clearSelect(selectEl) {
    // deselect all selected options
    selectEl
      .querySelectorAll("div.option.selected")
      .forEach(function (optionEl) {
        optionEl.classList.remove("selected");
      });

    // clear contents of selected values
    const selectedValuesEl = getSelectedValuesEl(selectEl);
    selectedValuesEl.innerHTML = "";

    showAsCleared(selectEl);
  }

  function handleSelectClick(event) {
    const selectEl = event.target.closest("div.select");
    const isMultiple = getSelectIsMultiple(selectEl);

    // handle the click being on an option in the dropdown
    const optionEl = event.target.closest("div.option");
    if (optionEl !== null) {
      const optionValue = optionEl.getAttribute("data-value");

      if (isMultiple) {
        const selectedValues = getMultiSelectValues(selectEl);
        if (!selectedValues.includes(optionValue)) {
          selectOption(selectEl, optionEl);
        } else {
          deselectOptionByOptionEl(selectEl, optionEl);
        }
      } else {
        const selectedValue = getSelectedValue(selectEl);
        if (selectedValue !== optionValue) {
          selectOption(selectEl, optionEl);
        }
      }

      return;
    }

    // handle the select being on the select's 'clear' icon
    const clearSelectEl = event.target.closest("span.clear-select");
    if (clearSelectEl !== null) {
      clearSelect(selectEl);

      return;
    }

    // for multiselect, handle the click being on a selected value chip's remove icon
    const selectedChipRemoveEl = event.target.closest("span.remove");
    if (selectedChipRemoveEl) {
      const chipEl = selectedChipRemoveEl.closest("span.selected-value-chip");
      deselectOptionByChipEl(selectEl, chipEl);

      return;
    }

    // finally, handle the click being on the button, in which case, toggle open
    const buttonEl = event.target.closest("button");
    if (buttonEl) {
      toggleOpenSelect(selectEl);

      return;
    }
  }

  document.addEventListener("click", function (event) {
    const isInsideSelect = event.target.closest("div.select");
    if (isInsideSelect) {
      event.preventDefault();
      event.stopPropagation();
      handleSelectClick(event);
    } else {
      // check if any selects are open and remove the open class if so
      document.querySelectorAll("div.select.open").forEach(function (selectEl) {
        selectEl.classList.remove("open");
      });
    }
  });

  // Toast start
  const toast = document.querySelector(".toast");
  if (toast) {
    toast.style.opacity = "1";
    toast.style.pointerEvents = "auto";
    setTimeout(() => {
      toast.style.opacity = "0";
    }, 4000);
  }
  // Toast end
})();

/*
 * END: SELECT
 * */

/*
 * START TABLE
 * */

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

/*
 * END TABLE
 * */
