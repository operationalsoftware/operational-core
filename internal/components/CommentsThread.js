(function () {
  const commentForms = Array.from(document.querySelectorAll(".comment-form"));
  if (commentForms.length === 0) return;

  commentForms.forEach((form) => {
    initSelectedFiles(form);
    initMentionAutocomplete(form);
  });
})();

function initSelectedFiles(form) {
  const fileInput = form.querySelector('[name="Files"]');
  if (!fileInput) return;

  const container = form.querySelector(".selected-files");
  if (!container) return;

  fileInput.addEventListener("change", (e) => {
    container.innerHTML = "";
    const maxFiles = parseInt(e.target.dataset.maxFiles || "10", 10);
    let files = Array.from(e.target.files || []);

    if (files.length > maxFiles) {
      alert(`You can only upload up to ${maxFiles} files.`);

      files = files.slice(0, maxFiles);

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
}

function initMentionAutocomplete(form) {
  const textarea = form.querySelector(".new-comment");
  const suggestionsEl = form.querySelector(".mention-suggestions");
  if (!textarea || !suggestionsEl) return;

  const endpoint = form.dataset.mentionEndpoint || "/comments/mentions/users";
  let activeIndex = -1;
  let suggestions = [];
  let mentionToken = null;
  let hideTimeout = null;
  let requestController = null;

  const hideSuggestions = () => {
    if (requestController) {
      requestController.abort();
      requestController = null;
    }
    suggestions = [];
    mentionToken = null;
    activeIndex = -1;
    suggestionsEl.hidden = true;
    suggestionsEl.innerHTML = "";
  };

  const renderSuggestions = () => {
    suggestionsEl.innerHTML = "";

    if (suggestions.length === 0) {
      const empty = document.createElement("div");
      empty.className = "mention-suggestion-empty";
      empty.textContent = "No users found";
      suggestionsEl.appendChild(empty);
      suggestionsEl.hidden = false;
      return;
    }

    const fragment = document.createDocumentFragment();

    suggestions.forEach((item, index) => {
      const button = document.createElement("button");
      button.type = "button";
      button.className = "mention-suggestion-item";
      button.dataset.index = String(index);
      if (index === activeIndex) {
        button.classList.add("active");
      }

      const username = document.createElement("span");
      username.className = "mention-suggestion-username";
      username.textContent = `@${item.username}`;

      const display = document.createElement("span");
      display.className = "mention-suggestion-display";
      display.textContent = item.displayName || "";

      button.appendChild(username);
      button.appendChild(display);
      fragment.appendChild(button);
    });

    suggestionsEl.appendChild(fragment);
    suggestionsEl.hidden = false;
  };

  const setActiveIndex = (nextIndex) => {
    if (nextIndex === activeIndex) return;
    activeIndex = nextIndex;

    const buttons = suggestionsEl.querySelectorAll(".mention-suggestion-item");
    buttons.forEach((button) => {
      const buttonIndex = Number.parseInt(button.dataset.index || "-1", 10);
      button.classList.toggle("active", buttonIndex === activeIndex);
    });
  };

  const getSuggestionIndexFromEvent = (event) => {
    const button = event.target.closest(".mention-suggestion-item");
    if (!button) return -1;

    const index = Number.parseInt(button.dataset.index || "-1", 10);
    if (Number.isNaN(index) || !suggestions[index]) return -1;
    return index;
  };

  const getMentionToken = () => {
    const cursorStart = textarea.selectionStart;
    const cursorEnd = textarea.selectionEnd;
    if (cursorStart == null || cursorEnd == null || cursorStart !== cursorEnd) {
      return null;
    }

    const beforeCursor = textarea.value.slice(0, cursorStart);
    const match = beforeCursor.match(/(?:^|[^A-Za-z0-9_])@([A-Za-z0-9_]*)$/);
    if (!match) return null;

    const query = match[1] || "";
    const mentionStart = beforeCursor.length - query.length - 1;

    return {
      query,
      start: mentionStart,
      end: cursorStart,
    };
  };

  const fetchSuggestions = async (query) => {
    if (requestController) {
      requestController.abort();
    }
    requestController = new AbortController();

    const url = `${endpoint}?q=${encodeURIComponent(query)}`;
    try {
      const res = await fetch(url, {
        headers: { "X-Requested-With": "fetch" },
        signal: requestController.signal,
      });

      if (!res.ok) {
        throw new Error(`mention search failed: ${res.status}`);
      }

      const data = await res.json();
      if (Array.isArray(data)) {
        return data;
      }
      return [];
    } catch (err) {
      if (err.name === "AbortError") {
        return null;
      }
      console.error(err);
      return [];
    }
  };

  const applySuggestion = (item) => {
    if (!mentionToken) return;

    const before = textarea.value.slice(0, mentionToken.start);
    const after = textarea.value.slice(mentionToken.end);
    const insertion = `@${item.username} `;

    textarea.value = before + insertion + after;
    const nextCursor = before.length + insertion.length;
    textarea.focus();
    textarea.setSelectionRange(nextCursor, nextCursor);
    hideSuggestions();
  };

  const refreshSuggestions = async () => {
    const token = getMentionToken();
    mentionToken = token;

    if (!token || token.query.length === 0) {
      hideSuggestions();
      return;
    }

    const result = await fetchSuggestions(token.query);
    if (result === null) {
      return;
    }

    const latestToken = getMentionToken();
    if (
      document.activeElement !== textarea ||
      !latestToken ||
      latestToken.query !== token.query ||
      latestToken.start !== token.start ||
      latestToken.end !== token.end
    ) {
      return;
    }

    mentionToken = latestToken;
    suggestions = result;
    activeIndex = suggestions.length > 0 ? 0 : -1;
    renderSuggestions();
  };

  textarea.addEventListener("input", () => {
    if (hideTimeout) {
      clearTimeout(hideTimeout);
      hideTimeout = null;
    }
    refreshSuggestions();
  });

  textarea.addEventListener("keydown", (event) => {
    if (suggestionsEl.hidden) return;

    if (event.key === "ArrowDown") {
      event.preventDefault();
      if (suggestions.length > 0) {
        setActiveIndex((activeIndex + 1) % suggestions.length);
      }
      return;
    }

    if (event.key === "ArrowUp") {
      event.preventDefault();
      if (suggestions.length > 0) {
        setActiveIndex((activeIndex - 1 + suggestions.length) % suggestions.length);
      }
      return;
    }

    if (event.key === "Enter" || event.key === "Tab") {
      if (activeIndex >= 0 && suggestions[activeIndex]) {
        event.preventDefault();
        applySuggestion(suggestions[activeIndex]);
      }
      return;
    }

    if (event.key === "Escape") {
      hideSuggestions();
    }
  });

  textarea.addEventListener("blur", () => {
    hideTimeout = setTimeout(hideSuggestions, 120);
  });

  suggestionsEl.addEventListener("mousedown", (event) => {
    const index = getSuggestionIndexFromEvent(event);
    if (index < 0) return;

    event.preventDefault();
    applySuggestion(suggestions[index]);
  });

  suggestionsEl.addEventListener("mouseover", (event) => {
    const index = getSuggestionIndexFromEvent(event);
    if (index < 0) return;
    setActiveIndex(index);
  });
}

async function submitComment(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);
  const submitBtn = form.querySelector("button[type='submit']");
  submitBtn.classList.add("loading");
  submitBtn.disabled = true;

  const comment = formData.get("Comment");
  const threadId = form.dataset.threadId;
  const hmacEnvelope = form.dataset.hmacEnvelope;
  if (!threadId) {
    console.error("Missing comment thread id on form");
    return;
  }
  if (!hmacEnvelope) {
    console.error("Missing HMAC envelope on form");
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
    body: JSON.stringify({
      comment,
      hmac: hmacEnvelope,
    }),
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
    const metaRes = await fetch(
      `/comments/${threadId}/${commentId}/attachment`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          hmac: attachmentHmac,
          filename: file.name,
          contentType: file.type,
          sizeBytes: file.size,
        }),
      }
    );
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
        console.warn(`Upload failed (status ${res.status}), retrying...`);
      }
    } catch (err) {
      console.warn(`Upload error: ${err.message}, retrying...`);
    }
    attempt++;
  }
  throw new Error(`Failed to upload ${file.name} after ${maxRetries} attempts`);
}
