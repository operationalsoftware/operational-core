const searchInput = document.getElementById("search-input");
const resultsContainer = document.querySelector(".search-results");
const checkboxes = document.querySelectorAll(".filter-checkbox");

let timeoutId;

function getFilters() {
  return Array.from(checkboxes)
    .filter((cb) => cb.checked)
    .map((cb) => cb.value);
}

function updateQueryAndSearch() {
  const query = searchInput.value.trim();
  const filters = getFilters();

  const newUrl = new URL(window.location);
  if (query) {
    newUrl.searchParams.set("q", query);
  } else {
    newUrl.searchParams.delete("q");
  }

  if (filters.length > 0) {
    newUrl.searchParams.set("types", filters.join(","));
  } else {
    newUrl.searchParams.delete("types");
  }

  window.history.replaceState({}, "", newUrl);

  // Debounce the API call
  clearTimeout(timeoutId);
  timeoutId = setTimeout(() => {
    if (query) {
      fetch(
        `/search-results?q=${encodeURIComponent(
          query
        )}&types=${encodeURIComponent(filters.join(","))}`
      )
        .then((res) => res.json())
        .then((data) => displayResults(data))
        .catch((err) => {
          console.error("Search failed", err);
          resultsContainer.innerHTML = `<p>Error fetching results</p>`;
        });
    } else {
      resultsContainer.innerHTML = "";
    }
  }, 300);
}

searchInput.addEventListener("input", async (e) => {
  updateQueryAndSearch();
});

checkboxes.forEach((cb) =>
  cb.addEventListener("change", () => {
    updateQueryAndSearch();
  })
);

function displayResults(results) {
  if (!results || Object.keys(results).length === 0) {
    resultsContainer.innerHTML = "<p>No results found.</p>";
    return;
  }

  let html = "";

  for (const [type, items] of Object.entries(results)) {
    if (items.length === 0) continue;

    html += `<h3 class="mt-4 text-lg font-semibold capitalize">${type} Results</h3>`;
    html += `<ul class="mb-4 space-y-2">`;

    items.forEach((item) => {
      // You can customize the display below based on the structure of each item
      const name = item.name || item.title || "Unnamed";
      html += `<li class="search-result-item p-2 border rounded-md">${name}</li>`;
    });

    html += `</ul>`;
  }

  resultsContainer.innerHTML = html || "<p>No results found.</p>";
}
