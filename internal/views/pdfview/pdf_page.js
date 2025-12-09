function handleTemplateNameChange(e) {
  const templateNameSelect = e.target;

  const form = templateNameSelect.closest("form");
  if (!form) return;

  e.preventDefault();

  // Build query string from form data
  const params = new URLSearchParams(new FormData(form));

  // Get current URL without query string
  const baseUrl = window.location.pathname;

  // Redirect to current URL with new query parameters
  window.location.href = `${baseUrl}?${params.toString()}`;
}

function initPrintRequirementSelections() {
  const selects = document.querySelectorAll(".print-requirement-select");
  if (!selects.length) return;

  selects.forEach((select) => {
    const requirement = select.dataset.requirementName || "";
    const defaultPrinter = select.dataset.defaultPrinter || "";
    if (!requirement) return;

    const storageKey = `printRequirement:${requirement}`;
    const savedValue = window.localStorage.getItem(storageKey);

    const hasOption = (value) =>
      Array.from(select.options).some((opt) => opt.value === value);

    if (savedValue && hasOption(savedValue)) {
      select.value = savedValue;
    } else if (defaultPrinter && hasOption(defaultPrinter)) {
      select.value = defaultPrinter;
    }

    select.addEventListener("change", () => {
      window.localStorage.setItem(storageKey, select.value || "");
    });
  });
}

function initPrintActions() {
  const printButtons = document.querySelectorAll(".print-requirement-action");
  printButtons.forEach((btn) => {
    btn.addEventListener("click", async () => {
      const requirement = btn.dataset.requirementName;
      const form = document.querySelector("form");
      const templateSelect = form?.querySelector('select[name="TemplateName"]');
      const inputField = form?.querySelector('textarea[name="InputData"]');
      const printerSelect = btn.parentElement?.querySelector(".print-requirement-select");
      if (!form || !templateSelect || !inputField || !printerSelect) return;

      const templateName = templateSelect.value;
      const inputData = inputField.value;
      const printerID = printerSelect.value;
      const printerName = printerSelect.selectedOptions[0]?.textContent || "";
      if (!templateName) {
        alert("Select a template before printing.");
        return;
      }
      if (!printerID) {
        alert("Select a printer before printing.");
        return;
      }

      const data = new FormData();
      data.append("TemplateName", templateName);
      data.append("InputData", inputData);
      data.append("RequirementName", requirement || templateName);
      data.append("PrinterName", printerName);
      data.append("PrinterID", printerSelect.value);

      try {
        const res = await fetch("/pdf/print", { method: "POST", body: data });
        if (!res.ok) throw new Error("Print failed");
        window.location.reload();
      } catch (err) {
        console.error(err);
        alert("Printing failed. Please try again.");
      }
    });
  });

  const reprintButtons = document.querySelectorAll(".print-log-reprint");
  reprintButtons.forEach((btn) => {
    btn.addEventListener("click", async () => {
      const logId = btn.dataset.printLogId;
      const row = btn.closest(".print-log-action");
      const select = row?.querySelector(".print-log-printer");
      const printerID = select?.value || "";
      const printerName = select?.selectedOptions[0]?.textContent || "";

      const data = new FormData();
      data.append("PrintLogID", logId);
      if (printerID) {
        data.append("PrinterID", printerID);
        data.append("PrinterName", printerName);
      }

      try {
        const res = await fetch("/pdf/print", { method: "POST", body: data });
        if (!res.ok) throw new Error("Reprint failed");
        window.location.reload();
      } catch (err) {
        console.error(err);
        alert("Reprint failed. Please try again.");
      }
    });
  });
}

document.addEventListener("DOMContentLoaded", initPrintRequirementSelections);
document.addEventListener("DOMContentLoaded", initPrintActions);
