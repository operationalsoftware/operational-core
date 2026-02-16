(() => {
  const ocrClient = window.OperationalOcr;
  if (!ocrClient) return;

  const containers = document.querySelectorAll("[data-ocr-input]");
  if (!containers.length) return;

  const setStatus = (el, message) => {
    if (el) {
      el.textContent = message || "";
    }
  };

  const setError = (el, message) => {
    if (el) {
      el.textContent = message || "";
    }
  };

  const updateInputValue = (input, value) => {
    input.value = value;
    input.dispatchEvent(new Event("input", { bubbles: true }));
    input.dispatchEvent(new Event("change", { bubbles: true }));
  };

  const applyReturnValues = () => {
    const url = new URL(window.location.href);
    let changed = false;

    containers.forEach((container) => {
      const input = container.querySelector("[data-ocr-field]");
      if (!input) return;

      const param = container.getAttribute("data-ocr-param") || input.name;
      if (!param) return;

      if (url.searchParams.has(param)) {
        updateInputValue(input, url.searchParams.get(param) || "");
        url.searchParams.delete(param);
        changed = true;
      }
    });

    if (changed) {
      window.history.replaceState({}, "", url.toString());
    }
  };

  const buildReturnUrl = (param) => {
    const url = new URL(window.location.href);
    if (param) {
      url.searchParams.delete(param);
    }
    return url.toString();
  };

  const storeOcrText = (text) => {
    const key = `ocr-text-${Date.now()}-${Math.random()
      .toString(16)
      .slice(2)}`;
    try {
      window.sessionStorage.setItem(key, text);
    } catch (err) {
      return "";
    }
    return key;
  };

  const sendToFixPage = ({ text, pattern, flags, param }) => {
    const storageKey = storeOcrText(text);
    const returnTo = buildReturnUrl(param);

    const fixUrl = new URL("/image-to-text/fix", window.location.origin);
    fixUrl.searchParams.set("return_to", returnTo);
    if (param) {
      fixUrl.searchParams.set("param", param);
    }
    fixUrl.searchParams.set("pattern", pattern);
    if (flags) {
      fixUrl.searchParams.set("flags", flags);
    }
    if (storageKey) {
      fixUrl.searchParams.set("storage", storageKey);
    } else {
      fixUrl.searchParams.set("text", text);
    }

    window.location.assign(fixUrl.toString());
  };

  const handleFile = async ({ file, input, container, statusEl, errorEl }) => {
    setError(errorEl, "");

    if (!file) return;

    const pattern = (container.getAttribute("data-ocr-pattern") || "").trim();
    const flags = (container.getAttribute("data-ocr-flags") || "").trim();
    const param = container.getAttribute("data-ocr-param") || input.name;

    if (!pattern) {
      setError(errorEl, "OCR regex pattern is required.");
      return;
    }

    try {
      setStatus(statusEl, "Preparing image...");
      const prepared = await ocrClient.preprocessImage(file);

      const worker = await ocrClient.getWorker();
      await worker.setParameters({
        preserve_interword_spaces: "1",
        tessedit_pageseg_mode: "6",
      });

      setStatus(statusEl, "Extracting text...");
      const primaryResult = await worker.recognize(prepared.blob);
      let finalResult = primaryResult;

      const primaryText =
        (primaryResult && primaryResult.data && primaryResult.data.text) || "";
      const primaryConfidence = ocrClient.computeMeanConfidence(primaryResult);

      if (!primaryText.trim() || primaryConfidence < 45) {
        await worker.setParameters({ tessedit_pageseg_mode: "11" });
        const fallbackResult = await worker.recognize(prepared.blob);
        finalResult = ocrClient.pickBestResult(primaryResult, fallbackResult);
      }

      const text = (finalResult && finalResult.data && finalResult.data.text) || "";
      if (!text.trim()) {
        throw new Error("No text detected. Try a clearer photo.");
      }

      const regex = ocrClient.parseRegex(pattern, flags);
      const value = ocrClient.extractFirstValue(text, regex);

      if (value) {
        updateInputValue(input, value);
        setStatus(statusEl, "Captured.");
        return;
      }

      setStatus(statusEl, "Needs review...");
      sendToFixPage({ text, pattern, flags, param });
    } catch (err) {
      const message = err && err.message ? err.message : "OCR failed.";
      setStatus(statusEl, "Ready.");
      setError(errorEl, message);
    }
  };

  applyReturnValues();

  containers.forEach((container) => {
    const input = container.querySelector("[data-ocr-field]");
    const trigger = container.querySelector("[data-ocr-trigger]");
    const fileInput = container.querySelector("[data-ocr-file]");
    const statusEl = container.querySelector("[data-ocr-status]");
    const errorEl = container.querySelector("[data-ocr-error]");

    if (!input || !trigger || !fileInput) return;

    trigger.addEventListener("click", () => {
      fileInput.click();
    });

    fileInput.addEventListener("change", (event) => {
      const file = event.target.files && event.target.files[0];
      handleFile({ file, input, container, statusEl, errorEl });
      fileInput.value = "";
    });
  });
})();
