function initResourceBulkSelection() {
  const bulkEditButton = document.querySelector(
    "[data-resource-bulk-edit-button]",
  );
  const selectAllCheckbox = document.querySelector(
    "[data-resource-select-all]",
  );
  const checkboxes = document.querySelectorAll("[data-resource-select]");

  if (!bulkEditButton || checkboxes.length === 0) return;

  const storageKey = "resources.bulkEdit.selectedIds";

  const loadSelected = () => {
    try {
      const raw = window.localStorage.getItem(storageKey);
      if (!raw) return [];
      const parsed = JSON.parse(raw);
      if (!Array.isArray(parsed)) return [];
      return parsed.map((id) => String(id));
    } catch (err) {
      return [];
    }
  };

  const saveSelected = (ids) => {
    try {
      window.localStorage.setItem(storageKey, JSON.stringify(ids));
    } catch (err) {
      // Ignore storage errors so selection still works within the page.
    }
  };

  const selected = new Set(loadSelected());

  const syncButtonState = () => {
    bulkEditButton.disabled = selected.size === 0;
  };

  const syncSelectAllState = () => {
    if (!selectAllCheckbox) return;

    let checkedCount = 0;
    checkboxes.forEach((checkbox) => {
      if (checkbox.checked) checkedCount += 1;
    });

    if (checkedCount === 0) {
      selectAllCheckbox.checked = false;
      selectAllCheckbox.indeterminate = false;
      return;
    }

    if (checkedCount === checkboxes.length) {
      selectAllCheckbox.checked = true;
      selectAllCheckbox.indeterminate = false;
      return;
    }

    selectAllCheckbox.checked = false;
    selectAllCheckbox.indeterminate = true;
  };

  checkboxes.forEach((checkbox) => {
    const resourceID = checkbox.getAttribute("data-resource-id");
    if (!resourceID) return;

    if (selected.has(resourceID)) {
      checkbox.checked = true;
    }

    checkbox.addEventListener("change", () => {
      if (checkbox.checked) {
        selected.add(resourceID);
      } else {
        selected.delete(resourceID);
      }
      saveSelected(Array.from(selected));
      syncButtonState();
      syncSelectAllState();
    });
  });

  syncButtonState();
  syncSelectAllState();

  if (selectAllCheckbox) {
    selectAllCheckbox.addEventListener("change", () => {
      checkboxes.forEach((checkbox) => {
        const resourceID = checkbox.getAttribute("data-resource-id");
        if (!resourceID) return;

        checkbox.checked = selectAllCheckbox.checked;
        if (selectAllCheckbox.checked) {
          selected.add(resourceID);
        } else {
          selected.delete(resourceID);
        }
      });
      saveSelected(Array.from(selected));
      syncButtonState();
      syncSelectAllState();
    });
  }

  bulkEditButton.addEventListener("click", () => {
    if (selected.size === 0) return;

    const baseURL =
      bulkEditButton.getAttribute("data-bulk-edit-url") ||
      "/resources/bulk-edit-service-schedules";
    const params = new URLSearchParams();

    Array.from(selected).forEach((resourceID) => {
      params.append("ResourceID", resourceID);
    });

    const query = params.toString();
    window.location.href = query ? `${baseURL}?${query}` : baseURL;
  });
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initResourceBulkSelection);
} else {
  initResourceBulkSelection();
}
