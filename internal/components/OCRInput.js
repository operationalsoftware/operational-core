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

  const isMobileCameraFlow = () => {
    const hasTouch = window.matchMedia("(pointer: coarse)").matches;
    const hasMediaCamera =
      !!navigator.mediaDevices && typeof navigator.mediaDevices.getUserMedia === "function";
    return hasTouch && hasMediaCamera && window.isSecureContext;
  };

  const stopStream = (stream) => {
    if (!stream) return;
    stream.getTracks().forEach((track) => track.stop());
  };

  const clamp = (value, min, max) => Math.max(min, Math.min(max, value));

  const canvasToBlob = (canvas) =>
    new Promise((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (!blob) {
          reject(new Error("Unable to capture image."));
          return;
        }
        resolve(blob);
      }, "image/png");
    });

  const createMobileCameraCapture = () => {
    let stream = null;
    let pendingResolve = null;

    const modal = document.createElement("div");
    modal.className = "ocr-camera-modal";
    modal.hidden = true;

    const sheet = document.createElement("div");
    sheet.className = "ocr-camera-sheet";

    const header = document.createElement("div");
    header.className = "ocr-camera-header";
    header.textContent = "Align text inside the rectangle";

    const preview = document.createElement("div");
    preview.className = "ocr-camera-preview";

    const video = document.createElement("video");
    video.className = "ocr-camera-video";
    video.autoplay = true;
    video.playsInline = true;
    video.muted = true;

    const overlay = document.createElement("div");
    overlay.className = "ocr-camera-overlay";

    const cropRect = document.createElement("div");
    cropRect.className = "ocr-camera-crop-rect";

    const actions = document.createElement("div");
    actions.className = "ocr-camera-actions";

    const cancelBtn = document.createElement("button");
    cancelBtn.className = "button secondary";
    cancelBtn.type = "button";
    cancelBtn.textContent = "Cancel";

    const captureBtn = document.createElement("button");
    captureBtn.className = "button";
    captureBtn.type = "button";
    captureBtn.textContent = "Capture";

    const status = document.createElement("div");
    status.className = "ocr-camera-status";
    status.textContent = "";

    actions.appendChild(cancelBtn);
    actions.appendChild(captureBtn);

    overlay.appendChild(cropRect);
    preview.appendChild(video);
    preview.appendChild(overlay);
    sheet.appendChild(header);
    sheet.appendChild(preview);
    sheet.appendChild(actions);
    sheet.appendChild(status);
    modal.appendChild(sheet);
    document.body.appendChild(modal);

    const close = (value) => {
      modal.hidden = true;
      document.body.classList.remove("ocr-camera-open");
      stopStream(stream);
      stream = null;
      video.srcObject = null;
      status.textContent = "";
      const resolve = pendingResolve;
      pendingResolve = null;
      if (resolve) resolve(value || null);
    };

    const getCropSourceBox = () => {
      const previewRect = preview.getBoundingClientRect();
      const frameRect = cropRect.getBoundingClientRect();
      const sourceWidth = video.videoWidth;
      const sourceHeight = video.videoHeight;
      if (!sourceWidth || !sourceHeight || previewRect.width <= 0 || previewRect.height <= 0) {
        return null;
      }

      const sourceAspect = sourceWidth / sourceHeight;
      const previewAspect = previewRect.width / previewRect.height;
      let scale = 1;
      let offsetX = 0;
      let offsetY = 0;

      if (sourceAspect > previewAspect) {
        scale = previewRect.height / sourceHeight;
        const drawnWidth = sourceWidth * scale;
        offsetX = (drawnWidth - previewRect.width) / 2;
      } else {
        scale = previewRect.width / sourceWidth;
        const drawnHeight = sourceHeight * scale;
        offsetY = (drawnHeight - previewRect.height) / 2;
      }

      const frameX = frameRect.left - previewRect.left;
      const frameY = frameRect.top - previewRect.top;
      const frameW = frameRect.width;
      const frameH = frameRect.height;

      const sx = clamp(Math.round((frameX + offsetX) / scale), 0, sourceWidth - 1);
      const sy = clamp(Math.round((frameY + offsetY) / scale), 0, sourceHeight - 1);
      const sw = clamp(Math.round(frameW / scale), 1, sourceWidth - sx);
      const sh = clamp(Math.round(frameH / scale), 1, sourceHeight - sy);

      if (sw <= 0 || sh <= 0) return null;
      return { sx, sy, sw, sh };
    };

    const capture = async () => {
      const cropBox = getCropSourceBox();
      if (!cropBox) {
        throw new Error("Unable to compute crop area.");
      }

      const canvas = document.createElement("canvas");
      canvas.width = cropBox.sw;
      canvas.height = cropBox.sh;
      const ctx = canvas.getContext("2d");
      if (!ctx) {
        throw new Error("Unable to capture image.");
      }

      ctx.drawImage(
        video,
        cropBox.sx,
        cropBox.sy,
        cropBox.sw,
        cropBox.sh,
        0,
        0,
        cropBox.sw,
        cropBox.sh
      );
      return await canvasToBlob(canvas);
    };

    cancelBtn.addEventListener("click", () => close(null));
    modal.addEventListener("click", (event) => {
      if (event.target === modal) {
        close(null);
      }
    });
    document.addEventListener("keydown", (event) => {
      if (!modal.hidden && event.key === "Escape") {
        close(null);
      }
    });
    window.addEventListener("pagehide", () => close(null));

    captureBtn.addEventListener("click", async () => {
      captureBtn.disabled = true;
      status.textContent = "Capturing...";
      try {
        const blob = await capture();
        close(blob);
      } catch (err) {
        status.textContent = "Could not capture image. Try again.";
      } finally {
        captureBtn.disabled = false;
      }
    });

    return {
      open: async () => {
        if (pendingResolve) return await Promise.resolve(null);
        modal.hidden = false;
        document.body.classList.add("ocr-camera-open");
        status.textContent = "Opening camera...";

        const openPromise = new Promise((resolve) => {
          pendingResolve = resolve;
        });

        try {
          stream = await navigator.mediaDevices.getUserMedia({
            video: {
              facingMode: { ideal: "environment" },
              width: { ideal: 1920 },
              height: { ideal: 1080 },
            },
            audio: false,
          });
          video.srcObject = stream;
          await video.play();
          status.textContent = "";
        } catch (err) {
          status.textContent = "Camera access failed.";
          close(null);
        }

        return await openPromise;
      },
    };
  };

  let mobileCameraCapture = null;
  const openMobileCamera = async () => {
    if (!mobileCameraCapture) {
      mobileCameraCapture = createMobileCameraCapture();
    }
    return await mobileCameraCapture.open();
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

  const parseRegexList = (container) => {
    const regexListRaw = container.getAttribute("data-ocr-regex-list") || "";
    let regexList = [];
    if (!regexListRaw) return regexList;
    try {
      const parsed = JSON.parse(regexListRaw);
      if (Array.isArray(parsed)) {
        regexList = parsed.filter(
          (item) => item && typeof item.pattern === "string" && item.pattern.trim()
        );
      }
    } catch (err) {
      regexList = [];
    }
    return regexList;
  };

  const handleImageSource = async ({ source, container }) => {
    if (!source) return;

    const targetName = container.getAttribute("data-ocr-name") || "";
    const input = findTargetInput(targetName);
    if (!input) return;
    const param = container.getAttribute("data-ocr-param") || targetName;
    const regexList = parseRegexList(container);
    if (!regexList.length) return;

    try {
      const prepared = await ocrClient.preprocessImage(source);

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

    trigger.addEventListener("click", async (event) => {
      event.preventDefault();
      if (isMobileCameraFlow()) {
        const capturedBlob = await openMobileCamera();
        if (capturedBlob) {
          handleImageSource({ source: capturedBlob, container });
        }
        return;
      }
      fileInput.click();
    });

    fileInput.addEventListener("change", (event) => {
      const file = event.target.files && event.target.files[0];
      handleImageSource({ source: file, container });
      fileInput.value = "";
    });
  });
})();
