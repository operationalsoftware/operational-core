(function () {
  const selectContainer = me(".custom-select");
  const select = me(".select-button");
  const optionText = me(".select-button span");
  const hiddenInput = me(".hidden-input");
  const arrow = me(".arrow", selectContainer);

  function updateDefaultText() {
    if (hiddenInput.value === "") {
      optionText.textContent = "Select an option";
    }
  }

  function replaceIcon() {
    // Create Remove Icon
    const closeIconString = `<svg class="remove-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z" /></svg>`;
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

  updateDefaultText();

  function setAriaAttribute() {
    select.setAttribute(
      "aria-expanded",
      select.getAttribute("aria-expanded") === "true" ? "false" : "true"
    );
  }

  selectContainer.on("click", (e) => {
    e.stopPropagation();

    if (e.target.matches(".select-button")) {
      setAriaAttribute();
      selectContainer.classToggle("active");
    }

    if (e.target.matches(".select-dropdown li")) {
      const parent = e.target.children[0];
      if (parent) {
        const label = parent.children[0].textContent;
        const value = parent.children[0].getAttribute("for");
        hiddenInput.value = value;
        optionText.textContent = label;
        selectContainer.classToggle("active");
        replaceIcon();
        setAriaAttribute();
      }
    }
  });

  select.on("click", (e) => {
    if (e.target.matches(".remove-icon")) {
      hiddenInput.value = "";
      optionText.textContent = "Select an option";
      replaceIcon();
      setAriaAttribute();
    }
  });

  me("body").on("click", () => {
    selectContainer.classRemove("active");
    setAriaAttribute();
  });
})();
