setInterval(() => {
  console.log("page reloading");
  window.location.reload();
}, 30000);

const teamSelect = document.querySelector('[data-name="AndonTeams"]');
const wrapper = document.getElementById("search-select-wrapper");
const form = document.getElementById("team-form");

const selectedFromBackend =
  wrapper.dataset.selected?.split(",").filter(Boolean).sort() || [];

if (selectedFromBackend.length === 0) {
  const storedTeams = JSON.parse(localStorage.getItem("andon-teams") || "[]");

  if (storedTeams.length > 0) {
    const searchParams = new URLSearchParams(window.location.search);
    const existingParam = searchParams.get("AndonTeams");

    const storedTeamsSorted = [...storedTeams].sort();

    if (!existingParam) {
      storedTeamsSorted.forEach((team) => {
        searchParams.append("AndonTeams", team);
      });

      const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
      window.location.replace(newUrl);
    }
  }
}

function handleTeamSelectChange(e) {
  const teamSelect = e.target;

  const selected = Array.from(
    teamSelect.querySelectorAll(".select-hidden-inputs input")
  ).map((el) => el.value);

  if (selected.length === 0) {
    localStorage.removeItem("andon-teams");
  } else {
    localStorage.setItem("andon-teams", JSON.stringify(selected));
  }

  form.submit();
}
