const searchInput = document.getElementById("search-input");
const resultsContainer = document.querySelector(".search-results");
const checkboxes = document.querySelectorAll(".filter-checkbox");
const searchForm = document.getElementById("search-form");

let timeoutId;

function getFilters() {
  return Array.from(checkboxes)
    .filter((cb) => cb.checked)
    .map((cb) => cb.value);
}

function searchAndUpdateQuery() {
  const query = searchInput.value.trim();
  const filters = getFilters();

  const newUrl = new URL(window.location);
  if (query) {
    newUrl.searchParams.set("Q", query);
  } else {
    newUrl.searchParams.delete("Q");
  }

  if (filters.length > 0) {
    newUrl.searchParams.set("E", filters.join(","));
  } else {
    newUrl.searchParams.delete("E");
  }

  // window.history.replaceState({}, "", newUrl);

  // Debounce the API call
  clearTimeout(timeoutId);
  timeoutId = setTimeout(() => {}, 300);
}

document.addEventListener("DOMContentLoaded", () => {
  const searchEntities = getParamEntities();

  if (searchEntities.length > 0) return;

  const storageSearchEntities = JSON.parse(
    localStorage.getItem("search-entities") || "[]"
  );

  checkboxes.forEach((checkbox) => {
    const searchEntity = checkbox.value;
    checkbox.checked =
      storageSearchEntities.length > 0
        ? storageSearchEntities.includes(searchEntity)
        : true;
  });
});

checkboxes.forEach((cb) =>
  cb.addEventListener("change", () => {
    const searchEntities = [...checkboxes]
      .filter((cb) => cb.checked)
      .map((cb) => cb.dataset.type);

    localStorage.setItem("search-entities", JSON.stringify(searchEntities));
  })
);

function getParamEntities() {
  const urlParams = new URLSearchParams(window.location.search);
  const searchEntities = urlParams.getAll("E");

  return searchEntities;
}
