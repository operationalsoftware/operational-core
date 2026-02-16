(() => {
  const textArea = document.querySelector("#ocr-fix-text");
  const exampleEl = document.querySelector("#ocr-fix-example");
  const statusEl = document.querySelector("#ocr-fix-status");
  const errorEl = document.querySelector("#ocr-fix-error");
  const submitBtn = document.querySelector("#ocr-fix-submit");
  const ocrClient = window.OperationalOcr;

  if (!textArea || !submitBtn || !ocrClient) {
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const returnTo = params.get("return_to") || "";
  const param = params.get("param") || "";
  const pattern = params.get("pattern") || "";
  const flags = params.get("flags") || "";
  const storageKey = params.get("storage") || "";
  const fallbackText = params.get("text") || "";

  let extractedText = "";
  let exampleText = "";
  if (storageKey) {
    try {
      const stored = window.sessionStorage.getItem(storageKey) || "";
      const parsed = stored ? JSON.parse(stored) : null;
      if (parsed && typeof parsed === "object") {
        extractedText = parsed.text || "";
        exampleText = parsed.example || "";
      } else {
        extractedText = stored;
      }
    } catch (err) {
      try {
        extractedText = window.sessionStorage.getItem(storageKey) || "";
      } catch (innerErr) {
        extractedText = "";
      }
    }
  }
  if (!extractedText) {
    extractedText = fallbackText;
  }
  if (extractedText) {
    textArea.value = extractedText;
  }
  if (exampleEl && exampleText) {
    exampleEl.textContent = exampleText
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

    const text = textArea.value || "";
    if (!text.trim()) {
      setError("Extracted text is required.");
      return;
    }

    let regex;
    try {
      regex = ocrClient.parseRegex(pattern, flags);
    } catch (err) {
      setError(err && err.message ? err.message : "Invalid regex.");
      return;
    }

    const value = ocrClient.extractFirstValue(text, regex);
    if (!value) {
      setError("No match found. Update the text and try again.");
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
})();
