function updateAndon(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const andonId = targetBtn.dataset.id;
  const andonAction = targetBtn.dataset.action;

  confirmUpdate = confirm(
    `Are you sure you want to ${andonAction} this Andon?`
  );

  if (confirmUpdate) {
    fetch(`/andons/${andonId}/${andonAction}/update`, {
      method: "POST",
    }).then((res) => {
      if (res.ok) {
        window.location.href = `/andons/${andonId}`;
      } else {
        alert("Failed to update Andon.");
      }
    });
  }
}
