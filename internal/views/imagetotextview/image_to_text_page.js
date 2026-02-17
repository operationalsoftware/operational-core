(() => {
  const returnUrlInput = document.querySelector("#return-url");
  const patternInput = document.querySelector("#regex-pattern");
  const flagsInput = document.querySelector("#regex-flags");
  const fileInput = document.querySelector("#ocr-file");
  const resetButton = document.querySelector("#ocr-reset");
  const previewImg = document.querySelector("#ocr-preview");
  const previewProcessedImg = document.querySelector("#ocr-preview-processed");
  const statusEl = document.querySelector("#ocr-status");
  const textOutput = document.querySelector("#ocr-text");
  const errorEl = document.querySelector("#ocr-error");
  const ocrClient = window.OperationalOcr;

  if (!returnUrlInput || !patternInput || !fileInput || !ocrClient) {
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const initialReturnUrl =
    params.get("return_to") || params.get("returnTo") || "";

  if (initialReturnUrl) {
    returnUrlInput.value = initialReturnUrl;
  }

  let activeObjectUrl = null;
  let activeProcessedUrl = null;

  const setStatus = (message) => {
    if (statusEl) {
      statusEl.textContent = message;
    }
  };

  const setError = (message) => {
    if (!errorEl) return;
    errorEl.textContent = message || "";
  };

  const clearResults = () => {
    if (textOutput) {
      textOutput.value = "";
    }
    setError("");
  };

  const parseRegex = () => {
    const pattern = patternInput.value.trim();
    if (!pattern) {
      throw new Error("Regex pattern is required.");
    }
    const flags = flagsInput ? flagsInput.value.trim() : "";
    try {
      return new RegExp(pattern, flags);
    } catch (err) {
      throw new Error("Invalid regex pattern or flags.");
    }
  };

  const buildReturnUrl = (groups) => {
    const base = returnUrlInput.value.trim();
    if (!base) {
      throw new Error("Return URL is required.");
    }

    let url;
    try {
      url = new URL(base, window.location.origin);
    } catch (err) {
      throw new Error("Return URL is not valid.");
    }

    Object.entries(groups).forEach(([key, value]) => {
      if (value == null) return;
      const cleaned = String(value).trim();
      if (!cleaned) return;
      url.searchParams.set(key, cleaned);
    });

    return url.toString();
  };

  const handleFile = async (file) => {
    clearResults();

    if (!file) return;

    if (activeObjectUrl) {
      URL.revokeObjectURL(activeObjectUrl);
      activeObjectUrl = null;
    }

    activeObjectUrl = URL.createObjectURL(file);
    if (previewImg) {
      previewImg.src = activeObjectUrl;
      previewImg.classList.remove("is-hidden");
    }
    if (previewProcessedImg) {
      previewProcessedImg.src = "";
      previewProcessedImg.classList.add("is-hidden");
    }

    setStatus("Preparing image...");

    try {
      setStatus("Preparing image...");
      const prepared = await ocrClient.preprocessImage(file);

      if (activeProcessedUrl) {
        URL.revokeObjectURL(activeProcessedUrl);
        activeProcessedUrl = null;
      }
      activeProcessedUrl = URL.createObjectURL(prepared.blob);
      if (previewProcessedImg) {
        previewProcessedImg.src = activeProcessedUrl;
        previewProcessedImg.classList.remove("is-hidden");
      }

      setStatus("Extracting text...");
      const { text, refinedBlob } = await ocrClient.recognizeWithBoxRefine(
        prepared.blob,
        {
          firstPassPSM: "11",
          secondPassPSM: "6",
          minConfidence: 45,
          margin: 16,
          imageWidth: prepared.width,
          imageHeight: prepared.height,
        }
      );

      if (refinedBlob && previewProcessedImg) {
        if (activeProcessedUrl) {
          URL.revokeObjectURL(activeProcessedUrl);
          activeProcessedUrl = null;
        }
        activeProcessedUrl = URL.createObjectURL(refinedBlob);
        previewProcessedImg.src = activeProcessedUrl;
        previewProcessedImg.classList.remove("is-hidden");
      }

      if (textOutput) {
        textOutput.value = text.trim();
      }

      if (!text.trim()) {
        throw new Error("No text detected. Try a clearer photo.");
      }

      const regex = parseRegex();
      const match = regex.exec(text);

      if (!match || !match.groups) {
        throw new Error("No named capture groups matched the OCR text.");
      }

      const groups = {};
      Object.entries(match.groups).forEach(([key, value]) => {
        if (value == null) return;
        const cleaned = String(value).trim();
        if (!cleaned) return;
        groups[key] = cleaned;
      });

      if (!Object.keys(groups).length) {
        throw new Error("Named capture groups were empty.");
      }

      const returnUrl = buildReturnUrl(groups);

      setStatus("Captured text. Returning to form...");

      setTimeout(() => {
        window.location.assign(returnUrl);
      }, 800);
    } catch (err) {
      const message = err && err.message ? err.message : "OCR failed.";
      setStatus("Ready.");
      setError(message);
    }
  };

  if (resetButton) {
    resetButton.addEventListener("click", () => {
      fileInput.value = "";
      clearResults();
      setStatus("Ready.");
      if (previewImg) {
        previewImg.src = "";
        previewImg.classList.add("is-hidden");
      }
      if (previewProcessedImg) {
        previewProcessedImg.src = "";
        previewProcessedImg.classList.add("is-hidden");
      }
      if (activeObjectUrl) {
        URL.revokeObjectURL(activeObjectUrl);
        activeObjectUrl = null;
      }
      if (activeProcessedUrl) {
        URL.revokeObjectURL(activeProcessedUrl);
        activeProcessedUrl = null;
      }
    });
  }

  fileInput.addEventListener("change", (event) => {
    const file = event.target.files && event.target.files[0];
    handleFile(file);
  });
})();
