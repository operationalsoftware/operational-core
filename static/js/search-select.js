function initSearchSelect(selectEl) {
  const mode = selectEl.dataset.mode;
  const input = selectEl.querySelector(".select-input");
  const inputSpan = selectEl.querySelector(".select-input span");
  const search = selectEl.querySelector(".select-search");
  const dropdown = selectEl.querySelector(".select-dropdown");
  const optionsList = selectEl.querySelector(".select-options");

  const selectedOptionEls = optionsList.querySelectorAll(
    ".select-option.selected"
  );

  let selected;

  if (mode === "multi") {
    selected = Array.from(selectedOptionEls).map((el) => el.dataset.value);
    inputSpan.textContent = selectedOptionEls.length
      ? Array.from(selectedOptionEls)
          .map((el) => el.textContent)
          .join(", ")
      : inputSpan.textContent;
  } else {
    const selectedOption = selectedOptionEls[0];
    if (selectedOption) {
      selected = selectedOption.dataset.value;
      inputSpan.textContent = selectedOption.textContent;
    } else {
      selected = null;
    }
  }

  // Set hidden input initially
  const name = selectEl.dataset.name;
  updateHiddenInputs(selectEl, selected, name);

  input.addEventListener("click", () => {
    dropdown.classList.toggle("open");
    search.focus();
  });

  input.addEventListener("keydown", (e) => {
    if (e.key === "ArrowDown" && !dropdown.classList.contains("open")) {
      e.preventDefault();
      dropdown.classList.add("open");
      search.focus();
    }
  });

  // input.addEventListener("mousedown", () => {
  //   clearTimeout(focusTimeout);
  // });

  search.addEventListener("input", () => {
    const term = search.value.toLowerCase();
    selectEl.dispatchEvent(
      new CustomEvent("load-options", {
        detail: { search: term },
      })
    );

    // OR filter locally:
    optionsList.querySelectorAll(".select-option").forEach((opt) => {
      opt.style.display = opt.textContent.toLowerCase().includes(term)
        ? ""
        : "none";
    });
  });

  optionsList.addEventListener("click", (e) => {
    const option = e.target.closest(".select-option");
    if (!option) return;

    const value = option.dataset.value;

    if (mode === "single") {
      selected = value;
      inputSpan.textContent = option.textContent;
      // input.innerHTML =
      //   option.textContent +
      //   '<svg class="icon"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z"></path></svg></svg>';
      dropdown.classList.remove("open");
    } else {
      const index = selected.indexOf(value);
      if (index > -1) {
        selected.splice(index, 1);
        option.classList.remove("selected");
      } else {
        selected.push(value);
        option.classList.add("selected");
      }
      inputSpan.textContent = selected.join(", ");
    }

    const name = selectEl.dataset.name;
    updateHiddenInputs(selectEl, selected, name);

    selectEl.dispatchEvent(new Event("change", { bubbles: true }));
  });

  let activeIndex = -1;

  function updateActiveOption(index) {
    const options = Array.from(
      optionsList.querySelectorAll(
        ".select-option:not([style*='display: none'])"
      )
    );
    if (options.length === 0) return;

    // Clamp index
    if (index < 0) index = 0;
    if (index >= options.length) index = options.length - 1;

    // Remove previous
    options.forEach((opt) => opt.classList.remove("active"));

    // Add to current
    const current = options[index];
    if (current) {
      current.classList.add("active");
      current.scrollIntoView({ block: "nearest" });
    }

    activeIndex = index;
  }

  search.addEventListener("keydown", (e) => {
    const options = Array.from(
      optionsList.querySelectorAll(
        ".select-option:not([style*='display: none'])"
      )
    );

    if (e.key === "ArrowDown") {
      e.preventDefault();
      updateActiveOption(activeIndex + 1);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      updateActiveOption(activeIndex - 1);
    } else if (e.key === "Enter") {
      e.preventDefault();
      const option = options[activeIndex];
      if (option) {
        option.click();
      }
    } else if (e.key === "Escape") {
      dropdown.classList.remove("open");
    }
  });

  document.addEventListener("click", (e) => {
    if (!selectEl.contains(e.target)) dropdown.classList.remove("open");
  });
}

function updateHiddenInputs(container, selected, inputName) {
  const inputsContainer = container.querySelector(".select-hidden-inputs");
  inputsContainer.innerHTML = "";

  if (Array.isArray(selected)) {
    selected.forEach((value) => {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = inputName;
      input.value = value;
      inputsContainer.appendChild(input);
    });
  } else {
    const input = document.createElement("input");
    input.type = "hidden";
    input.name = inputName;
    input.value = selected;
    inputsContainer.appendChild(input);
  }
}

document.addEventListener("DOMContentLoaded", () =>
  document.querySelectorAll(".search-select").forEach(initSearchSelect)
);
