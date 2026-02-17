(() => {
  const ocrClient = window.OperationalOcr;
  if (!ocrClient) return;

  const containers = document.querySelectorAll("[data-ocr-input]");
  if (!containers.length) return;

  const findTargetInput = (name) => {
    if (!name) return null;
    return document.querySelector(`[name="${CSS.escape(name)}"]`);
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
      const param = container.getAttribute("data-ocr-param") || "";
      if (!param || !url.searchParams.has(param)) return;

      const targetName = container.getAttribute("data-ocr-name") || param;
      const input = findTargetInput(targetName);
      if (!input) return;

      updateInputValue(input, url.searchParams.get(param) || "");
      url.searchParams.delete(param);
      changed = true;
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

  const storeOcrText = (text, example) => {
    const key = `ocr-text-${Date.now()}-${Math.random()
      .toString(16)
      .slice(2)}`;
    try {
      const payload = { text };
      if (example) {
        payload.example = example;
      }
      window.sessionStorage.setItem(key, JSON.stringify(payload));
    } catch (err) {
      return "";
    }
    return key;
  };

  const sendToFixPage = ({ text, pattern, flags, param, example }) => {
    const storageKey = storeOcrText(text, example);
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

  const handleFile = async ({ file, container }) => {
    if (!file) return;

    const targetName = container.getAttribute("data-ocr-name") || "";
    const input = findTargetInput(targetName);
    if (!input) return;
    const pattern = (container.getAttribute("data-ocr-pattern") || "").trim();
    const flags = (container.getAttribute("data-ocr-flags") || "").trim();
    const param = container.getAttribute("data-ocr-param") || targetName;
    const example = (container.getAttribute("data-ocr-example") || "").trim();

    if (!pattern) return;

    try {
      const prepared = await ocrClient.preprocessImage(file);

      const { text } = await ocrClient.recognizeWithBoxRefine(prepared.blob, {
        firstPassPSM: "11",
        secondPassPSM: "6",
        minConfidence: 45,
        margin: 16,
        imageWidth: prepared.width,
        imageHeight: prepared.height,
      });
      if (!text.trim()) {
        throw new Error("No text detected. Try a clearer photo.");
      }

      const regex = ocrClient.parseRegex(pattern, flags);
      const value = ocrClient.extractFirstValue(text, regex);

      if (value) {
        updateInputValue(input, value);
        return;
      }

      sendToFixPage({ text, pattern, flags, param, example });
    } catch (err) {
      const message = err && err.message ? err.message : "OCR failed.";
      console.error("OCR error:", message);
    }
  };

  applyReturnValues();

  containers.forEach((container) => {
    const trigger = container.querySelector("[data-ocr-trigger]");
    const fileInput = container.querySelector("[data-ocr-file]");
    if (!trigger || !fileInput) return;

    trigger.addEventListener("click", () => {
      fileInput.click();
    });

    fileInput.addEventListener("change", (event) => {
      const file = event.target.files && event.target.files[0];
      handleFile({ file, container });
      fileInput.value = "";
    });
  });
})();
