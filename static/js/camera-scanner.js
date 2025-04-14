const scannerWrapper = document.getElementById("scanner-wrapper");
const scannerEl = document.getElementById("scanner");
const cameraSelect = document.getElementById("camera-select");
const cancelBtn = document.getElementById("cancel-button");

const LOCAL_STORAGE_KEY = "selectedCamera";
let selectedCamera = localStorage.getItem(LOCAL_STORAGE_KEY);
let html5QrCode;

const urlParams = new URLSearchParams(window.location.search);
const fieldValue = urlParams.get("field");

const onScanSuccess = (data) => {
  try {
    const returnToUrlObj = new URL(document.referrer);
    returnToUrlObj.searchParams.append(fieldValue, data);
    window.location.href = returnToUrlObj.toString();
  } catch (e) {
    console.error("Redirect failed:", e);
  }

  stopScanner();
};

function startScanner(cameraId) {
  html5QrCode = new Html5Qrcode("scanner");

  html5QrCode
    .start(
      { deviceId: { exact: cameraId } },
      {
        fps: 10,
        // qrbox: { width: 250, height: 250 }, // Helps focus the scan area
      },
      (decodedText, decodedResult) => {
        console.log("✅ QR Code Scanned:", decodedText);
        onScanSuccess(decodedText);

        // Optionally stop scanner after successful scan
        // stopScanner();
      },
      (errorMessage) => {
        // scanning failure — can be ignored or logged
        // console.warn("Scan error", errorMessage);
      }
    )
    .catch((err) => {
      console.error("🚫 Scanner start error:", err);
    });
}

function stopScanner() {
  if (html5QrCode) {
    html5QrCode
      .stop()
      .then(() => {
        html5QrCode.clear();
        html5QrCode = null;
      })
      .catch(console.error);
  }
}

async function populateCameras() {
  try {
    const devices = await Html5Qrcode.getCameras();
    cameraSelect.innerHTML = `<option value="">Select Camera</option>`;
    devices.forEach((device) => {
      const option = document.createElement("option");
      option.value = device.id;
      option.textContent = device.label || `Camera ${device.id}`;
      cameraSelect.appendChild(option);
    });

    if (selectedCamera) {
      cameraSelect.value = selectedCamera;
      startScanner(selectedCamera);
    }
  } catch (err) {
    console.error("Failed to get cameras", err);
  }
}

cameraSelect.addEventListener("change", () => {
  selectedCamera = cameraSelect.value;
  localStorage.setItem(LOCAL_STORAGE_KEY, selectedCamera);
  stopScanner();
  if (selectedCamera) {
    startScanner(selectedCamera);
  }
});

cancelBtn.addEventListener("click", () => {
  stopScanner();
  scannerWrapper.style.display = "none";
  window.history.back();
});

// Show scanner and init cameras
function showScanner(fullscreen = true) {
  if (fullscreen) {
    scannerWrapper.classList.add("fullscreen");
  } else {
    scannerWrapper.classList.remove("fullscreen");
  }
  scannerWrapper.style.display = "block";
  populateCameras();
}

// Call this when you want to start scanning
showScanner(true);
