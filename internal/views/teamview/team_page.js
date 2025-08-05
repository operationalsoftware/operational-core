function updateAndon(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const userId = targetBtn.dataset.id;
  const teamId = targetBtn.dataset.teamId;
  const username = targetBtn.dataset.username;

  confirmUpdate = confirm(
    `Are you sure you want to delete "${username}" from this team?`
  );

  if (confirmUpdate) {
    fetch(`/teams/${teamId}/delete/${userId}`, {
      method: "DELETE",
    }).then((res) => {
      if (res.ok) {
        window.location.href = `/teams/${teamId}`;
      } else {
        alert("Failed to update Andon.");
      }
    });
  }
}
