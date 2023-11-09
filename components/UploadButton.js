(function () {
  const input = me(".upload-button input[type=file]");
  const label = me(".upload-button label");
  const info = me(".upload-button .file-info");
  input.on("change", (e) => {
    const file = e.target.files[0];
    if (file) {
      label.textContent = file.name;
      info.textContent = `${file.size} bytes`;
    } else {
      label.textContent = "Choose file";
      info.textContent = "No file chosen";
    }
  });
})();
