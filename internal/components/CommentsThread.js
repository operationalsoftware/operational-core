(function () {
  // Selected files display (guard if element not present)
  const fileInput = document.querySelector('[name="Files"]');
  if (!fileInput) return;

  fileInput.addEventListener("change", (e) => {
    const container = document.getElementById("selected-files");
    if (!container) return;
    container.innerHTML = "";
    const maxFiles = parseInt(e.target.dataset.maxFiles || "10", 10);
    let files = Array.from(e.target.files);

    if (files.length > maxFiles) {
      alert(`You can only upload up to ${maxFiles} files.`);

      // keep only first allowed files
      files = files.slice(0, maxFiles);

      // reset input to only hold allowed files
      const dataTransfer = new DataTransfer();
      files.forEach((file) => dataTransfer.items.add(file));
      e.target.files = dataTransfer.files;
    }

    if (files.length === 0) {
      container.textContent = "No files selected";
      return;
    }

    files.forEach((file) => {
      const div = document.createElement("div");
      div.textContent = file.name;
      container.appendChild(div);
    });
  });
})();

async function submitComment(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);
  const submitBtn = form.querySelector("button[type='submit']");
  submitBtn.classList.add("loading");
  submitBtn.disabled = true;

  const comment = formData.get("Comment");
  const threadId = form.dataset.threadId;
  const envelopeAttr = form.dataset.hmacEnvelope;
  if (!threadId) {
    console.error("Missing comment thread id on form");
    return;
  }
  if (!envelopeAttr) {
    console.error("Missing HMAC envelope on form");
    return;
  }
  let commentHmac;
  try {
    commentHmac = JSON.parse(envelopeAttr);
  } catch (e) {
    console.error("Invalid HMAC envelope JSON", e);
    return;
  }

  // Filter empty file entries and 0 byte files
  const files = formData
    .getAll("Files")
    .filter((file) => file instanceof File && file.name && file.size > 0);

  // Add comment to centralized endpoint using thread id
  const commentRes = await fetch(`/comments/${threadId}/add`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ comment, hmac: commentHmac }),
  });
  if (!commentRes.ok) {
    submitBtn.classList.remove("loading");
    submitBtn.disabled = false;
    try {
      const err = await commentRes.json();
      alert(
        "Failed to add comment: " +
          (err?.errors ? JSON.stringify(err.errors) : commentRes.statusText)
      );
    } catch (_) {
      alert("Failed to add comment");
    }
    return;
  }
  const { commentId, attachmentHmac } = await commentRes.json();

  for (const file of files) {
    const metaRes = await fetch(`/comments/${commentId}/attachment`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        filename: file.name,
        contentType: file.type,
        sizeBytes: file.size,
        hmac: attachmentHmac,
      }),
    });
    const { fileId, signedUrl } = await metaRes.json();

    // Upload with retry
    await uploadWithRetry(signedUrl, file, 3);

    // confirm upload
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
