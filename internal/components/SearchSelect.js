(function () {
  const selectEl = document.currentScript.closest(".search-select");
  if (!selectEl) return;

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

  search.addEventListener("input", async () => {
    const term = search.value.trim();

    selectEl.dispatchEvent(
      new CustomEvent("load-options", {
        detail: { search: term },
      })
    );

    // --- Dynamic fetch logic ---
    const endpoint = selectEl.dataset.optionsEndpoint;
    if (endpoint) {
      const searchQueryParam =
        selectEl.dataset.searchQueryParam || "SearchText";
      const form = selectEl.closest("form");
      const selectedValues = Array.from(
        selectEl.querySelectorAll(".select-hidden-inputs input")
      ).map((i) => i.value);

      const formData = new FormData(form || undefined);
      // send selected under the field's real name, and (optionally) as CSV for convenience
      const fieldName = selectEl.dataset.name;
      const nonEmpty = selectedValues.filter((v) => v != null && v !== "");

      // for compatibility with typical form parsers:
      // - multi: repeat the same key multiple times
      // - single: one key/value
      if (mode === "multi") {
        nonEmpty.forEach((v) => formData.append(fieldName, v));
      } else if (nonEmpty[0]) {
        formData.append(fieldName, nonEmpty[0]);
      }

      formData.append(searchQueryParam, term);

      const url = endpoint + "?" + new URLSearchParams(formData);

      // Abort previous request if still running
      if (search._currentRequest) search._currentRequest.abort();
      search._currentRequest = new AbortController();
      const signal = search._currentRequest.signal;

      try {
        const html = await fetch(url, { signal }).then((r) => r.text());
        optionsList.innerHTML = html;
      } catch (e) {
        if (e.name !== "AbortError") {
          console.error("SearchSelect fetch error:", e);
        }
      }
      return; // skip local filtering if we fetched from server
    }

    // --- Local filtering fallback ---
    const lowerTerm = term.toLowerCase();
    optionsList.querySelectorAll(".select-option").forEach((opt) => {
      opt.style.display = opt.textContent.toLowerCase().includes(lowerTerm)
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

    index = Math.max(0, Math.min(index, options.length - 1));
    options.forEach((opt) => opt.classList.remove("active"));

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
      if (option) option.click();
    } else if (e.key === "Escape") {
      dropdown.classList.remove("open");
    }
  });

  document.addEventListener("click", (e) => {
    if (!selectEl.contains(e.target)) {
      dropdown.classList.remove("open");
    }
  });

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
})();
