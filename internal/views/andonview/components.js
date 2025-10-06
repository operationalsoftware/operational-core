function updateAndon(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const andonId = targetBtn.dataset.id;
  const andonAction = targetBtn.dataset.action;
  const returnTo = targetBtn.dataset.returnTo;

  confirmUpdate = confirm(
    `Are you sure you want to ${andonAction} this Andon?`
  );

  if (confirmUpdate) {
    fetch(`/andons/${andonId}/${andonAction}/update`, {
      method: "POST",
    }).then((res) => {
      if (res.ok) {
        if (returnTo) {
          window.location.href = returnTo;
          return;
        }
        window.location.href = "/andons";
      } else {
        alert("Failed to update Andon");
      }
    });
  }
}
