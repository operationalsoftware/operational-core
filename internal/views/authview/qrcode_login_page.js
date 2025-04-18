const qrcodeForm = document.getElementById("qrcode-login-form");
const formInput = document.querySelector(".qrcode-form-input");
// Cancel button
const cancelButton = document.querySelector(".cancel-button");

// Get Query params
const urlParams = new URLSearchParams(window.location.search);
const encryptedCredentials = urlParams.get("EncryptedCredentials");

document.addEventListener("DOMContentLoaded", function () {
  if (encryptedCredentials) {
    const url = new URL(window.location);
    url.searchParams.delete("EncryptedCredentials");
    window.history.replaceState({}, document.title, url);

    qrcodeForm.submit();
  }
});

// Handle input typing in qrcode form
let typingTimer;
const doneTypingInterval = 500; // ms

formInput.addEventListener("input", (e) => {
  clearTimeout(typingTimer);
  typingTimer = setTimeout(() => {
    qrcodeForm.submit();
  }, doneTypingInterval);
});

cancelButton.addEventListener("click", () => {
  window.location.href = "/auth/password";
});
