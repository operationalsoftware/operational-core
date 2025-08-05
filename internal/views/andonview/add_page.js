const selectEl = document.querySelector(".search-select");
const assignedTeamInput = document.getElementById("assigned-team");

function updateAssignedTeam() {
  const selectedOption = selectEl.querySelector(".select-option.selected");
  if (selectedOption) {
    assignedTeamInput.value = selectedOption.dataset.team || "";
  } else {
    assignedTeamInput.value = "";
  }
}

// Initial call on page load
updateAssignedTeam();

// Then listen to clicks on options to update
selectEl.querySelector(".select-options").addEventListener("click", (e) => {
  const option = e.target.closest(".select-option");
  if (!option) return;

  // Remove selected class from all options (single mode)
  selectEl.querySelectorAll(".select-option.selected").forEach((el) => {
    el.classList.remove("selected");
  });

  // Add selected class to clicked
  option.classList.add("selected");

  updateAssignedTeam();
});
