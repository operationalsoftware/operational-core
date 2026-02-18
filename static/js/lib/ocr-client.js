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

    // Auto-crop disabled: rely on box-refine crop after first OCR pass.
    const blob = await canvasToBlob(canvas);

    return {
      blob,
      width: canvas.width,
      height: canvas.height,
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

  // Just to avoid duplicate flags in regexp
  const normalizeFlags = (flags, required) => {
    const set = new Set((flags || "").split(""));
    (required || []).forEach((flag) => set.add(flag));
    return Array.from(set).join("");
  };

  const collectMatches = (text, regex, strict) => {
    const values = [];
    if (!regex) return values;
    const flags = normalizeFlags(
      regex.flags,
      strict ? ["g", "m"] : ["g"]
    );
    const source = strict ? `^(?:${regex.source})$` : regex.source;
    const matcher = new RegExp(source, flags);
    let match = matcher.exec(text || "");
    while (match) {
      const value =
        match && match[0] != null && String(match[0]).trim()
          ? String(match[0]).trim()
          : "";
      if (value) {
        values.push(value);
      }
      if (match[0] === "") {
        matcher.lastIndex += 1; // lastIndex doesn't advance on empty match, so move forward manually to avoid infinite loop
      }
      match = matcher.exec(text || "");
    }
    return values;
  };

  const extractBestValue = (text, regex) => {
    const values = collectMatches(text, regex, false);
    if (values.length === 1) return values[0];

    const strictValues = collectMatches(text, regex, true);
    if (strictValues.length >= 1) return strictValues[0];

    return "";
  };

  // Pull bounding boxes from lines (preferred) or words and keep the ones
  // above a confidence threshold to avoid noise.
  const getHighConfidenceBoxes = (result, minConfidence) => {
    if (!result || !result.data) return [];
    const lines = result.data.lines || [];
    const words = result.data.words || [];
    const source = lines.length ? lines : words;
    return source
      .filter(
        (item) =>
          item &&
          item.bbox &&
          typeof item.confidence === "number" &&
          item.confidence >= minConfidence
      )
      .map((item) => item.bbox);
  };

  // Merge bounding boxes into a single rectangle and add a margin.
  const mergeBoxes = (boxes, width, height, margin) => {
    if (!boxes || !boxes.length) return null;
    let minX = Number.POSITIVE_INFINITY;
    let minY = Number.POSITIVE_INFINITY;
    let maxX = 0;
    let maxY = 0;

    boxes.forEach((bbox) => {
      minX = Math.min(minX, bbox.x0);
      minY = Math.min(minY, bbox.y0);
      maxX = Math.max(maxX, bbox.x1);
      maxY = Math.max(maxY, bbox.y1);
    });

    if (!Number.isFinite(minX) || !Number.isFinite(minY)) return null;

    const paddedMinX = Math.max(0, Math.floor(minX - margin));
    const paddedMinY = Math.max(0, Math.floor(minY - margin));
    const paddedMaxX = Math.min(width, Math.ceil(maxX + margin));
    const paddedMaxY = Math.min(height, Math.ceil(maxY + margin));

    const cropWidth = paddedMaxX - paddedMinX;
    const cropHeight = paddedMaxY - paddedMinY;

    if (cropWidth <= 0 || cropHeight <= 0) return null;

    return {
      x: paddedMinX,
      y: paddedMinY,
      width: cropWidth,
      height: cropHeight,
    };
  };

  const cropBlobToBox = async (blob, box) => {
    const bitmap = await createImageBitmap(blob);
    const canvas = document.createElement("canvas");
    canvas.width = box.width;
    canvas.height = box.height;
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      throw new Error("Unable to prepare image.");
    }
    ctx.drawImage(
      bitmap,
      box.x,
      box.y,
      box.width,
      box.height,
      0,
      0,
      box.width,
      box.height
    );
    return await canvasToBlob(canvas);
  };

  // Two-pass OCR using bounding boxes to refine the crop between passes.
  const recognizeWithBoxRefine = async (
    blob,
    {
      firstPassPSM = "11",
      secondPassPSM = "6",
      minConfidence = 45,
      margin = 16,
      imageWidth,
      imageHeight,
    } = {}
  ) => {
    const worker = await getWorker();

    // Pass 1: sparse text mode to discover bounding boxes.
    await worker.setParameters({
      preserve_interword_spaces: "1",
      tessedit_pageseg_mode: firstPassPSM,
    });
    const firstResult = await worker.recognize(blob);

    const firstText =
      (firstResult && firstResult.data && firstResult.data.text) || "";

    const boxes = getHighConfidenceBoxes(firstResult, minConfidence);
    const size = firstResult.data && firstResult.data.imageSize;
    const width = size && size.width ? size.width : imageWidth;
    const height = size && size.height ? size.height : imageHeight;

    if (!width || !height) {
      return {
        text: firstText.trim(),
        result: firstResult,
      };
    }

    const box = mergeBoxes(boxes, width, height, margin);

    if (!box) {
      return {
        text: firstText.trim(),
        result: firstResult,
      };
    }

    // Pass 2: crop to detected region and re-run OCR in block mode.
    const croppedBlob = await cropBlobToBox(blob, box);
    await worker.setParameters({
      preserve_interword_spaces: "1",
      tessedit_pageseg_mode: secondPassPSM,
    });
    const secondResult = await worker.recognize(croppedBlob);
    const secondText =
      (secondResult && secondResult.data && secondResult.data.text) || "";
    const secondConfidence = computeMeanConfidence(secondResult);

    if (secondText.trim() && secondConfidence >= minConfidence) {
      return {
        text: secondText.trim(),
        result: secondResult,
        refinedBlob: croppedBlob,
      };
    }

    return {
      text: firstText.trim(),
      result: firstResult,
      refinedBlob: croppedBlob,
    };
  };

  window.OperationalOcr = {
    getWorker,
    preprocessImage,
    computeMeanConfidence,
    pickBestResult,
    parseRegex,
    extractBestValue,
    recognizeWithBoxRefine,
  };
})();
