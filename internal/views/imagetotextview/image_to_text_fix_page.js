(() => {
  const textArea = document.querySelector("#ocr-fix-text");
  const patternInput = document.querySelector("#ocr-fix-pattern");
  const flagsInput = document.querySelector("#ocr-fix-flags");
  const statusEl = document.querySelector("#ocr-fix-status");
  const errorEl = document.querySelector("#ocr-fix-error");
  const submitBtn = document.querySelector("#ocr-fix-submit");
  const ocrClient = window.OperationalOcr;

  if (!textArea || !patternInput || !submitBtn || !ocrClient) {
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const returnTo = params.get("return_to") || "";
  const param = params.get("param") || "";
  const pattern = params.get("pattern") || "";
  const flags = params.get("flags") || "";
  const storageKey = params.get("storage") || "";
  const fallbackText = params.get("text") || "";

  if (pattern) {
    patternInput.value = pattern;
  }
  if (flags && flagsInput) {
    flagsInput.value = flags;
  }

  let extractedText = "";
  if (storageKey) {
    try {
      extractedText = window.sessionStorage.getItem(storageKey) || "";
    } catch (err) {
      extractedText = "";
    }
  }
  if (!extractedText) {
    extractedText = fallbackText;
  }
  if (extractedText) {
    textArea.value = extractedText;
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
      regex = ocrClient.parseRegex(
        patternInput.value,
        flagsInput ? flagsInput.value : ""
      );
    } catch (err) {
      setError(err && err.message ? err.message : "Invalid regex.");
      return;
    }

    const value = ocrClient.extractFirstValue(text, regex);
    if (!value) {
      setError("No match found. Update the text or regex and try again.");
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
