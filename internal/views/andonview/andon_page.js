function updateAndon(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const andonId = targetBtn.dataset.id;
  const andonAction = targetBtn.dataset.action;

  confirmUpdate = confirm(
    `Are you sure you want to ${andonAction} Andon \u2013 ${andonId}?`
  );

  if (confirmUpdate) {
    fetch(`/andons/${andonId}/${andonAction}/update`, {
      method: "POST",
    }).then((res) => {
      if (res.ok) {
        window.location.href = `/andons/${andonId}`;
      } else {
        alert("Failed to update Andon.");
      }
    });
  }
}

async function submitComment(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);

  const entityId = formData.get("EntityID");
  const comment = formData.get("Comment");
  const files = formData.getAll("files");

  // Add comment
  const commentRes = await fetch(`/andons/${entityId}/comments/add`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ comment }),
  });
  const { commentId } = await commentRes.json();

  for (const file of files) {
    const metaRes = await fetch(`/andons/comment/${commentId}/attachment`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        filename: file.name,
        contentType: file.type,
        sizeBytes: file.size,
        entity: "comment",
        entityId: commentId,
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
