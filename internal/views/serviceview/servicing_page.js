const storageKey = "service-ownership-team-ids";
const paramKey = "ServiceOwnershipTeamIDs";
const form = document.getElementById("service-team-form");
const wrapper = document.getElementById("service-team-select-wrapper");
const teamSelect = document.querySelector('[data-name="ServiceOwnershipTeamIDs"]');

const selectedFromBackend =
  wrapper?.dataset.selected?.split(",").filter(Boolean).sort() || [];

if (form && wrapper && teamSelect) {
  if (selectedFromBackend.length === 0) {
    const storedTeams = JSON.parse(
      localStorage.getItem(storageKey) || "[]"
    );

    if (storedTeams.length > 0) {
      const searchParams = new URLSearchParams(window.location.search);
      const existingParams = searchParams.getAll(paramKey);
      const storedTeamsSorted = [...storedTeams].sort();

      if (existingParams.length === 0) {
        storedTeamsSorted.forEach((team) => {
          searchParams.append(paramKey, team);
        });

        const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
        window.location.replace(newUrl);
      }
    }
  }
}

function handleServiceTeamSelectChange(event) {
  if (!form || !teamSelect) {
    return;
  }

  const selected = Array.from(
    teamSelect.querySelectorAll(".select-hidden-inputs input")
  ).map((el) => el.value);

  if (selected.length === 0) {
    localStorage.removeItem(storageKey);
  } else {
    localStorage.setItem(storageKey, JSON.stringify(selected));
  }

  form.submit();
}

window.handleServiceTeamSelectChange = handleServiceTeamSelectChange;
