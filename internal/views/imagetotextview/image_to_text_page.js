(() => {
  const returnUrlInput = document.querySelector("#return-url");
  const patternInput = document.querySelector("#regex-pattern");
  const flagsInput = document.querySelector("#regex-flags");
  const fileInput = document.querySelector("#ocr-file");
  const resetButton = document.querySelector("#ocr-reset");
  const previewImg = document.querySelector("#ocr-preview");
  const statusEl = document.querySelector("#ocr-status");
  const textOutput = document.querySelector("#ocr-text");
  const errorEl = document.querySelector("#ocr-error");

  if (!returnUrlInput || !patternInput || !fileInput) {
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const initialReturnUrl =
    params.get("return_to") || params.get("returnTo") || "";

  if (initialReturnUrl) {
    returnUrlInput.value = initialReturnUrl;
  }


  let workerPromise = null;
  let activeObjectUrl = null;

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

  const getWorker = async () => {
    if (!window.Tesseract || !window.Tesseract.createWorker) {
      throw new Error("Tesseract.js failed to load.");
    }
    if (!workerPromise) {
      workerPromise = window.Tesseract.createWorker("eng");
    }
    return workerPromise;
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

    setStatus("Loading OCR engine...");

    try {
      const worker = await getWorker();
      setStatus("Extracting text...");
      const result = await worker.recognize(file);
      const text = (result && result.data && result.data.text) || "";

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
    });
  }

  fileInput.addEventListener("change", (event) => {
    const file = event.target.files && event.target.files[0];
    handleFile(file);
  });
})();
