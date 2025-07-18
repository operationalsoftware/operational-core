function filterList(e) {
  const input = e.target;
  const searchText = e.target.value.toLowerCase();
  const list = input.closest(".sku-list");

  const itemsContainer = list.querySelector(".sku-list-items");
  const items = itemsContainer.querySelectorAll(".sku-list-item");

  items.forEach((item) => {
    const label = item
      .querySelector(".sku-item-label")
      .textContent.toLowerCase();
    if (label.includes(searchText)) {
      item.style.display = "";
    } else {
      item.style.display = "none";
    }
  });
}

function deleteSKUItem(e) {
  e.preventDefault();
  const btn = e.target;
  const skuItem = btn.closest(".sku-list-item");
  const skuField = skuItem.dataset.field;
  const skuLabel = skuItem.dataset.label;
  const code = skuItem.dataset.code;

  confirmDelete = confirm(
    `Are you sure you want to delete ${code} \u2013 ${skuLabel} ?`
  );

  if (confirmDelete) {
    fetch(`/stock-items/sku-config/${skuField}/${code}`, {
      method: "DELETE",
    }).then((res) => {
      if (res.ok) {
        skuItem.remove();
        window.location.href = "/stock-items/sku-config?toast=deleted";
      } else {
        alert("Failed to delete SKU config.");
      }
    });
  }
}
