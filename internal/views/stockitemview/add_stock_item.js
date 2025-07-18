const skuGenerateForm = document.getElementById("sku-generate-form");

async function handleSkuGenerate(form) {
  const formData = new FormData(form);
  const params = new URLSearchParams(formData).toString();

  const response = await fetch("/stock-items/sku-generator?" + params, {
    method: "GET",
    headers: {
      "X-Requested-With": "XMLHttpRequest",
    },
  });

  const html = await response.text();
  document.getElementById("sku-preview").innerHTML = html;
}

skuGenerateForm.addEventListener("change", async function (e) {
  const form = e.target.closest("form");
  handleSkuGenerate(form);
});
