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
  let activeProcessedUrl = null;

  const clamp = (value, min, max) => Math.max(min, Math.min(max, value));

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
      workerPromise = window.Tesseract.createWorker("eng", 1, {
        logger: (message) => console.log(message),
      });
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

  const canvasToBlob = (canvas) =>
    new Promise((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (!blob) {
          reject(new Error("Unable to prepare image."));
          return;
        }
        resolve(blob);
      }, "image/png");
    });

  const autoCropCanvas = (sourceCanvas) => {
    const ctx = sourceCanvas.getContext("2d", { willReadFrequently: true });
    if (!ctx) return sourceCanvas;

    const width = sourceCanvas.width;
    const height = sourceCanvas.height;
    const imageData = ctx.getImageData(0, 0, width, height);
    const data = imageData.data;

    let minX = width;
    let minY = height;
    let maxX = 0;
    let maxY = 0;
    let inkPixels = 0;

    for (let y = 0; y < height; y++) {
      const rowStart = y * width * 4;
      for (let x = 0; x < width; x++) {
        const idx = rowStart + x * 4;
        const gray = data[idx];
        if (gray < 200) {
          inkPixels += 1;
          if (x < minX) minX = x;
          if (x > maxX) maxX = x;
          if (y < minY) minY = y;
          if (y > maxY) maxY = y;
        }
      }
    }

    if (inkPixels < width * height * 0.001) {
      return sourceCanvas;
    }

    const marginX = Math.floor(width * 0.03);
    const marginY = Math.floor(height * 0.03);
    minX = Math.max(0, minX - marginX);
    maxX = Math.min(width - 1, maxX + marginX);
    minY = Math.max(0, minY - marginY);
    maxY = Math.min(height - 1, maxY + marginY);

    const cropWidth = maxX - minX + 1;
    const cropHeight = maxY - minY + 1;

    if (cropWidth <= 0 || cropHeight <= 0) {
      return sourceCanvas;
    }

    const cropped = document.createElement("canvas");
    cropped.width = cropWidth;
    cropped.height = cropHeight;
    const croppedCtx = cropped.getContext("2d");
    if (!croppedCtx) return sourceCanvas;

    croppedCtx.drawImage(
      sourceCanvas,
      minX,
      minY,
      cropWidth,
      cropHeight,
      0,
      0,
      cropWidth,
      cropHeight
    );

    return cropped;
  };

  const preprocessImage = async (file) => {
    const imageBitmap = await createImageBitmap(file);
    
    const maxWidth = 1600;
    const scale =
      imageBitmap.width > maxWidth ? maxWidth / imageBitmap.width : 1;
    const targetWidth = Math.max(1, Math.round(imageBitmap.width * scale));
    const targetHeight = Math.max(1, Math.round(imageBitmap.height * scale));

    const canvas = document.createElement("canvas");
    canvas.width = targetWidth;
    canvas.height = targetHeight;

    const ctx = canvas.getContext("2d", { willReadFrequently: true });
    if (!ctx) {
      throw new Error("Unable to prepare image.");
    }

    ctx.fillStyle = "#fff";
    ctx.fillRect(0, 0, targetWidth, targetHeight);
    ctx.imageSmoothingEnabled = true;
    ctx.imageSmoothingQuality = "high";
    ctx.drawImage(imageBitmap, 0, 0, targetWidth, targetHeight);

    const imageData = ctx.getImageData(0, 0, targetWidth, targetHeight);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
      const r = data[i];
      const g = data[i + 1];
      const b = data[i + 2];
      const gray = Math.round(0.299 * r + 0.587 * g + 0.114 * b);
      data[i] = gray;
      data[i + 1] = gray;
      data[i + 2] = gray;
    }

    const contrastFactor = 1.3;
    for (let i = 0; i < data.length; i += 4) {
      const centered = data[i] - 128;
      const boosted = clamp(128 + centered * contrastFactor, 0, 255);
      data[i] = boosted;
      data[i + 1] = boosted;
      data[i + 2] = boosted;
    }

    ctx.putImageData(imageData, 0, 0);

    const croppedCanvas = autoCropCanvas(canvas);
    const blob = await canvasToBlob(croppedCanvas);

    return {
      blob,
      width: croppedCanvas.width,
      height: croppedCanvas.height,
    };
  };

  const computeMeanConfidence = (result) => {
    if (!result || !result.data) return 0;
    const words = result.data.words || [];
    const valid = words.filter(
      (w) => typeof w.confidence === "number" && w.confidence >= 0
    );
    if (valid.length > 0) {
      const sum = valid.reduce((acc, w) => acc + w.confidence, 0);
      return sum / valid.length;
    }
    const lines = result.data.lines || [];
    const lineValid = lines.filter(
      (l) => typeof l.confidence === "number" && l.confidence >= 0
    );
    if (lineValid.length > 0) {
      const sum = lineValid.reduce((acc, l) => acc + l.confidence, 0);
      return sum / lineValid.length;
    }
    return 0;
  };

  const pickBestResult = (primary, fallback) => {
    const primaryText = (primary && primary.data && primary.data.text) || "";
    const fallbackText = (fallback && fallback.data && fallback.data.text) || "";
    const primaryConfidence = computeMeanConfidence(primary);
    const fallbackConfidence = computeMeanConfidence(fallback);

    if (!primaryText.trim() && fallbackText.trim()) {
      return fallback;
    }

    if (fallbackConfidence > primaryConfidence + 2) {
      return fallback;
    }

    if (primaryConfidence === 0 && fallbackConfidence === 0) {
      return fallbackText.length > primaryText.length ? fallback : primary;
    }

    return primary;
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
      const prepared = await preprocessImage(file);

      if (activeProcessedUrl) {
        URL.revokeObjectURL(activeProcessedUrl);
        activeProcessedUrl = null;
      }
      activeProcessedUrl = URL.createObjectURL(prepared.blob);
      if (previewProcessedImg) {
        previewProcessedImg.src = activeProcessedUrl;
        previewProcessedImg.classList.remove("is-hidden");
      }

      const worker = await getWorker();
      await worker.setParameters({
        preserve_interword_spaces: "1",
        tessedit_pageseg_mode: "6",
      });

      setStatus("Extracting text...");
      const primaryResult = await worker.recognize(prepared.blob);
      let finalResult = primaryResult;

      const primaryText =
        (primaryResult && primaryResult.data && primaryResult.data.text) || "";
      const primaryConfidence = computeMeanConfidence(primaryResult);

      if (!primaryText.trim() || primaryConfidence < 45) {
        await worker.setParameters({ tessedit_pageseg_mode: "11" }); // Sparse text. May improve detection of small amounts of text or single lines.
        const fallbackResult = await worker.recognize(prepared.blob);
        finalResult = pickBestResult(primaryResult, fallbackResult);
      }

      const text = (finalResult && finalResult.data && finalResult.data.text) || "";

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
