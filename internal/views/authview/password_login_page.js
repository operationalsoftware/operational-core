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
    // await loginUserWithEncryptedCredentials(tagContent);
    // refetchUser();
  } catch (e) {
    console.error(e.message);
  }
});
