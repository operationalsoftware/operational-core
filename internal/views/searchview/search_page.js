const searchInput = document.getElementById("search-input");
const resultsContainer = document.querySelector(".search-results");
const checkboxes = document.querySelectorAll(".filter-checkbox");
const searchForm = document.getElementById("search-form");
const hiddenTypesInput = document.getElementById("types-input");

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
  timeoutId = setTimeout(() => {
    // searchForm.submit();
    // if (query) {
    //   fetch(
    //     `/search-results?q=${encodeURIComponent(
    //       query
    //     )}&types=${encodeURIComponent(filters.join(","))}`
    //   )
    //     .then((res) => res.json())
    //     .then((data) => displayResults(data))
    //     .catch((err) => {
    //       console.error("Search failed", err);
    //       injectHTML(
    //         ".search-results",
    //         `<p class="placeholder">Error fetching results</p>`
    //       );
    //     });
    // } else {
    //   injectHTML(
    //     ".search-results",
    //     `<p class="placeholder">No Search results.</p>`
    //   );
    // }
  }, 300);
}

document.addEventListener("DOMContentLoaded", () => {
  // searchAndUpdateQuery();

  const searchEntities = getParamEntities();

  console.log(searchEntities);

  if (searchEntities.length === 0) {
    const searchQuery = JSON.parse(localStorage.getItem("searchQuery"));
    console.log("cangedd", searchEntities);
    console.log("local", searchQuery);
  }
});

checkboxes.forEach((cb) =>
  cb.addEventListener("change", () => {
    const searchEntities = getParamEntities();

    console.log("cangedd", searchEntities);

    localStorage.setItem("searchQuery", JSON.stringify(searchEntities));
  })
);

function getParamEntities() {
  const urlParams = new URLSearchParams(window.location.search);
  const searchEntities = urlParams.getAll("E");

  return searchEntities;
}
