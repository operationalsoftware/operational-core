(() => {
  if (window.OperationalOcr) return;

  let workerPromise = null;

  const clamp = (value, min, max) => Math.max(min, Math.min(max, value));

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

  const preprocess = async (imageBitmap) => {
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

  const preprocessImage = async (file) => {
    const imageBitmap = await createImageBitmap(file);
    return await preprocess(imageBitmap);
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

  const parseRegex = (pattern, flags) => {
    const cleaned = (pattern || "").trim();
    if (!cleaned) {
      throw new Error("Regex pattern is required.");
    }
    try {
      return new RegExp(cleaned, (flags || "").trim());
    } catch (err) {
      throw new Error("Invalid regex pattern or flags.");
    }
  };

  const extractNamedGroups = (text, regex) => {
    const match = regex.exec(text || "");
    if (!match || !match.groups) return {};

    const groups = {};
    Object.entries(match.groups).forEach(([key, value]) => {
      if (value == null) return;
      const cleaned = String(value).trim();
      if (!cleaned) return;
      groups[key] = cleaned;
    });

    return groups;
  };

  const extractFirstValue = (text, regex) => {
    const match = regex.exec(text || "");
    if (!match) return "";

    if (match.groups) {
      for (const key of Object.keys(match.groups)) {
        const value = match.groups[key];
        if (value != null && String(value).trim()) {
          return String(value).trim();
        }
      }
    }

    if (match[1] && String(match[1]).trim()) {
      return String(match[1]).trim();
    }

    if (match[0] && String(match[0]).trim()) {
      return String(match[0]).trim();
    }

    return "";
  };

  window.OperationalOcr = {
    getWorker,
    preprocessImage,
    computeMeanConfidence,
    pickBestResult,
    parseRegex,
    extractNamedGroups,
    extractFirstValue,
  };
})();
