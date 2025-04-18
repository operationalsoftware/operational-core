const passwordGenerateBtn = document.querySelector(".generate-password-btn");
const passwordInput = document.querySelector('[name="Password"]');
const confirmPasswordInput = document.querySelector('[name="ConfirmPassword"]');

// Get query params
const urlParams = new URLSearchParams(window.location.search);
const encryptedCredentials = urlParams.get("EncryptedCredentials");

document.addEventListener("DOMContentLoaded", async () => {
  if (encryptedCredentials) {
    // NFC buttons
    const writeNfcBtn = document.querySelector(".write-nfc-btn");
    const readonlyNfcBtn = document.querySelector(".nfc-readonly-btn");

    writeNfcBtn.disabled = true;
    readonlyNfcBtn.disabled = true;

    const hasNFC = await AppNFC.checkDeviceHasNFC();
    if (!hasNFC) {
      alert("This device does not support NFC");
      return;
    }

    writeNfcBtn.disabled = false;

    writeNfcBtn.addEventListener("click", async () => {
      try {
        await AppNFC.writeNFCTag({
          recordType: "text",
          tagContent: encryptedCredentials,
        });
        alert("NFC tag written successfully");
        writeNfcBtn.disabled = true;
        readonlyNfcBtn.disabled = false;
      } catch (e) {
        alert(e.message);
      }
    });

    readonlyNfcBtn.addEventListener("click", async () => {
      AppNFC.makeNFCReadOnly()
        .then(() => {
          alert("NFC tag made read-only successfully");
        })
        .catch((e) => {
          alert(e.message);
        });
    });
  }
});

// Random password generation
passwordGenerateBtn.addEventListener("click", () => {
  const randPassword = generateRandomPassword(12);
  passwordInput.value = randPassword;
  confirmPasswordInput.value = randPassword;
});

document.querySelectorAll(".toggle-password-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    const target = btn.dataset.target;
    const input = document.querySelector(`[name="${target}"]`);
    const eyeOpen = btn.querySelector(".eye-open-icon");
    const eyeClosed = btn.querySelector(".eye-closed-icon");

    if (input.type === "password") {
      input.type = "text";
      eyeOpen.classList.add("hidden");
      eyeClosed.classList.remove("hidden");
    } else {
      input.type = "password";
      eyeOpen.classList.remove("hidden");
      eyeClosed.classList.add("hidden");
    }
  });
});

function generateRandomPassword(length) {
  if (length < 8) {
    throw new Error("Password length must be at least 8 characters");
  }

  const lowercase = "abcdefghijklmnopqrstuvwxyz";
  const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
  const numbers = "0123456789";
  const symbols = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~";
  const allChars = lowercase + uppercase + numbers + symbols;

  const getRandomChar = (charset) =>
    charset[Math.floor(Math.random() * charset.length)];

  // Ensure at least one of each required character type
  let password = [
    getRandomChar(lowercase),
    getRandomChar(uppercase),
    getRandomChar(numbers),
    getRandomChar(symbols),
  ];

  // Fill the rest of the password with random characters
  for (let i = 4; i < length; i++) {
    password.push(getRandomChar(allChars));
  }

  // Shuffle the password
  for (let i = password.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [password[i], password[j]] = [password[j], password[i]];
  }

  return password.join("");
}
