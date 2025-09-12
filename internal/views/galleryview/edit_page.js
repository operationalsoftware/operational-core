document.addEventListener("DOMContentLoaded", () => {
  const galleryItemsUploadBtn = document.querySelector(".upload-gallery-btn");
  const grid = document.querySelector(".gallery-grid");
  if (!grid) return;

  let draggingEl = null;

  // Make items draggable and attach dragstart/dragend via event delegation
  grid.querySelectorAll(".gallery-item").forEach((item) => {
    item.setAttribute("draggable", "true");

    item.addEventListener("dragstart", (e) => {
      draggingEl = item;
      item.classList.add("dragging");

      // Save old position at drag start
      draggingEl.dataset.oldPosition =
        [...grid.querySelectorAll(".gallery-item")].indexOf(item) + 1;

      // Required for Firefox to allow dragging
      try {
        e.dataTransfer.setData("text/plain", item.dataset.id || "");
      } catch (err) {
        // ignore (some browsers throw here if blocked)
      }
      e.dataTransfer.effectAllowed = "move";
    });

    item.addEventListener("dragend", () => {
      if (draggingEl) {
        draggingEl.classList.remove("dragging");

        const items = [...grid.querySelectorAll(".gallery-item")];
        const newPos = items.indexOf(draggingEl) + 1;
        const oldPos = parseInt(draggingEl.dataset.oldPosition, 10);

        saveNewOrder({
          gallery_item_id: parseInt(draggingEl.dataset.id, 10),
          old_position: oldPos,
          new_position: newPos,
        });

        draggingEl = null;
      }
    });
  });

  // Single dragover listener on the grid container
  grid.addEventListener("dragover", (e) => {
    e.preventDefault(); // allow drop

    const pointerX = e.clientX;
    const pointerY = e.clientY;

    const dragging = grid.querySelector(".dragging");
    if (!dragging) return;

    const afterElement = getClosestElement(grid, pointerX, pointerY, dragging);

    if (!afterElement) {
      grid.appendChild(dragging);
    } else {
      grid.insertBefore(dragging, afterElement);
    }
  });

  // Handles file selection
  document.querySelector('[name="Files"]').addEventListener("change", (e) => {
    const container = document.getElementById("selected-files");
    container.innerHTML = "";
    const maxFiles = parseInt(e.target.dataset.maxFiles || "10", 10);
    let files = Array.from(e.target.files);

    if (files.length > maxFiles) {
      alert(`You can only upload up to ${maxFiles} files.`);

      // keep only first 10 files
      files = files.slice(0, maxFiles);

      // reset input to only hold allowed files
      const dataTransfer = new DataTransfer();
      files.forEach((file) => dataTransfer.items.add(file));
      e.target.files = dataTransfer.files;
    }

    if (files.length === 0) {
      galleryItemsUploadBtn.disabled = true;
      container.textContent = "No files selected";
      return;
    }
    galleryItemsUploadBtn.disabled = false;

    files.forEach((file) => {
      const div = document.createElement("div");
      div.textContent = file.name;
      container.appendChild(div);
    });
  });
});

/**
 * Return the child element in container that is closest to the pointer.
 * Skip the currently dragging element.
 * Works for both grid and list layouts by measuring distance to child centers.
 */
function getClosestElement(container, x, y, draggingElement) {
  const children = [...container.querySelectorAll(".gallery-item")].filter(
    (el) => el !== draggingElement
  );

  if (children.length === 0) return null;

  let closest = null;
  let minDistSq = Infinity;

  for (const child of children) {
    const rect = child.getBoundingClientRect();
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;

    const dx = x - cx;
    const dy = y - cy;
    const distSq = dx * dx + dy * dy;

    if (distSq < minDistSq) {
      minDistSq = distSq;
      closest = child;
    }
  }

  // If pointer is below/right of every child, we might want to append.
  // We simply return the closest child — caller may append or insertBefore.
  return closest;
}

async function saveNewOrder(payload) {
  const grid = document.querySelector(".gallery-grid");
  const galleryId = grid.dataset.galleryId;

  const url = `/gallery/${galleryId}/reorder` + window.location.search;

  const res = await fetch(url, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    alert("Failed to save new order");
  }
}

function deleteItem(galleryId, itemId, position) {
  confirmDelete = confirm(
    "Are you sure you want to delete this gallery item ?"
  );

  if (confirmDelete) {
    fetch(
      `/gallery/${galleryId}/item/${itemId}/${position}` +
        window.location.search,
      {
        method: "DELETE",
      }
    ).then(() => location.reload());
  }
}

async function submitGalleryItems(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);
  const submitBtn = form.querySelector("button[type='submit']");
  submitBtn.classList.add("loading");
  submitBtn.disabled = true;

  // Filter empty file entries and 0 byte files
  const files = formData
    .getAll("Files")
    .filter((file) => file instanceof File && file.name && file.size > 0);

  const baseEndpoint = window.location.pathname + window.location.search;

  for (const file of files) {
    const metaRes = await fetch(baseEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        filename: file.name,
        contentType: file.type,
        sizeBytes: file.size,
      }),
    });
    const { fileId, signedUrl } = await metaRes.json();

    await uploadWithRetry(signedUrl, file, 3);

    await fetch(`/files/${fileId}/complete`, { method: "GET" });
  }

  window.location.reload();
}

async function uploadWithRetry(signedUrl, file, maxRetries = 3) {
  const encodedFilename = encodeURIComponent(file.name);

  const isImage = file.type.startsWith("image/");
  const isPDF = file.type === "application/pdf";

  const headers = {
    "Content-Type": file.type || "application/octet-stream",
  };

  if (isImage || isPDF) {
    headers[
      "Content-Disposition"
    ] = `inline; filename*=UTF-8''${encodedFilename}`;
  } else {
    headers[
      "Content-Disposition"
    ] = `attachment; filename*=UTF-8''${encodedFilename}`;
  }

  let attempt = 0;
  while (attempt < maxRetries) {
    try {
      const res = await fetch(signedUrl, {
        method: "PUT",
        headers,
        body: file,
      });

      if (res.ok) {
        return true;
      } else {
        console.warn(`Upload failed (status ${res.status}), retrying…`);
      }
    } catch (err) {
      console.warn(`Upload error: ${err.message}, retrying…`);
    }
    attempt++;
  }
  throw new Error(`Failed to upload ${file.name} after ${maxRetries} attempts`);
}
