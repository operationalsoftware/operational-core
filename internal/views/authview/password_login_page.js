const nfcLoginBtn = document.querySelector(".nfc-login-button");

nfcLoginBtn.addEventListener("click", async () => {
  try {
    const hasNFC = await AppNFC.checkDeviceHasNFC();
    if (!hasNFC) {
      console.info("This device does not support NFC");
      //   setNfcEnabled(false);
      return;
    }
    const { tagContent } = await AppNFC.readNFCTag();

    const form = document.querySelector("#login-form");

    const encryptedInput = form.querySelector("input[name='EncryptedCredentials']");
    if (!encryptedInput) {
      console.error("EncryptedCredentials input not found");
      return;
    }

    encryptedInput.value = tagContent;

    form.submit();
  } catch (e) {
    console.error(e.message);
  }
});
