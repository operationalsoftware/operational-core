const searchInput = document.getElementById("search-input");
const resultsContainer = document.querySelector(".search-results");
const checkboxes = document.querySelectorAll(".filter-checkbox");

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
          resultsContainer.innerHTML = `<p class="placeholder">Error fetching results</p>`;
        });
    } else {
      resultsContainer.innerHTML = `<p class="placeholder">No Search results.</p>`;
    }
  }, 300);
}

document.addEventListener("DOMContentLoaded", () => {
  searchAndUpdateQuery();
});

searchInput.addEventListener("input", async (e) => {
  searchAndUpdateQuery();
});

checkboxes.forEach((cb) =>
  cb.addEventListener("change", () => {
    searchAndUpdateQuery();
  })
);

function displayResults(results) {
  if (!results || Object.keys(results).length === 0) {
    resultsContainer.innerHTML = `<p class="placeholder">No results found.</p>`;
    return;
  }

  console.log(results);

  let html = "";

  for (const [type, items] of Object.entries(results)) {
    if (!items || items.length === 0) continue;

    const title = type.charAt(0).toUpperCase() + type.slice(1);

    html += `<h3 class="result-type-heading">${title} Results</h3>`;
    html += `<ul class="result-group">`;

    items.forEach(({ data }) => {
      if (type === "user") {
        const fullName = `${data.first_name} ${data.last_name}`;
        html += `
          <li class="search-result-item">
            <strong>${fullName}</strong> <br>
            <span>Username: ${data.username}</span><br>
            <span>Email: ${data.email}</span>
          </li>
        `;
      } else if (type === "batch") {
        html += `
          <li class="search-result-item">
            <strong>Batch #: ${data.batch_number}</strong><br>
            <span>Works Order #: ${data.works_order_number}</span><br>
            <span>Part #: ${data.part_number}</span>
          </li>
        `;
      }
    });

    html += `</ul>`;
  }

  resultsContainer.innerHTML =
    html || `<p class="placeholder">No results found.</p>`;
}
