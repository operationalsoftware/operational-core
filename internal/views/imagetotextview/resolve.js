(() => {
  const tagsContainer = document.querySelector("#ocr-fix-tags");
  const inputEl = document.querySelector("#ocr-fix-input");
  const exampleEl = document.querySelector("#ocr-fix-example");
  const statusEl = document.querySelector("#ocr-fix-status");
  const errorEl = document.querySelector("#ocr-fix-error");
  const backBtn = document.querySelector("#ocr-fix-back");
  const submitBtn = document.querySelector("#ocr-fix-submit");
  const ocrClient = window.OperationalOcr;

  if (!tagsContainer || !inputEl || !submitBtn || !ocrClient) {
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const returnTo = params.get("ReturnTo") || "";
  const param = params.get("ParamName") || "";
  const storageKey = params.get("Storage") || "";

  let extractedText = "";
  let regexList = [];
  if (storageKey) {
    try {
      const stored = window.sessionStorage.getItem(storageKey) || "";
      const parsed = stored ? JSON.parse(stored) : null;
      if (parsed && typeof parsed === "object") {
        extractedText = parsed.text || "";
        if (Array.isArray(parsed.regexList)) {
          regexList = parsed.regexList.filter(
            (item) =>
              item &&
              typeof item.pattern === "string" &&
              item.pattern.trim()
          );
        }
      }
    } catch (err) {
      extractedText = "";
    }
  }

  const sanitizeToken = (token) => {
    const cleaned = token.replace(/^[^A-Za-z0-9]+|[^A-Za-z0-9]+$/g, "");
    if (!cleaned || cleaned === "|") return "";
    return cleaned;
  };

  const normalizeToken = (token) => sanitizeToken(token).toLowerCase();

  const tokenize = (text) =>
    text
      .split(/\s+/)
      .map((token) => sanitizeToken(token))
      .filter(Boolean);

  const uniqueTokens = (tokens) => Array.from(new Set(tokens));

  const updateActiveTags = () => {
    const currentTokens = tokenize(inputEl.value);
    const tokenSet = new Set(currentTokens.map((token) => token.toLowerCase()));
    tagsContainer.querySelectorAll("[data-token]").forEach((tag) => {
      const token = tag.getAttribute("data-token-normalized");
      if (tokenSet.has(token)) {
        tag.classList.add("is-active");
      } else {
        tag.classList.remove("is-active");
      }
    });
  };

  const renderTags = (tokens) => {
    tagsContainer.innerHTML = "";
    if (!tokens.length) {
      const empty = document.createElement("span");
      empty.className = "image-to-text-fix-example";
      empty.textContent = "No tokens detected from OCR.";
      tagsContainer.appendChild(empty);
      return;
    }

    tokens.forEach((token) => {
      const tag = document.createElement("button");
      tag.type = "button";
      tag.className = "image-to-text-fix-tag";
      tag.textContent = token;
      tag.setAttribute("data-token", token);
      tag.setAttribute("data-token-normalized", normalizeToken(token));
      tag.addEventListener("click", () => {
        tag.classList.add("is-animating");
        setTimeout(() => tag.classList.remove("is-animating"), 140);
        const currentTokens = tokenize(inputEl.value);
        const normalized = normalizeToken(token);
        const exists = currentTokens.some(
          (current) => normalizeToken(current) === normalized
        );
        const nextTokens = exists
          ? currentTokens.filter(
              (current) => normalizeToken(current) !== normalized
            )
          : currentTokens.concat(token);
        inputEl.value = nextTokens.join("");
        updateActiveTags();
      });
      tagsContainer.appendChild(tag);
    });

    updateActiveTags();
  };

  renderTags(uniqueTokens(tokenize(extractedText)));

  if (exampleEl) {
    const examples = regexList
      .map((entry) => entry && entry.example)
      .filter((entry) => typeof entry === "string" && entry.trim());
    if (examples.length) {
      exampleEl.textContent = `Examples: ${examples.join(", ")}`;
    }
  }

  const setStatus = (message) => {
    if (statusEl) {
      statusEl.textContent = message || "";
    }
  };

  const setError = (message) => {
    if (errorEl) {
      errorEl.textContent = message || "";
    }
  };

  const buildReturnUrl = (value) => {
    if (!returnTo) {
      throw new Error("Return URL is required.");
    }
    if (!param) {
      throw new Error("Target field is missing.");
    }

    const url = new URL(returnTo, window.location.origin);
    url.searchParams.set(param, value);
    return url.toString();
  };

  submitBtn.addEventListener("click", () => {
    setError("");

    const text = inputEl.value || "";
    if (!text.trim()) {
      setError("Selected text is required.");
      return;
    }

    let value = "";
    try {
      if (!regexList.length) {
        setError("No regex patterns available for matching.");
        return;
      }
      for (const entry of regexList) {
        const regex = ocrClient.parseRegex(`^(?:${entry.pattern})$`, entry.flags || "");
        regex.lastIndex = 0; // reset regex state in case of global flag
        if (regex.test(text)) {
          value = text.trim();
          break;
        }
      }
    } catch (err) {
      setError(err && err.message ? err.message : "Invalid regex.");
      return;
    }
    if (!value) {
      setError("No match found. Update the selected text and try again.");
      return;
    }

    try {
      const returnUrl = buildReturnUrl(value);
      if (storageKey) {
        try {
          window.sessionStorage.removeItem(storageKey);
        } catch (err) {
          // ignore
        }
      }
      setStatus("Returning to form...");
      window.location.assign(returnUrl);
    } catch (err) {
      setError(err && err.message ? err.message : "Unable to return to form.");
    }
  });

  inputEl.addEventListener("input", () => {
    updateActiveTags();
  });

  if (backBtn && returnTo) {
    backBtn.setAttribute("href", returnTo);
  }
})();
