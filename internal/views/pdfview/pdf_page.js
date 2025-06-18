document.querySelector("form").addEventListener("submit", async function (e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);

  const res = await fetch(form.action || window.location.href, {
    method: "POST",
    body: formData,
  });

  if (!res.ok) {
    alert("Failed to generate PDF");
    return;
  }

  const blob = await res.blob();
  const blobUrl = URL.createObjectURL(blob);
  window.open(blobUrl, "_blank"); // ✅ opens in new tab
});
