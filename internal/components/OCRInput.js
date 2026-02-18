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

  const storeOcrText = (text, regexList) => {
    const key = `ocr-text-${Date.now()}-${Math.random()
      .toString(16)
      .slice(2)}`;
    try {
      const payload = { text };
      if (Array.isArray(regexList) && regexList.length) {
        payload.regexList = regexList;
      }
      window.sessionStorage.setItem(key, JSON.stringify(payload));
    } catch (err) {
      return "";
    }
    return key;
  };

  const sendToFixPage = ({ text, param, regexList }) => {
    const storageKey = storeOcrText(text, regexList);
    const returnTo = buildReturnUrl(param);

    const fixUrl = new URL("/image-to-text/resolve", window.location.origin);
    fixUrl.searchParams.set("ReturnTo", returnTo);
    if (param) {
      fixUrl.searchParams.set("ParamName", param);
    }
    if (storageKey) {
      fixUrl.searchParams.set("Storage", storageKey);
    }

    window.location.assign(fixUrl.toString());
  };

  const handleFile = async ({ file, container }) => {
    if (!file) return;

    const targetName = container.getAttribute("data-ocr-name") || "";
    const input = findTargetInput(targetName);
    if (!input) return;
    const param = container.getAttribute("data-ocr-param") || targetName;
    const regexListRaw = container.getAttribute("data-ocr-regex-list") || "";
    let regexList = [];
    if (regexListRaw) {
      try {
        const parsed = JSON.parse(regexListRaw);
        if (Array.isArray(parsed)) {
          regexList = parsed.filter(
            (item) =>
              item &&
              typeof item.pattern === "string" &&
              item.pattern.trim()
          );
        }
      } catch (err) {
        regexList = [];
      }
    }
    if (!regexList.length) return;

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

      let value = "";
      for (const entry of regexList) {
        const regex = ocrClient.parseRegex(entry.pattern, entry.flags || "");
        const matched = ocrClient.extractBestValue(text, regex);
        if (matched) {
          value = matched;
          break;
        }
      }

      if (value) {
        updateInputValue(input, value);
        return;
      }

      sendToFixPage({ text, param, regexList });
    } catch (err) {
      const message = err && err.message ? err.message : "OCR failed.";
      console.error("OCR error:", message);
    }
  };

  applyReturnValues();

  containers.forEach((container) => {
    if (container.dataset.ocrInitialized === "true") return;
    container.dataset.ocrInitialized = "true";
    const trigger = container.querySelector("[data-ocr-trigger]");
    const fileInput = container.querySelector("[data-ocr-file]");
    if (!trigger || !fileInput) return;

    trigger.addEventListener("click", (event) => {
      event.preventDefault();
      fileInput.click();
    });

    fileInput.addEventListener("change", (event) => {
      const file = event.target.files && event.target.files[0];
      handleFile({ file, container });
      fileInput.value = "";
    });
  });
})();
